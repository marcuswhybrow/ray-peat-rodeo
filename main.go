package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"html/template"

	"github.com/PuerkitoBio/goquery"
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
	"github.com/yuin/goldmark/extension"
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

func getFilePathsInDir(directory string) []string {
	filePaths := []string{}
	fs.WalkDir(os.DirFS(directory), ".", func(filePath string, d fs.DirEntry, err error) error {
		utils.PanicOnErr(err)
		if !d.IsDir() {
			filePaths = append(filePaths, path.Join(directory, filePath))
		}
		return nil
	})
	return filePaths
}

func parseMarkdownAndGetContext(filePath string, gm goldmark.Markdown) (string, parser.Context) {
	var html bytes.Buffer
	context := parser.NewContext()
	rawMarkdown := utils.ReturnOrPanic(os.ReadFile(filePath))
	utils.PanicOnErr(gm.Convert(rawMarkdown, &html, parser.WithContext(context)))
	return html.String(), context
}

func decodeFrontmatter(context parser.Context) DocumentFrontMatter {
	var frontMatter DocumentFrontMatter
	utils.PanicOnErr(mapstructure.Decode(meta.Get(context), &frontMatter))
	return frontMatter
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
		case "check":
			const RPF_URL = "https://raypeatforum.com"
			RPF_AUDIO_INTERVIEW_TRANSCRIPTS_URL := utils.ReturnOrPanic(url.JoinPath(RPF_URL, "/community/forums/audio-interview-transcripts.73"))

			fmt.Println("Scraping ", RPF_AUDIO_INTERVIEW_TRANSCRIPTS_URL, "to determine the number of pages in this forum category...")

			pageOne := utils.ReturnOrPanic(goquery.NewDocument(RPF_AUDIO_INTERVIEW_TRANSCRIPTS_URL))
			totalPages := utils.ReturnOrPanic(strconv.Atoi(pageOne.Find(".pageNav-main .pageNav-page a").Last().Text()))

			fmt.Println(totalPages, " pages found. Scraping all thread details...")

			forumThreadPages := make(chan goquery.Document, totalPages)
			forumThreadPages <- *pageOne

			for n := 2; n <= totalPages; n++ {
				rpfUrl := utils.ReturnOrPanic(url.JoinPath(RPF_AUDIO_INTERVIEW_TRANSCRIPTS_URL, fmt.Sprint("page-", n)))
				go func() {
					forumThreadPages <- *utils.ReturnOrPanic(goquery.NewDocument(rpfUrl))
				}()
			}

			transcriptionSourcesOnRayPeatRodeo := func() map[string]bool {
				sources := map[string]bool{}
				gm := goldmark.New(
					goldmark.WithExtensions(
						meta.New(meta.WithStoresInDocument()),
					),
				)
				filePaths := getFilePathsInDir(DOCUMENTS_DIR)
				for _, filePath := range filePaths {
					_, context := parseMarkdownAndGetContext(filePath, gm)
					frontMatter := decodeFrontmatter(context)
					sources[frontMatter.Transcription.Source] = true
				}
				return sources
			}()

			type ForumThread struct {
				Title              string
				Url                string
				StartDate          time.Time
				IsFoundInDocuments bool
			}

			forumThreads := []ForumThread{}
			numberOnRayPeatRodeo := 0
			var oldestThreadNotInDocuments *ForumThread = nil
			for i := 0; i < totalPages; i++ {
				(<-forumThreadPages).Find(".structItem-cell.structItem-cell--main").Each(func(i int, s *goquery.Selection) {
					titleElem := s.Find(".structItem-title a")
					threadPath := titleElem.AttrOr("href", "")
					startDateStr, _ := s.Find(".structItem-startDate a time").Attr("datetime")
					threadUrl := utils.ReturnOrPanic(url.JoinPath(RPF_URL, threadPath))
					isFoundInDocuments := transcriptionSourcesOnRayPeatRodeo[threadUrl]
					startDate := utils.ReturnOrPanic(time.Parse("2006-01-02T15:04:05", startDateStr[:19]))
					forumThread := ForumThread{
						Title:              titleElem.Text(),
						Url:                threadUrl,
						StartDate:          startDate,
						IsFoundInDocuments: isFoundInDocuments,
					}
					forumThreads = append(forumThreads, forumThread)
					if isFoundInDocuments {
						numberOnRayPeatRodeo++
					} else if oldestThreadNotInDocuments == nil || oldestThreadNotInDocuments.StartDate.After(startDate) {
						oldestThreadNotInDocuments = &forumThread
					}
				})
			}

			sort.Slice(forumThreads, func(i, j int) bool {
				return forumThreads[i].StartDate.After(forumThreads[j].StartDate)
			})

			fmt.Println()

			for _, forumThread := range forumThreads {
				if forumThread.IsFoundInDocuments {
					fmt.Print("[x] ")
				} else {
					fmt.Print("[ ] ")
				}
				fmt.Println(forumThread.Title, "(", forumThread.StartDate.Format("2006-01-02"), ")")
				fmt.Println("   ", forumThread.Url+"\n")
			}

			percent := fmt.Sprintf("%00d", (numberOnRayPeatRodeo/len(forumThreads))*100)

			fmt.Println("Success! Found", len(forumThreads), "forum threads, the number found in", DOCUMENTS_DIR, "was", numberOnRayPeatRodeo, "("+percent+"%)")
			if oldestThreadNotInDocuments != nil {
				fmt.Print("Open oldest thread not in ", DOCUMENTS_DIR, " in your web browser? [Y/n]: ")
				var answer string
				fmt.Scanln(&answer)
				answer = strings.ToLower(answer)
				if answer == "y" || answer == "" {
					switch runtime.GOOS {
					case "linux":
						utils.PanicOnErr(exec.Command("xdg-open", oldestThreadNotInDocuments.Url).Start())
					case "windows":
						utils.PanicOnErr(exec.Command("rundll32", "url.dll,FileProtocolHandler", oldestThreadNotInDocuments.Url).Start())
					case "darwin":
						utils.PanicOnErr(exec.Command("open", oldestThreadNotInDocuments.Url).Start())
					default:
						panic("unsupported platform")
					}
				}
			}
			os.Exit(1)
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
				extension.Typographer,
				meta.New(meta.WithStoresInDocument()),
				sidenotes.Sidenotes,
				citations.Citations,
				timecodes.Timecodes,
				speakers.Speakers,
			),
		)
		markdownFiles := getFilePathsInDir(DOCUMENTS_DIR)
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
				postMarkdownHtml, context := parseMarkdownAndGetContext(filePath, markdown)
				frontMatter := decodeFrontmatter(context)
				citations := citations.Get(context)
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
