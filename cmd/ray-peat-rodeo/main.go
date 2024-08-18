package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/blog"
	rprCatalog "github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/check"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/global"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/utils"
)

func main() {
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

	if err := os.MkdirAll(global.BUILD_OUTPUT, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// err := os.RemoveAll(OUTPUT)
	// if err != nil {
	// 	log.Panicf("Failed to clean output directory: %v", err)
	// }

	// 🗃 Catalog

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

	// When an asset filename changes, it's URL changes.
	// It's nice to redirect old URL's to the new ones.
	// N.B. this data is currently collected, but not acted upon
	redirections := map[string][]*rprCatalog.Asset{}
	redirectionsMutex := sync.RWMutex{}

	utils.Parallel(catalog.Assets, func(file *rprCatalog.Asset) error {
		err := file.WriteHtml()
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

	err = catalog.WriteMentionPages()
	if err != nil {
		log.Fatal("Failed to build mention pages:", err)

	}
	err = catalog.WritePopups()
	if err != nil {
		log.Fatal("Failed to build mention popup page:", err)

	}
	err = catalog.WriteSeriesPages()
	if err != nil {
		log.Fatal("Failed to build asset series pages:", err)
	}

	slices.SortFunc(completedAssets, rprCatalog.SortAssetsByDateAdded)

	progress := float32(len(completedAssets)) / float32(len(catalog.Assets))

	var latestFile *rprCatalog.Asset = nil
	if len(completedAssets) > 0 {
		latestFile = completedAssets[0]
	}

	blogPosts, err := blog.Write(catalog)
	if err != nil {
		log.Fatal("Failed to write blog:", err)
	}

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

	// 🏠 Homepage

	indexPage, _ := utils.MakePage(".")
	component = Index(catalog.Assets, latestFile, progress, blogPosts[0])
	component.Render(context.Background(), indexPage)

	err = catalog.HttpCache.Write()
	if err != nil {
		log.Fatal("Failed to write HTTP cache:", err)
	}

	// 🏁 Done

	fmt.Printf("\nFinished in %v.\n", time.Since(start))
}
