package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"html/template"

	"github.com/gosimple/slug"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/citations"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/sidenotes"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/speakers"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/timecodes"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/otiai10/copy"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

const (
	BUILD_DIR     = "build"
	TEMPLATES_DIR = "lib/templates"
	DOCUMENTS_DIR = "./documents"
)

var BUILD_START = time.Now()

type DocumentFrontMatter struct {
	Title          string
	Series         string
	Source         string
	SourceDuration string
	Speakers       map[string]string
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
	Speakers   []string
	Source     struct {
		Url      string
		Duration time.Duration
	}
	Transcription struct {
		Url    string
		Date   *time.Time
		Author string
	}
	Citations   citations.AggregatedCitations
	EditLink    string
	ContactLink string
	ProjectLink string
	Date        time.Time
	Contents    template.HTML
	Global      GlobalData
}

type Citations struct {
	Count         int
	People        map[citations.Person][]Document
	SciencePapers map[citations.SciencePaper][]Document
	ExternalLinks map[citations.ExternalLink][]Document
	Books         map[citations.Book][]Document
}

type GlobalData struct {
	Documents   []Document
	Citations   Citations
	ProjectLink string
	ContactLink string
	BuildTime   time.Time
	ProjectName string
}

func main() {
	if len(os.Args) >= 2 {
		arg := os.Args[1]
		switch arg {
		case "dev":
			utils.DownloadBinariesIfAbsentAndExecuteLast("pagefind", "devd", "modd")
			os.Exit(1)
		case "clean":
			utils.RemoveBinaryDir()
			fmt.Println("Removed " + utils.BINARY_DIR)
			utils.PanicOnErr(os.RemoveAll(BUILD_DIR))
			fmt.Println("Removed " + BUILD_DIR)
			os.Exit(1)
		case "build":
			break
		default:
			panic(fmt.Sprintf("Unrecognised argument '%s' options are 'build', 'dev', or'clean'", arg))
		}
	}

	fmt.Printf("Building to ./%s\n", BUILD_DIR)
	utils.PanicOnErr(os.RemoveAll(BUILD_DIR))

	documents := func() []Document {
		documentRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})-(.*).md`)
		markdown := goldmark.New(
			goldmark.WithExtensions(
				meta.New(meta.WithStoresInDocument()),
				sidenotes.Sidenotes,
				citations.Citations,
				timecodes.Timecodes,
				speakers.Speakers,
			),
		)
		markdownFiles := func() []string {
			markdownFiles := []string{}
			fs.WalkDir(os.DirFS(DOCUMENTS_DIR), ".", func(filePath string, d fs.DirEntry, err error) error {
				utils.PanicOnErr(err)
				if !d.IsDir() {
					markdownFiles = append(markdownFiles, path.Join(DOCUMENTS_DIR, filePath))
				}
				return nil
			})
			return markdownFiles
		}()
		var wg sync.WaitGroup
		documentsChannel := make(chan Document, len(markdownFiles))
		for _, filePath := range markdownFiles {
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				document := Document{}
				outputFileName, date, slug := func() (string, time.Time, string) {
					matches := documentRegex.FindStringSubmatch(path.Base(filePath))
					if len(matches) < 2 {
						panic(filePath + ": filename does not match pattern YYYY-MM-DD-title-md")
					}
					slug := slug.Make(matches[2])
					date, err := time.Parse("2006-01-02", matches[1])
					utils.PanicOnErr(err)
					return slug + "/index.html", date, slug
				}()
				postMarkdownHtml, frontMatter, citations := func() (string, DocumentFrontMatter, citations.AggregatedCitations) {
					var html bytes.Buffer
					context := parser.NewContext()
					rawMarkdown := utils.ReturnOrPanic(os.ReadFile(filePath))
					utils.PanicOnErr(markdown.Convert(rawMarkdown, &html, parser.WithContext(context)))
					frontMatter := func() DocumentFrontMatter {
						data := func() map[string]interface{} {
							if data := meta.Get(context); data != nil {
								return data
							}
							return map[string]interface{}{}
						}()
						var frontMatter DocumentFrontMatter
						err := mapstructure.Decode(data, &frontMatter)
						utils.PanicOnErr(err)
						return frontMatter
					}()
					return html.String(), frontMatter, citations.Get(context)
				}()
				document.InputPath = filePath
				document.EditLink = "https://github.com/marcuswhybrow/ray-peat-rodeo/edit/main/" + document.InputPath
				document.OutputPath = outputFileName
				document.Slug = slug
				document.Date = date
				document.Title = frontMatter.Title
				document.Series = frontMatter.Series
				document.Source.Url = frontMatter.Source
				document.Speakers = func() []string {
					values := make([]string, 0, len(frontMatter.Speakers))
					for _, v := range frontMatter.Speakers {
						values = append(values, v)
					}
					return values
				}()
				document.Transcription.Url = frontMatter.Transcription.Source
				document.Transcription.Author = frontMatter.Transcription.Author
				document.Transcription.Date = func() *time.Time {
					if frontMatter.Transcription.Date == "" {
						return nil
					}
					transcriptionDate, err := time.Parse("2006-01-02", frontMatter.Transcription.Date)
					utils.PanicOnErr(err)
					return &transcriptionDate
				}()
				document.Citations = citations
				document.Contents = func() template.HTML {
					var postGoTemplateHtml bytes.Buffer
					utils.ReturnOrPanic(template.New("markdown").Parse(string(postMarkdownHtml))).Execute(&postGoTemplateHtml, document)
					return template.HTML(postGoTemplateHtml.String())
				}()
				documentsChannel <- document
			}(filePath)
		}
		wg.Wait()
		close(documentsChannel)
		documents := []Document{}
		for document := range documentsChannel {
			documents = append(documents, document)
		}
		sort.Slice(documents, func(i, j int) bool {
			return documents[i].Date.After(documents[j].Date)
		})
		return documents
	}()

	globalData := GlobalData{
		Documents: documents,
		Citations: func() Citations {
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
		}(),
		ProjectLink: "https://github.com/marcuswhybrow/ray-peat-rodeo",
		ContactLink: "/contact",
		BuildTime:   BUILD_START,
		ProjectName: "Ray Peat Rodeo",
	}

	type Page struct {
		outputPath string
		data       any
		template   string
	}

	pageData := struct{ Global GlobalData }{Global: globalData}
	pages := []Page{
		{
			outputPath: "index.html",
			data:       pageData,
			template:   "home.tmpl",
		},
		{
			outputPath: "contact/index.html",
			data:       pageData,
			template:   "contact.tmpl",
		},
	}
	for _, document := range documents {
		document.Global = globalData
		pages = append(pages, Page{
			outputPath: document.OutputPath,
			data:       document,
			template:   "document.tmpl",
		})
	}

	var wg sync.WaitGroup
	for _, page := range pages {
		wg.Add(1)
		go func(page Page) {
			defer wg.Done()
			templatePath := path.Join(TEMPLATES_DIR, page.template)
			outputPath := func() string {
				outputPath := path.Join(BUILD_DIR, page.outputPath)
				utils.PanicOnErr(os.MkdirAll(filepath.Dir(outputPath), os.ModePerm))
				return outputPath
			}()
			outputFile := utils.ReturnOrPanic(os.Create(outputPath))
			defer outputFile.Close()
			utils.ReturnOrPanic(outputFile.WriteString(func() string {
				templates := utils.ReturnOrPanic(template.ParseFiles("lib/templates/base.tmpl", templatePath))
				var html bytes.Buffer
				utils.PanicOnErr(templates.ExecuteTemplate(&html, "base", page.data))
				return html.String()
			}()))
			fmt.Println(outputPath)
		}(page)
	}
	wg.Wait()

	utils.PanicOnErr(copy.Copy("lib/assets", path.Join(BUILD_DIR, "assets")))
	fmt.Println(time.Since(BUILD_START))
	utils.DownloadBinariesIfAbsentAndExecuteLast("pagefind --source " + BUILD_DIR)
	fmt.Println(time.Since(BUILD_START))
}
