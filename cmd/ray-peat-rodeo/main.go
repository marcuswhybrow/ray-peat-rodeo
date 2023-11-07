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

	"github.com/marcuswhybrow/ray-peat-rodeo/cmd/ray-peat-rodeo/templates"
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
		log.Fatalf("failed to create output directory: %v", err)
	}

	markdownParser := goldmark.New(goldmark.WithExtensions(
		gmExtension.Typographer,
		meta.New(meta.WithStoresInDocument()),
		extension.Timecodes,
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

			markdown, err := os.ReadFile(filePath)
			if err != nil {
				filesChannel <- Result[*File]{value: (*File)(nil), err: err}
				return
			}

			var html bytes.Buffer
			parserContext := parser.NewContext()
			markdownParser.Convert(markdown, &html, parser.WithContext(parserContext))

			var frontMatter FrontMatter
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

			base := templates.Base(frontMatter.Source.Title, html.String())
			base.Render(context.Background(), outFile)

			file := &File{}
			file.path = filePath
			filesChannel <- Result[*File]{file, nil}
		}(filePath)
	}

	waitGroup.Wait()

	close(filesChannel)
	for result := range filesChannel {
		if result.err != nil {
			log.Panicf("Failed to construct file: '%v'", result.err)
		}

		log.Printf("Wrote '%v'", result.value.path)
	}

	fmt.Println("Done.")
}

type File struct {
	path string
}

type Result[T any] struct {
	value T
	err   error
}

type FrontMatter struct {
	Source struct {
		Title    string
		Series   string
		Url      string
		Duration string
	}
	Speakers      map[string]string
	Transcription struct {
		Source string
		Date   string
		Author string
	}
}
