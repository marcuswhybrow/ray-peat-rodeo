package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"html/template"

	"github.com/gosimple/slug"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/citations"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/sidenotes"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/speakers"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/timecodes"
	"github.com/mitchellh/mapstructure"
	"github.com/otiai10/copy"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

var buildPath = func(buildPath string) string {
	removeErr := os.RemoveAll(buildPath)
	panicOn(removeErr)
	return buildPath
}("build")

func writePage(filePath string, contents []byte) string {
	outputPath := path.Join(buildPath, filePath)
	{
		mkdirErr := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		panicOn(mkdirErr)
	}
	{
		file, createErr := os.Create(outputPath)
		panicOn(createErr)
		defer file.Close()
		_, writeErr := file.Write(contents)
		panicOn(writeErr)
		fmt.Println(outputPath)
	}
	return outputPath
}

var templatesPath = "lib/templates"

func templates(templatePaths ...string) *template.Template {
	result := make([]string, len(templatePaths))
	for i, templatePath := range templatePaths {
		result[i] = path.Join(templatesPath, templatePath+".tmpl")
	}
	template, err := template.ParseFiles(result...)
	panicOn(err)
	return template
}

var markdown = goldmark.New(
	goldmark.WithExtensions(
		meta.New(meta.WithStoresInDocument()),
		sidenotes.Sidenotes,
		citations.Citations,
		timecodes.Timecodes,
		speakers.Speakers,
	),
)

type DocumentFrontMatter struct {
	Title          string
	Series         string
	Source         string
	SourceDuration string
	Transcription  struct {
		Source string
		Date   string
		Author string
	}
}

type Document struct {
	InputPath  string
	OutputPath string
	Title      string
	Series     string
	Slug       string
	Source     struct {
		Url      string
		Duration time.Duration
	}
	Transcription struct {
		Url    string
		Date   *time.Time
		Author string
	}
	Citations   citations.CitationsContext
	EditLink    string
	ContactLink string
	ProjectLink string
	Date        time.Time
	Contents    template.HTML
	BuildDate   time.Time
}

var BuildDate = time.Now()
var ContactLink = "/contact"
var ProjectLink = "https://github.com/marcuswhybrow/ray-peat-rodeo"

func main() {
	documents := func() []Document {
		documentRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})-(.*).md`)
		documents := []Document{}
		panicOn(filepath.WalkDir("documents", func(filePath string, info fs.DirEntry, walkDirFileErr error) error {
			document := Document{}
			document.InputPath = filePath
			document.EditLink = "https://github.com/marcuswhybrow/ray-peat-rodeo/edit/main/" + document.InputPath
			document.ContactLink = ContactLink
			document.BuildDate = BuildDate
			panicOn(walkDirFileErr)
			if info.IsDir() {
				return nil
			}
			fileName := info.Name()
			outputFileName, date, slug := func() (string, time.Time, string) {
				matches := documentRegex.FindStringSubmatch(fileName)
				if len(matches) < 2 {
					panic(filePath + ": filename does not match pattern YYYY-MM-DD-title-md")
				}
				slug := slug.Make(matches[2])
				date, err := time.Parse("2006-01-02", matches[1])
				panicOn(err)
				return slug + "/index.html", date, slug
			}()
			document.OutputPath = outputFileName
			document.Slug = slug
			document.Date = date
			markdownInput := func() []byte {
				preTemplate, readFileErr := os.ReadFile(filePath)
				panicOn(readFileErr)
				t, err := template.New("markdown").Parse(string(preTemplate))
				panicOn(err)
				var postTemplate bytes.Buffer
				t.Execute(&postTemplate, document)
				return postTemplate.Bytes()
			}()
			finalHtml := func() []byte {
				documentHtml, documentContext := func() (string, parser.Context) {
					var html bytes.Buffer
					context := parser.NewContext()
					markdownErr := markdown.Convert(markdownInput, &html, parser.WithContext(context))
					panicOn(markdownErr)
					return html.String(), context
				}()
				frontMatter := func() DocumentFrontMatter {
					data := func() map[string]interface{} {
						if data := meta.Get(documentContext); data != nil {
							return data
						}
						return map[string]interface{}{}
					}()
					var frontMatter DocumentFrontMatter
					err := mapstructure.Decode(data, &frontMatter)
					panicOn(err)
					return frontMatter
				}()
				document.Title = frontMatter.Title
				document.Series = frontMatter.Series
				document.Source.Url = frontMatter.Source
				document.Transcription.Url = frontMatter.Transcription.Source
				document.Transcription.Author = frontMatter.Transcription.Author
				document.Transcription.Date = func() *time.Time {
					if frontMatter.Transcription.Date == "" {
						return nil
					}
					transcriptionDate, err := time.Parse("2006-01-02", frontMatter.Transcription.Date)
					panicOn(err)
					return &transcriptionDate
				}()
				document.Citations = citations.Get(documentContext)
				document.Contents = template.HTML(documentHtml)
				finalHtml := func() []byte {
					var html bytes.Buffer
					t := templates("base", "document")
					panicOn(t.ExecuteTemplate(&html, "base", document))
					return html.Bytes()
				}()
				return finalHtml
			}()

			writePage(outputFileName, finalHtml)
			documents = append(documents, document)
			return nil
		}))

		return documents
	}()

	type Citations struct {
		Count         int
		People        map[citations.Person][]Document
		SciencePapers map[citations.SciencePaper][]Document
		ExternalLinks map[citations.ExternalLink][]Document
		Books         map[citations.Book][]Document
	}

	totalCitations := func() Citations {
		c := Citations{
			Count:         0,
			People:        map[citations.Person][]Document{},
			SciencePapers: map[citations.SciencePaper][]Document{},
			ExternalLinks: map[citations.ExternalLink][]Document{},
			Books:         map[citations.Book][]Document{},
		}

		for _, document := range documents {
			for _, book := range document.Citations.Books {
				c.Books[book] = append(c.Books[book], document)
				c.Count += 1
			}
			for _, person := range document.Citations.People {
				c.People[person] = append(c.People[person], document)
				c.Count += 1
			}
			for _, sciencePaper := range document.Citations.SciencePapers {
				c.SciencePapers[sciencePaper] = append(c.SciencePapers[sciencePaper], document)
				c.Count += 1
			}
			for _, externalLink := range document.Citations.ExternalLinks {
				c.ExternalLinks[externalLink] = append(c.ExternalLinks[externalLink], document)
				c.Count += 1
			}
		}
		return c
	}()

	writePageFromTemplate := func(pageOutpuPath, pageTemplateName string) string {
		return writePage(pageOutpuPath, func() []byte {
			t := templates("base", pageTemplateName)
			var html bytes.Buffer
			panicOn(t.ExecuteTemplate(&html, "base", map[string]interface{}{
				"documents":   documents,
				"citations":   totalCitations,
				"ProjectLink": ProjectLink,
				"ContactLink": ContactLink,
				"BuildDate":   BuildDate,
			}))
			return html.Bytes()
		}())
	}

	writePageFromTemplate("index.html", "home")
	writePageFromTemplate("contact/index.html", "contact")

	err := copy.Copy("lib/assets", path.Join(buildPath, "assets"))
	panicOn(err)
}
