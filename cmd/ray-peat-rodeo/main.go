package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/extension"
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	gmExtension "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func main() {
	fmt.Println("Building Ray Peat Rodeo...")

	assets := "./assets"
	build := "./build"

	if err := os.MkdirAll("build", os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	markdownParser := goldmark.New(goldmark.WithExtensions(
		extension.Mentions,
		gmExtension.Typographer,
		meta.New(meta.WithStoresInDocument()),
		extension.Timecodes,
		extension.Speakers,
		extension.Sidenotes,
		extension.GitHubIssues,
	))

	log.Printf("Scanning files in %v\n", assets)

	filePaths := []string{}
	err := fs.WalkDir(os.DirFS(assets), ".", func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Failed to obtain dir entry: %v", err)
		}

		if !entry.IsDir() {
			base := path.Base(filePath)
			baseLower := strings.ToLower(base)

			if baseLower == "readme.md" {
				return nil
			}

			outPath := path.Join(assets, filePath)
			filePaths = append(filePaths, outPath)
		}

		return nil
	})

	if err != nil {
		log.Panicf("Failed to read assets: %v", err)
	}

	// err = os.RemoveAll(build)
	// if err != nil {
	// 	log.Panicf("Failed to remove build directory: %v", err)
	// }

	var waitGroup sync.WaitGroup

	filesChannel := make(chan Result[*File], len(filePaths))
	for _, filePath := range filePaths {
		waitGroup.Add(1)

		go func(filePath string) {
			defer waitGroup.Done()

			fileName := filepath.Base(filePath)
			fileStem := strings.TrimSuffix(fileName, filepath.Ext(filePath))

			markdownBytes, err := os.ReadFile(filePath)
			if err != nil {
				filesChannel <- Result[*File]{value: (*File)(nil), err: err}
				return
			}

			var html bytes.Buffer
			parserContext := parser.NewContext()
			err = markdownParser.Convert(markdownBytes, &html, parser.WithContext(parserContext))
			if err != nil {
				log.Panicf("Failed to parse markdown: %v\n", err)
			}

			var frontMatter markdown.FrontMatter
			err = mapstructure.Decode(meta.Get(parserContext), &frontMatter)
			if err != nil {
				filesChannel <- Result[*File]{value: (*File)(nil), err: err}
				return
			}

			outPath := path.Join(build, fileStem, "index.html")
			os.MkdirAll(filepath.Dir(outPath), 0755)
			outFile, err := os.Create(outPath)
			if err != nil {
				filesChannel <- Result[*File]{value: (*File)(nil), err: err}
				return
			}

			parentPath := path.Dir(filePath)
			parentName := path.Base(parentPath)

			file := &File{}
			file.Path = filePath
			file.Date = fileStem[:11]
			file.Permalink = "/" + fileStem
			file.FrontMatter = frontMatter
			file.IsTodo = parentName == "todo"
			file.Markdown = string(markdownBytes)
			file.Html = html.String()

			RenderChat(file).Render(context.Background(), outFile)

			filesChannel <- Result[*File]{file, nil}
		}(filePath)
	}

	waitGroup.Wait()

	var files []*File
	close(filesChannel)
	for result := range filesChannel {
		if result.err != nil {
			log.Panicf("Failed to construct file: '%v'", result.err)
		}

		files = append(files, result.value)
		log.Printf("Wrote '%v'", result.value.Path)
	}

	indexFile, err := os.Create(path.Join(build, "index.html"))
	if err != nil {
		log.Panicf("Failed to create index.html: %v", err)
	}

	sort.Sort(ByDate(files))

	var latest []*File
	for _, file := range files {
		if !file.IsTodo {
			latest = append(latest, file)
		}
	}
	if len(latest) > 4 {
		latest = latest[:4]
	}
	sort.Sort(ByTranscriptionDate(latest))

	var humanTranscripts []*File
	for _, file := range files {
		if file.IsTodo && file.FrontMatter.Transcription.Kind == "text" {
			humanTranscripts = append(humanTranscripts, file)
		}
	}
	if len(humanTranscripts) > 2 {
		humanTranscripts = humanTranscripts[:2]
	}

	var aiTranscripts []*File
	for _, file := range files {
		if file.IsTodo && file.FrontMatter.Transcription.Kind == "auto-generated" {
			aiTranscripts = append(aiTranscripts, file)
		}
	}

	Index(latest, humanTranscripts).Render(context.Background(), indexFile)

	fmt.Println("Done.")
}

type File struct {
	FrontMatter   markdown.FrontMatter
	IsTodo        bool
	Path          string
	Date          string
	Markdown      string
	Html          string
	IssueCount    int
	MentionCounts map[string]int
	EditPermalink string
	Permalink     string
}

type ByDate []*File

func (f ByDate) Len() int { return len(f) }
func (f ByDate) Less(i, j int) bool {
	return f[i].Date > f[j].Date
}
func (f ByDate) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type ByTranscriptionDate []*File

func (f ByTranscriptionDate) Len() int { return len(f) }
func (f ByTranscriptionDate) Less(i, j int) bool {
	return f[i].FrontMatter.Transcription.Date > f[j].FrontMatter.Transcription.Date
}
func (f ByTranscriptionDate) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type Result[T any] struct {
	value T
	err   error
}
