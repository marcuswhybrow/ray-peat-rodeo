package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	git "github.com/libgit2/git2go/v34"
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

	// 🗃 Catalog

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

	catalog := &Catalog{
		MarkdownParser:    markdownParser,
		HttpCache:         httpCache,
		AvatarPaths:       avatarPaths,
		ByMentionable:     ByMentionable[ByFile[Mentions]]{},
		ByMentionablePart: ByPart[ByPart[ByFile[Mentions]]]{},
		Files:             []*File{},
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

	parallel(catalog.Files, func(file *File) error {
		file.Render()
		if err != nil {
			return fmt.Errorf("Failed to render file '%v': %v", file.Path, err)
		}
		return nil
	})

	catalog.SortFilesByDate()
	catalog.RenderMentionPages()
	catalog.RenderPopups()

	slices.SortFunc(completedFiles, filesByDateAdded)

	progress := float32(len(completedFiles)) / float32(len(catalog.Files))

	var latestFile *File = nil
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

	blogPage, _ := makePage("blog")
	component := BlogArchive(blogPosts)
	component.Render(context.Background(), blogPage)

	// 🏠 Homepage

	repo, err := git.InitRepository(".", false)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialise git repository for current working directory. %v", err))
	}
	repoHead, err := repo.Head()
	if err != nil {
		panic(fmt.Sprintf("Failed to get HEAD of git respository in current working directory. %v", err))
	}
	repoHeadObj, err := repoHead.Peel(git.ObjectCommit)
	if err != nil {
		panic(fmt.Sprintf("HEAD of repository in the current working directory is not a commit. %v", err))
	}
	repoHeadCommit, err := repoHeadObj.AsCommit()
	if err != nil {
		panic(fmt.Sprintf("Could not get HEAD of repository in current working directory as a commit. %v", err))
	}

	githubUserSearch := "https://api.github.com/search/users?q=" + repoHeadCommit.Author().Email
	githubLogin := <-httpCache.GetJSON(githubUserSearch, "login", func(res *http.Response) string {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(fmt.Sprintf("Failed to read body of HTTP response for url '%v': %v", githubUserSearch, err))
		}

		data := GithubUserSearchData{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal JSON response for url '%v': %v", githubUserSearch, err))
		}

		if len(data.Items) > 0 {
			return data.Items[0].Login
		}

		// Email address matches no given GitHub has no account
		return ""
	})

	githubAvatar := ""
	if len(githubLogin) > 0 {
		githubAvatar = fmt.Sprintf("https://github.com/%v.png", githubLogin)
	}

	gitMessageSanitized := ""
	gitMessageLines := strings.Split(repoHeadCommit.Message(), "\n")
	if len(gitMessageLines) > 0 {
		gitMessageSanitized = gitMessageLines[0]
	}

	latestCommit := GitCommit{
		SanitizedMessage: gitMessageSanitized,
		Commit:           repoHeadCommit,
		GitHub: GitCommitGitHubData{
			CommitLink:            "https://github.com/marcuswhybrow/ray-peat-rodeo/commit/" + repoHeadCommit.Id().String(),
			AuthorRepoCommitsLink: "https://github.com/marcuswhybrow/ray-peat-rodeo/commits?author=" + githubLogin,
			AuthorProfileLink:     "https://github.com/" + githubLogin,
			AuthorLogin:           githubLogin,
			AuthorAvatar:          githubAvatar,
		},
	}

	indexPage, _ := makePage(".")
	component = Index(latestCommit, catalog.Files, latestFile, progress, latestBlogPost)
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
			fmt.Printf(" for %v", url)
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
func makePage(outPath string) (*os.File, string) {
	return makeFile(path.Join(outPath, "index.html"))
}

// Convenience function to output file
func makeFile(outPath string) (*os.File, string) {
	buildPath := path.Join(OUTPUT, outPath)
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

// Descibes data from GitHub's REST API relevant to a particular Git commit
type GitCommitGitHubData struct {
	CommitLink            string
	AuthorRepoCommitsLink string
	AuthorProfileLink     string
	AuthorLogin           string
	AuthorAvatar          string
}

// Describes a Git commit
type GitCommit struct {
	SanitizedMessage string
	Commit           *git.Commit
	GitHub           GitCommitGitHubData
}

// Descibes data returned by GitHub's REST API user search endpoint
type GithubUserSearchData struct {
	Items []GithubUserSearchUser
}

// Descibes data returned by GitHub's REST API user search for a single user.
type GithubUserSearchUser struct {
	ID       int64
	Login    string
	HTML_url string
}
