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

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/extension"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/parser"
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	gmExtension "github.com/yuin/goldmark/extension"
	gparser "github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v3"
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

	cachePath := "./internal/http_cache.yml"
	cacheBytes, err := os.ReadFile(cachePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read cache file '%v': %v", cachePath, err))
	}

	cacheData := map[string]map[string]string{}
	err = yaml.Unmarshal(cacheBytes, cacheData)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse YAML contents of cache file '%v': %v", cachePath, err))
	}

	httpCache := cache.NewHTTPCache(cacheData)

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

			id := fileStem
			permalink := "/" + id

			var html bytes.Buffer
			parserContext := gparser.NewContext()
			parserContext.Set(markdown.PermalinkKey, permalink)
			parserContext.Set(markdown.IDKey, id)
			parserContext.Set(markdown.HTTPCache, httpCache)

			err = markdownParser.Convert(markdownBytes, &html, gparser.WithContext(parserContext))
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
			file.ID = id
			file.Path = filePath
			file.OutPath = outPath
			file.Date = fileStem[:11]
			file.Permalink = "/" + fileStem
			file.FrontMatter = frontMatter
			file.IsTodo = parentName == "todo"
			file.Markdown = markdownBytes
			file.Html = html.Bytes()
			file.Mentions = parser.GetMentions(parserContext)

			RenderChat(file).Render(context.Background(), outFile)

			filesChannel <- Result[*File]{file, nil}
		}(filePath)
	}

	waitGroup.Wait()

	cacheHitsData := httpCache.GetHits()
	cacheHits, err := yaml.Marshal(cacheHitsData)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal cache hits to YAML: %v", err))
	}

	err = os.WriteFile(cachePath, cacheHits, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to write cache hits to file '%v': %v", cachePath, err))
	}

	var allFiles []*File
	close(filesChannel)
	for result := range filesChannel {
		if result.err != nil {
			log.Panicf("Failed to construct file: '%v'", result.err)
		}

		file := result.value
		allFiles = append(allFiles, file)
		log.Printf("Wrote '%v'\n", file.OutPath)
	}

	mentions := map[ast.MentionPart]map[ast.MentionPart]map[*File][]*ast.Mention{}
	for _, file := range allFiles {
		for _, mention := range file.Mentions {
			secondaries := mentions[mention.Primary]
			if secondaries == nil {
				secondaries = map[ast.MentionPart]map[*File][]*ast.Mention{}
			}
			filesWithMention := secondaries[mention.Secondary]
			if filesWithMention == nil {
				filesWithMention = map[*File][]*ast.Mention{}
			}
			filesWithMention[file] = append(filesWithMention[file], mention)
			secondaries[mention.Secondary] = filesWithMention
			mentions[mention.Primary] = secondaries
		}
	}

	for primaryMentionPart, secondaries := range mentions {
		title := primaryMentionPart.CardinalFirst()
		title = strings.ToLower(title)
		title = strings.ReplaceAll(title, " ", "-")
		path := path.Join(build, title, "index.html")

		parent := filepath.Dir(path)
		os.MkdirAll(parent, 0755)
		mentionFile, err := os.Create(path)
		if err != nil {
			log.Panicf("Failed to create HTML for mention: %v", err)
		}

		empty := ast.MentionPart{Cardinal: "", Prefix: ""}
		primaries := secondaries[empty]
		delete(secondaries, empty)

		component := MentionPage(primaryMentionPart, primaries, secondaries)
		component.Render(context.Background(), mentionFile)
		log.Printf("Wrote %v", path)
	}

	indexFile, err := os.Create(path.Join(build, "index.html"))
	if err != nil {
		log.Panicf("Failed to create index.html: %v", err)
	}

	sort.Sort(ByDate(allFiles))

	var latest []*File
	for _, file := range allFiles {
		if !file.IsTodo {
			latest = append(latest, file)
		}
	}
	if len(latest) > 4 {
		latest = latest[:4]
	}
	sort.Sort(ByTranscriptionDate(latest))

	var humanTranscripts []*File
	for _, file := range allFiles {
		if file.IsTodo && file.FrontMatter.Transcription.Kind == "text" {
			humanTranscripts = append(humanTranscripts, file)
		}
	}
	if len(humanTranscripts) > 2 {
		humanTranscripts = humanTranscripts[:2]
	}

	var aiTranscripts []*File
	for _, file := range allFiles {
		if file.IsTodo && file.FrontMatter.Transcription.Kind == "auto-generated" {
			aiTranscripts = append(aiTranscripts, file)
		}
	}

	x := Index(latest, humanTranscripts)
	x.Render(context.Background(), indexFile)

	fmt.Println("Done.")
}

type Result[T any] struct {
	value T
	err   error
}

type File struct {
	FrontMatter   markdown.FrontMatter
	IsTodo        bool
	Path          string
	ID            string
	OutPath       string
	Date          string
	Markdown      []byte
	Html          []byte
	IssueCount    int
	MentionCounts map[string]int
	EditPermalink string
	Permalink     string
	Mentions      []*ast.Mention
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
