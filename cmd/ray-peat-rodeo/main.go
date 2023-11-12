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

	files := []string{}
	err := fs.WalkDir(os.DirFS(assets), ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Failed to obtain dir entry: %v", err)
		}

		if !d.IsDir() {
			files = append(files, path.Join(assets, filePath))
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

	filesChannel := make(chan Result[*File], len(files))
	for _, filePath := range files {
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
			markdownParser.Convert(markdownBytes, &html, parser.WithContext(parserContext))

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

			file := &File{}
			file.Path = filePath
			file.FrontMatter = frontMatter
			file.IsTodo = false
			file.Markdown = string(markdownBytes)
			file.Html = html.String()

			RenderChat(file).Render(context.Background(), outFile)

			filesChannel <- Result[*File]{file, nil}
		}(filePath)
	}

	waitGroup.Wait()

	close(filesChannel)
	for result := range filesChannel {
		if result.err != nil {
			log.Panicf("Failed to construct file: '%v'", result.err)
		}

		log.Printf("Wrote '%v'", result.value.Path)
	}

	fmt.Println("Done.")
}

type File struct {
	FrontMatter   markdown.FrontMatter
	IsTodo        bool
	Path          string
	Markdown      string
	Html          string
	IssueCount    int
	MentionCounts map[string]int
	EditPermalink string
}

type Result[T any] struct {
	value T
	err   error
}
