package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"time"

	rprCatalog "github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/check"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/global"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/utils"
)

func main() {
	// tStart := time.Now()

	if len(os.Args) >= 2 {
		subcommand := os.Args[1]
		switch subcommand {
		case "check":
			check.Check()
			return
		default:
			fmt.Printf("'%v' is not a valid subcommand. Try: check\n", subcommand)
			return
		}
	}

	start := time.Now()

	fmt.Println("Running Ray Peat Rodeo")

	workDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to determine current working direction: %v", err)
	}
	fmt.Printf("Source: \"%v\"\n", workDir)
	fmt.Printf("Output: \"%v\"\n", global.BUILD_OUTPUT)

	if err := os.RemoveAll(global.BUILD_OUTPUT); err != nil {
		log.Fatalf("Failed to clean output directory: %v", err)
	}

	if err := os.MkdirAll(global.BUILD_OUTPUT, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// err := os.RemoveAll(OUTPUT)
	// if err != nil {
	// 	log.Panicf("Failed to clean output directory: %v", err)
	// }

	// 🗃 Catalog

	// log.Printf("Preamble", time.Since(tStart))

	fmt.Println("\n[Files]")
	fmt.Printf("Source \"%v\"\n", global.ASSETS)

	// The catalog is a singleton for global data derived from assets.
	// It's in charge of creating assets from the source markdown files.
	// In so doing, it builds an in memory store of higher-order data.
	// This higher-order data is used to create other pages for our readers.
	catalog := rprCatalog.NewCatalog(global.ASSETS)

	allAssets := catalog.Assets
	completedAssets := catalog.GetCompletedAssets()

	fmt.Printf("Found %v markdown files of which %v are completed.\n", len(allAssets), len(completedAssets))

	// 📝 Write files

	// log.Printf("Catalog", time.Since(tStart))

	// When an asset filename changes, it's URL changes.
	// It's nice to redirect old URL's to the new ones.
	// N.B. this data is currently collected, but not acted upon
	redirections := map[string][]*rprCatalog.Asset{}
	redirectionsMutex := sync.RWMutex{}

	utils.Parallel(catalog.Assets, func(file *rprCatalog.Asset) error {
		err := file.WriteHtml(catalog, false)
		if err != nil {
			return fmt.Errorf("Failed to render file '%v': %v", file.Path, err)
		}

		for _, prevPath := range file.FrontMatter.RayPeatRodeo.PrevPaths {
			existing, ok := redirections[prevPath]
			if !ok {
				existing = []*rprCatalog.Asset{}
			}
			redirectionsMutex.Lock()
			redirections[prevPath] = append(existing, file)
			redirectionsMutex.Unlock()
		}
		return nil
	})

	// err = catalog.WriteMentionPages()
	// if err != nil {
	// 	log.Fatal("Failed to build mention pages:", err)

	// }
	// err = catalog.WritePopups()
	// if err != nil {
	// 	log.Fatal("Failed to build mention popup page:", err)

	// }
	// err = catalog.WriteSeriesPages()
	// if err != nil {
	// 	log.Fatal("Failed to build asset series pages:", err)
	// }

	// log.Printf("Write Files", time.Since(tStart))

	slices.SortFunc(completedAssets, rprCatalog.SortAssetsByDateAdded)

	// log.Printf("Sort Files", time.Since(tStart))

	// var latestFile *rprCatalog.Asset = nil
	// if len(completedAssets) > 0 {
	// 	latestFile = completedAssets[0]
	// }

	// blogPosts, err := blog.Write(catalog)
	// if err != nil {
	// 	log.Fatal("Failed to write blog:", err)
	// }

	// Stats

	statsPage, _ := utils.MakePage("stats")
	missing := map[*rprCatalog.Asset][]int{}
	prefixes := []string{
		"https://wiki.chadnet.org",
		"https://www.toxinless.com",
		"https://github.com/0x2447196/raypeatarchive",
		"https://expulsia.com",
	}
	for _, asset := range catalog.Assets {
		allFound := true
		results := []int{}
		for _, prefix := range prefixes {
			foundCount := 0
			for _, url := range asset.GetAllURLsUnescaped() {
				if strings.HasPrefix(url, prefix) {
					foundCount += 1
				}
			}
			results = append(results, foundCount)
			if foundCount == 0 {
				allFound = false
			}
		}
		if !allFound {
			missing[asset] = results
		}
	}
	component := Stats(prefixes, missing)
	component.Render(context.Background(), statsPage)

	// log.Printf("Stats", time.Since(tStart))

	// 🏠 Homepage

	indexPage, _ := utils.MakePage(".")
	component = Index(catalog)
	component.Render(context.Background(), indexPage)

	// log.Printf("Homepage", time.Since(tStart))

	// Ray Peat Page

	rayPeatPage, _ := utils.MakePage("ray-peat")
	component = RayPeatPage()
	component.Render(context.Background(), rayPeatPage)

	// log.Printf("Ray Peat Page", time.Since(tStart))

	// Catalog cache

	err = catalog.HttpCache.Write()
	if err != nil {
		log.Fatal("Failed to write HTTP cache:", err)
	}

	// log.Printf("Catalog Cache", time.Since(tStart))

	// JSON
	searchData := []SearchAsset{}
	for _, asset := range catalog.Assets {

		contributors := []SearchContributor{}
		for _, contributor := range asset.GetFilterableSpeakers() {
			contributors = append(contributors, SearchContributor{
				Name:   contributor.GetName(),
				Avatar: contributor.GetAvatarPath(),
			})
		}

		issues := []SearchIssue{}
		for _, issue := range asset.Issues {
			issues = append(issues, SearchIssue{
				Id:    issue.Id,
				Title: issue.Title,
				Url:   issue.Url,
			})
		}

		sections := []SearchSection{}
		for _, section := range asset.Sections {
			timecode := (func() *SearchTimecode {
				if section.Timecode == nil {
					return nil
				}

				return &SearchTimecode{
					Hours:   section.Timecode.Hours,
					Minutes: section.Timecode.Minutes,
					Seconds: section.Timecode.Seconds,
				}
			})()

			sections = append(sections, SearchSection{
				Title:    section.Title,
				Prefix:   section.Prefix,
				Level:    section.Level,
				ID:       section.ID,
				Timecode: timecode,
			})
		}

		searchData = append(searchData, SearchAsset{
			Path:         asset.UrlAbsPath,
			Title:        asset.FrontMatter.Source.Title,
			Series:       asset.FrontMatter.Source.Series,
			Kind:         asset.FrontMatter.Source.Kind,
			Date:         asset.Date,
			Contributors: contributors,
			Issues:       issues,
			Sections:     sections,
		})
	}

	b, err := json.Marshal(searchData)
	err = os.WriteFile(path.Join(global.BUILD_OUTPUT, "search.json"), b, 0755)
	if err != nil {
		log.Fatal("Failed to write search JSON:", err)
	}

	// log.Printf("JSON", time.Since(tStart))

	// 🏁 Done

	fmt.Printf("\nFinished in %v.\n", time.Since(start))

	// log.Printf("Done", time.Since(tStart))
}

type SearchContributor struct {
	Avatar string
	Name   string
}

type SearchAsset struct {
	Path         string
	Title        string
	Series       string
	Kind         string
	Date         string
	Contributors []SearchContributor
	Issues       []SearchIssue
	Sections     []SearchSection
}

type SearchIssue struct {
	Id    int
	Title string
	Url   string
}

type SearchSection struct {
	Prefix   []string
	Title    string
	Level    int
	ID       string
	Timecode *SearchTimecode
}

type SearchTimecode struct {
	Hours   int
	Minutes int
	Seconds int
}
