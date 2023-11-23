package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/extension"
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	gmExtension "github.com/yuin/goldmark/extension"
	gparser "github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v3"
)

const BUILD = "./build"
const ASSETS = "./assets"
const CACHE_PATH = "./internal/http_cache.yml"

func main() {
	fmt.Println("Building Ray Peat Rodeo...")

	if err := os.MkdirAll("build", os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	markdownParser := goldmark.New(
		goldmark.WithExtensions(
			extension.Mentions,
			gmExtension.Typographer,
			meta.New(meta.WithStoresInDocument()),
			extension.Timecodes,
			extension.Speakers,
			extension.Sidenotes,
			extension.GitHubIssues,
		),
	)

	log.Printf("Scanning files in %v\n", ASSETS)

	filePaths := []string{}
	err := fs.WalkDir(os.DirFS(ASSETS), ".", func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Failed to obtain dir entry: %v", err)
		}

		if !entry.IsDir() {
			base := path.Base(filePath)
			baseLower := strings.ToLower(base)

			if baseLower == "readme.md" {
				return nil
			}

			outPath := path.Join(ASSETS, filePath)
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

	// INIT HTTP CACHE

	cacheBytes, err := os.ReadFile(CACHE_PATH)
	if err != nil {
		panic(fmt.Sprintf("Failed to read cache file '%v': %v", CACHE_PATH, err))
	}

	cacheData := map[string]map[string]string{}
	err = yaml.Unmarshal(cacheBytes, cacheData)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse YAML contents of cache file '%v': %v", CACHE_PATH, err))
	}

	httpCache := cache.NewHTTPCache(cacheData)

	// PARSE MARKDOWN ASYNC

	var waitGroup sync.WaitGroup

	filesChannel := make(chan *File, len(filePaths))
	for _, filePath := range filePaths {
		waitGroup.Add(1)

		go func(filePath string) {
			defer waitGroup.Done()

			file := NewFile(filePath)

			parserContext := gparser.NewContext()
			parserContext.Set(ast.FileKey, file)
			parserContext.Set(ast.HTTPCacheKey, httpCache)

			var html bytes.Buffer
			err = markdownParser.Convert(file.Markdown, &html, gparser.WithContext(parserContext))
			if err != nil {
				log.Panicf("Failed to parse markdown for file '%v': %v", file.Path, err)
			}

			var frontMatter ast.FrontMatter
			err = mapstructure.Decode(meta.Get(parserContext), &frontMatter)
			if err != nil {
				log.Panicf("Failed to decode markdown front matter for '%v': %v", file.Path, err)
			}
			file.FrontMatter = frontMatter
			file.Html = html.Bytes()

			os.MkdirAll(filepath.Dir(file.OutPath), 0755)
			outFile, err := os.Create(file.OutPath)
			if err != nil {
				log.Panicf("Failed to create file '%v': %v", file.OutPath, err)
			}

			RenderChat(file).Render(context.Background(), outFile)

			filesChannel <- file
		}(filePath)
	}

	waitGroup.Wait()
	close(filesChannel)

	allFiles := NewFiles()
	for file := range filesChannel {
		allFiles.Add(file)
		log.Printf("Wrote '%v'\n", file.OutPath)
	}

	// CATALOG

	for primary, mentionsByFileBySecondary := range allFiles.Catalog {
		primaries := mentionsByFileBySecondary[ast.EmptyMentionablePart]
		delete(mentionsByFileBySecondary, ast.EmptyMentionablePart)

		file, outPath := createBuildHTMLFile(unencode(primary.ID()))
		component := MentionPage(primary, primaries, mentionsByFileBySecondary, httpCache)
		component.Render(context.Background(), file)
		log.Printf("Wrote %v", outPath)
	}

	// POPUPS

	for mentionable, mentionsByFile := range allFiles.Popups {
		location, _ := strings.CutPrefix(mentionable.PopupPermalink(), "/")

		otherMentionables := []ast.Mentionable{}
		for m := range allFiles.Popups {
			samePrimary := m.Primary == mentionable.Primary
			identical := m.Secondary == mentionable.Secondary
			if samePrimary && !identical {
				otherMentionables = append(otherMentionables, m)
			}
		}

		f, outPath := createBuildHTMLFile(unencode(location))
		component := MentionablePopup(mentionable, mentionsByFile, otherMentionables)
		component.Render(context.Background(), f)
		log.Printf("Wrote %v", outPath)
	}

	// HOMEPAGE

	indexFile, err := os.Create(path.Join(BUILD, "index.html"))
	if err != nil {
		log.Panicf("Failed to create index.html: %v", err)
	}

	sort.Sort(ByDate(allFiles.Files))

	var done []*File
	var todo []*File
	for _, file := range allFiles.Files {
		if !file.IsTodo {
			done = append(done, file)
		} else {
			todo = append(todo, file)
		}
	}

	latest := make([]*File, len(done))
	copy(latest, done)
	sort.Sort(ByAddedDate(latest))
	sort.Sort(ByDate(done))
	sort.Sort(ByDate(todo))

	var humanTranscripts []*File
	for _, file := range allFiles.Files {
		if file.IsTodo && file.FrontMatter.Transcription.Kind == "text" {
			humanTranscripts = append(humanTranscripts, file)
		}
	}
	if len(humanTranscripts) > 2 {
		humanTranscripts = humanTranscripts[:2]
	}

	var aiTranscripts []*File
	for _, file := range allFiles.Files {
		if file.IsTodo && file.FrontMatter.Transcription.Kind == "auto-generated" {
			aiTranscripts = append(aiTranscripts, file)
		}
	}

	progress := float32(len(latest)) / float32(len(allFiles.Files))

	component := Index(allFiles.Files, latest, humanTranscripts, progress)
	component.Render(context.Background(), indexFile)

	fmt.Println("\nEpilogue")
	fmt.Println("  ✅ Created Ray Peat Rodeo HTML in " + BUILD)

	// DUMP HTTP CACHE

	misses := httpCache.GetRequestsMissed()
	if len(misses) == 0 {
		fmt.Println("  ✅ HTTP Cache fulfilled all requests")
	} else {
		fmt.Printf("  ❌ HTTP Cache rectified %v miss(es)\n", len(misses))
		for url, keys := range misses {
			fmt.Print("    - Missed ")
			for i, key := range keys {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("'%v'", key)
			}
			fmt.Printf(" for %v", url)
		}
	}

	cacheRequests := httpCache.GetRequestsMade()
	newCache, err := yaml.Marshal(cacheRequests)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal cache hits to YAML: %v", err))
	}

	err = os.WriteFile(CACHE_PATH, newCache, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to write cache hits to file '%v': %v", CACHE_PATH, err))
	}

	fmt.Print("\nDone. 🎉\n\n")
}

// Convenience function to create path and file in build directory
func createBuildFile(outPath string) (*os.File, string) {
	buildPath := path.Join(BUILD, outPath)
	parent := filepath.Dir(buildPath)

	err := os.MkdirAll(parent, 0755)
	if err != nil {
		log.Panicf("Failed to create directory '%v': %v", parent, err)
	}

	f, err := os.Create(buildPath)
	if err != nil {
		log.Panicf("Failed to create file '%v': %v", buildPath, err)
	}

	return f, buildPath
}

// Convenience function to create dir at path with index.html inside
func createBuildHTMLFile(outPath string) (*os.File, string) {
	return createBuildFile(path.Join(outPath, "index.html"))
}

func unencode(filePath string) string {
	str, err := url.QueryUnescape(filePath)
	if err != nil {
		log.Panicf("Failed to unescape path '%v': %v", filePath, err)
	}
	return str
}
