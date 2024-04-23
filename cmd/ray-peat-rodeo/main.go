package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/extension"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	gmExtension "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

const OUTPUT = "./build"
const ASSETS = "./assets"
const CACHE_PATH = "./internal/http_cache.yml"

func main() {
	start := time.Now()

	fmt.Println("Running Ray Peat Rodeo")

	workDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to determine current working direction: %v", err)
	}
	fmt.Printf("Source: \"%v\"\n", workDir)
	fmt.Printf("Output: \"%v\"\n", OUTPUT)

	if err := os.MkdirAll(OUTPUT, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// err := os.RemoveAll(OUTPUT)
	// if err != nil {
	// 	log.Panicf("Failed to clean output directory: %v", err)
	// }

	// 📶 Get HTTP Cache from file

	cacheData, err := cache.DataFromYAMLFile(CACHE_PATH)
	if err != nil {
		log.Panicf("Failed to read cache file '%v': %v", CACHE_PATH, err)
	}

	httpCache := cache.NewHTTPCache(cacheData)

	// 👨👱 Speaker Avatars

	avatarPaths := NewAvatarPaths()

	// 🖹🧩 Markdown + Extensions

	markdownParser := goldmark.New(
		goldmark.WithRendererOptions(html.WithUnsafe()),
		goldmark.WithExtensions(
			meta.New(meta.WithStoresInDocument()),
			gmExtension.Typographer,
			extension.Mentions,
			extension.Timecodes,
			extension.Speakers,
			extension.Sidenotes,
			extension.GitHubIssues,
		),
	)

	// 🗃 Catalog

	// The catalog is a singleton for global data derived from assets.
	// It's in charge of creating assets from the source markdown files.
	// In so doing, it builds an in memory store of higher-order data.
	// This higher-order data is used to create other pages for our readers.
	catalog := &Catalog{
		MarkdownParser:    markdownParser,
		HttpCache:         httpCache,
		AvatarPaths:       avatarPaths,
		ByMentionable:     ByMentionable[ByAsset[Mentions]]{},
		ByMentionablePart: ByPart[ByPart[ByAsset[Mentions]]]{},
		Assets:            []*Asset{},
	}

	// 📂 Read Files

	fmt.Println("\n[Files]")
	fmt.Printf("Source \"%v\"\n", ASSETS)

	filePaths := files(ASSETS, ".", func(filePath string) (*string, error) {
		base := path.Base(filePath)
		baseLower := strings.ToLower(base)

		if baseLower == "readme.md" {
			return nil, nil
		}

		outPath := path.Join(ASSETS, filePath)
		return &outPath, nil
	})

	parallel(filePaths, func(filePath string) error {
		err := catalog.NewFile(filePath)
		if err != nil {
			return fmt.Errorf("Failed to retrieve file '%v': %v", filePath, err)
		}
		return nil
	})

	completedFiles := catalog.CompletedFiles()
	fmt.Printf("Found %v markdown files of which %v are completed.\n", len(filePaths), len(completedFiles))

	// 👉 Future processing to be added here 👈

	// 📝 Write files

	// When an asset filename changes, it's URL changes.
	// It's nice to redirect old URL's to the new ones.
	// N.B. this data is currently collected, but not acted upon
	redirections := map[string][]*Asset{}

	parallel(catalog.Assets, func(file *Asset) error {
		file.Write()
		if err != nil {
			return fmt.Errorf("Failed to render file '%v': %v", file.Path, err)
		}

		for _, prevPath := range file.FrontMatter.RayPeatRodeo.PrevPaths {
			existing, ok := redirections[prevPath]
			if !ok {
				existing = []*Asset{}
			}
			redirections[prevPath] = append(existing, file)
		}
		return nil
	})

	catalog.SortFilesByDate()
	catalog.RenderMentionPages()
	catalog.RenderPopups()

	slices.SortFunc(completedFiles, filesByDateAdded)

	progress := float32(len(completedFiles)) / float32(len(catalog.Assets))

	var latestFile *Asset = nil
	if len(completedFiles) > 0 {
		latestFile = completedFiles[0]
	}

	// 📢 Blog

	postPaths := files(".", "internal/blog", func(filePath string) (*string, error) {
		ext := filepath.Ext(filePath)
		if strings.ToLower(ext) != ".md" {
			return nil, nil
		}

		return &filePath, nil
	})

	blogPosts := parallel(postPaths, func(filePath string) *BlogPost {
		blogPost := NewBlogPost(filePath, avatarPaths)
		blogPost.Render()
		return blogPost
	})

	latestBlogPost := blogPosts[0]

	blogPage, _ := MakePage("blog")
	component := BlogArchive(blogPosts)
	component.Render(context.Background(), blogPage)

	// 🏠 Homepage

	indexPage, _ := MakePage(".")
	component = Index(catalog.Assets, latestFile, progress, latestBlogPost)
	component.Render(context.Background(), indexPage)

	// 📶 HTTP Cache

	fmt.Println("\n[HTTP Requests]")

	httpCacheMisses := httpCache.GetRequestsMissed()
	cacheRequests := httpCache.GetRequestsMade()

	if len(httpCacheMisses) == 0 {
		fmt.Printf("HTTP Cache fulfilled %v requests.\n", len(cacheRequests))
	} else {
		fmt.Printf("❌ HTTP Cache rectified %v miss(es):\n", len(httpCacheMisses))
		for url, keys := range httpCacheMisses {
			fmt.Print(" - Missed ")
			for i, key := range keys {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("'%v'", key)
			}
			fmt.Printf(" for %v\n", url)
		}
	}

	newCache, err := yaml.Marshal(cacheRequests)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal cache hits to YAML: %v", err))
	}

	err = os.WriteFile(CACHE_PATH, newCache, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to write cache hits to file '%v': %v", CACHE_PATH, err))
	}

	// 🏁 Done

	fmt.Printf("\nFinished in %v.\n", time.Since(start))
}

// Convenience function to output HTML page
func MakePage(outPath string) (*os.File, string) {
	return MakeFile(path.Join(outPath, "index.html"))
}

var builtFiles []string
var builtFilesMutex sync.RWMutex

// Convenience function to output file
func MakeFile(outPath string) (*os.File, string) {
	buildPath := path.Join(OUTPUT, outPath)
	parent := filepath.Dir(buildPath)

	err := os.MkdirAll(parent, 0755)
	if err != nil {
		log.Panicf("Failed to create directory '%v': %v", parent, err)
	}

	builtFilesMutex.Lock()
	if slices.Contains(builtFiles, buildPath) {
		log.Panicf(
			"Multiple writes attempted to the same build path: %v\n"+
				"  Common reasons for this include:\n"+
				"    - Two files in ./assets that have different dates in the filename, but the same wording after the date.\n"+
				"    - Two mentions that have the same name, but different capitalization.\n"+
				"    - A file in ./assets that has the same wording after the date as the wording of a mention.\n",
			buildPath,
		)
	}
	builtFiles = append(builtFiles, buildPath)
	builtFilesMutex.Unlock()

	f, err := os.Create(buildPath)
	if err != nil {
		log.Panicf("Failed to create file '%v': %v", buildPath, err)
	}

	return f, buildPath
}

// Runs func in parallel for each entry in slice and awaits all results
func parallel[Item, Result any](slice []Item, f func(Item) Result) []Result {
	var waitGroup sync.WaitGroup

	count := len(slice)
	results := make(chan Result, count)
	waitGroup.Add(count)

	for _, item := range slice {
		go func(i Item) {
			defer waitGroup.Done()
			results <- f(i)
		}(item)
	}

	waitGroup.Wait()
	close(results)

	allResults := []Result{}
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}

func files[Result any](pwd, scope string, f func(filePath string) (*Result, error)) []Result {
	results := []Result{}

	err := fs.WalkDir(os.DirFS(pwd), scope, func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Failed to walk dir: %v", err)
		}

		if !entry.IsDir() {
			result, err := f(filePath)
			if err != nil {
				return err
			}

			if result != nil {
				results = append(results, *result)
			}
		}

		return nil
	})

	if err != nil {
		log.Panicf("Failed to read directory '%v': %v", path.Join(pwd, scope), err)
	}

	return results
}
