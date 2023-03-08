package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

var BUILD_START = time.Now()

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func returnOrPanic[A any](a A, err error) A {
	panicOnErr(err)
	return a
}

var buildPath = func(buildPath string) string {
	panicOnErr(os.RemoveAll(buildPath))
	return buildPath
}("build")

func writePage(filePath string, contents []byte) string {
	outputPath := path.Join(buildPath, filePath)
	{
		panicOnErr(os.MkdirAll(filepath.Dir(outputPath), os.ModePerm))
	}
	{
		file := returnOrPanic(os.Create(outputPath))
		defer file.Close()
		returnOrPanic(file.Write(contents))
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
	return returnOrPanic(template.ParseFiles(result...))
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
	Citations   citations.CitationsContext
	EditLink    string
	ContactLink string
	ProjectLink string
	Date        time.Time
	Contents    template.HTML
	BuildDate   time.Time
}

const CONTACT_LINK = "/contact"
const PROJECT_LINK = "https://github.com/marcuswhybrow/ray-peat-rodeo"
const PROJECT_NAME = "Ray Peat Rodeo"
const BIN_DIR = "./lib/bin"
const DOCUMENTS_DIR = "./documents"

var BIN_TARGZ_URLS = map[string]string{
	"modd":     "https://github.com/cortesi/modd/releases/download/v0.8/modd-0.8-linux64.tgz",
	"devd":     "https://github.com/cortesi/devd/releases/download/v0.9/devd-0.9-linux64.tgz",
	"pagefind": "https://github.com/CloudCannon/pagefind/releases/download/v0.12.0/pagefind-v0.12.0-x86_64-unknown-linux-musl.tar.gz",
}

func downloadBinaryIfAbsent(binaryName string) bool {
	fileName := path.Join(BIN_DIR, binaryName)
	tarGzUrl := BIN_TARGZ_URLS[binaryName]
	if tarGzUrl == "" {
		panic(fmt.Sprintf("'%s' has no associated tar.gz download URL defined in BIN_TARGZ_URLS", binaryName))
	}
	if _, err := os.Stat(fileName); err == nil {
		return false
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Binary %s not installed. Downloading from %s\n", binaryName, tarGzUrl)
		resp := returnOrPanic(http.Get(tarGzUrl))
		defer resp.Body.Close()
		uncompressedStream := returnOrPanic(gzip.NewReader(resp.Body))
		tarReader := tar.NewReader(uncompressedStream)
		for true {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			panicOnErr(err)
			if header.Typeflag == tar.TypeReg && path.Base(header.Name) == path.Base(fileName) {
				panicOnErr(os.MkdirAll(path.Dir(fileName), 0755))
				outFile := returnOrPanic(os.Create(fileName))
				outFile.Chmod(755)
				returnOrPanic(io.Copy(outFile, tarReader))
				outFile.Close()
				fmt.Printf("Binary %s installed to %s\n", binaryName, fileName)
				return true
			}
		}
		panic(fmt.Sprintf("%s could not be found in %s", fileName, tarGzUrl))
	} else {
		panic(err)
	}
}

func downloadBinariesIfAbsentAndExecuteLast(commandsWithArgs ...string) {
	{
		var wg sync.WaitGroup
		for _, commandWithArgs := range commandsWithArgs {
			wg.Add(1)
			command := strings.Split(commandWithArgs, " ")[0]
			go func() {
				defer wg.Done()
				downloadBinaryIfAbsent(command)
			}()
		}
		wg.Wait()
	}
	lastCommandWithArgs := commandsWithArgs[len(commandsWithArgs)-1]
	lastCommand, lastCommandArgs := func() (string, []string) {
		lastCommandParts := strings.Split(lastCommandWithArgs, " ")
		return lastCommandParts[0], func() []string {
			if len(lastCommandParts) > 1 {
				return lastCommandParts[1:]
			}
			return []string{}
		}()
	}()
	modd := exec.Command(path.Join(BIN_DIR, lastCommand), lastCommandArgs...)
	stdout := returnOrPanic(modd.StdoutPipe())
	var wg sync.WaitGroup
	wg.Add(1)
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			log.Print(scanner.Text())
		}
		wg.Done()
	}()
	fmt.Printf("./%s\n", lastCommandWithArgs)
	panicOnErr(modd.Start())
	wg.Wait()
	panicOnErr(modd.Wait())
}

func main() {
	if len(os.Args) >= 2 {
		arg := os.Args[1]
		switch arg {
		case "dev":
			downloadBinariesIfAbsentAndExecuteLast("pagefind", "devd", "modd")
			os.Exit(1)
		case "clean":
			panicOnErr(os.RemoveAll(BIN_DIR))
			fmt.Println("Removed " + BIN_DIR)
			panicOnErr(os.RemoveAll(buildPath))
			fmt.Println("Removed " + buildPath)
			os.Exit(1)
		case "build":
			break
		default:
			panic(fmt.Sprintf("Unrecognised argument '%s' options are 'build', 'dev', or'clean'", arg))
		}
	}

	fmt.Printf("Building %s to %s\n", PROJECT_NAME, buildPath)

	documentsChannel := func() chan Document {
		documentRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})-(.*).md`)
		markdownFiles := func() []string {
			markdownFiles := []string{}
			fs.WalkDir(os.DirFS(DOCUMENTS_DIR), ".", func(filePath string, d fs.DirEntry, err error) error {
				panicOnErr(err)
				if !d.IsDir() {
					markdownFiles = append(markdownFiles, path.Join(DOCUMENTS_DIR, filePath))
				}
				return nil
			})
			return markdownFiles
		}()
		documentsChannel := make(chan Document, len(markdownFiles))
		var wg sync.WaitGroup
		for _, filePath := range markdownFiles {
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				document := Document{}
				document.InputPath = filePath
				document.EditLink = "https://github.com/marcuswhybrow/ray-peat-rodeo/edit/main/" + document.InputPath
				document.ContactLink = CONTACT_LINK
				document.BuildDate = BUILD_START
				outputFileName, date, slug := func() (string, time.Time, string) {
					matches := documentRegex.FindStringSubmatch(path.Base(filePath))
					if len(matches) < 2 {
						panic(filePath + ": filename does not match pattern YYYY-MM-DD-title-md")
					}
					slug := slug.Make(matches[2])
					date, err := time.Parse("2006-01-02", matches[1])
					panicOnErr(err)
					return slug + "/index.html", date, slug
				}()
				document.OutputPath = outputFileName
				document.Slug = slug
				document.Date = date
				markdownInput := func() []byte {
					preTemplate, readFileErr := os.ReadFile(filePath)
					panicOnErr(readFileErr)
					t, err := template.New("markdown").Parse(string(preTemplate))
					panicOnErr(err)
					var postTemplate bytes.Buffer
					t.Execute(&postTemplate, document)
					return postTemplate.Bytes()
				}()
				finalHtml := func() []byte {
					documentHtml, documentContext := func() (string, parser.Context) {
						var html bytes.Buffer
						context := parser.NewContext()
						markdownErr := markdown.Convert(markdownInput, &html, parser.WithContext(context))
						panicOnErr(markdownErr)
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
						panicOnErr(err)
						return frontMatter
					}()
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
						panicOnErr(err)
						return &transcriptionDate
					}()
					document.Citations = citations.Get(documentContext)
					document.Contents = template.HTML(documentHtml)
					finalHtml := func() []byte {
						var html bytes.Buffer
						t := templates("base", "document")
						panicOnErr(t.ExecuteTemplate(&html, "base", document))
						return html.Bytes()
					}()
					return finalHtml
				}()
				writePage(outputFileName, finalHtml)
				documentsChannel <- document
			}(filePath)
		}
		wg.Wait()
		close(documentsChannel)
		return documentsChannel
	}()

	type Citations struct {
		Count         int
		People        map[citations.Person][]Document
		SciencePapers map[citations.SciencePaper][]Document
		ExternalLinks map[citations.ExternalLink][]Document
		Books         map[citations.Book][]Document
	}

	documents, totalCitations := func() ([]Document, Citations) {
		c := Citations{
			Count:         0,
			People:        map[citations.Person][]Document{},
			SciencePapers: map[citations.SciencePaper][]Document{},
			ExternalLinks: map[citations.ExternalLink][]Document{},
			Books:         map[citations.Book][]Document{},
		}

		documents := []Document{}

		for document := range documentsChannel {
			documents = append(documents, document)
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
		return documents, c
	}()

	writePageFromTemplate := func(pageOutpuPath, pageTemplateName string) string {
		return writePage(pageOutpuPath, func() []byte {
			t := templates("base", pageTemplateName)
			var html bytes.Buffer
			panicOnErr(t.ExecuteTemplate(&html, "base", map[string]interface{}{
				"documents":   documents,
				"citations":   totalCitations,
				"ProjectLink": PROJECT_LINK,
				"ContactLink": CONTACT_LINK,
				"BuildDate":   BUILD_START,
			}))
			return html.Bytes()
		}())
	}

	writePageFromTemplate("index.html", "home")
	writePageFromTemplate("contact/index.html", "contact")
	panicOnErr(copy.Copy("lib/assets", path.Join(buildPath, "assets")))
	fmt.Printf("Pages built in %s\n", time.Since(BUILD_START))

	downloadBinariesIfAbsentAndExecuteLast("pagefind --source " + buildPath)
	fmt.Printf("Build completed in %s\n", time.Since(BUILD_START))
}
