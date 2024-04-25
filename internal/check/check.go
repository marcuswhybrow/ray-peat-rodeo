package check

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	rprCatalog "github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
)

func Check() error {
	err := spinner.New().TitleStyle(huh.ThemeBase().Form).
		Title("Checking Ray Peat Rodeo").
		Run()

	if err != nil {
		log.Fatal(err)
	}

	// 🗃️ Identify existing assets

	catalog := rprCatalog.NewCatalog("./assets")

	fmt.Printf("Found %v assets.\n", len(catalog.Assets))

	allGot := []*ScrapedAsset{}
	allMatchNeed := []*ScrapedAssetMatch{}
	allNeed := []*ScrapedAsset{}

	handle := func(got []*ScrapedAsset, matchNeed []*ScrapedAssetMatch, need []*ScrapedAsset) {
		printResults(got, matchNeed, need)
		allGot = append(allGot, got...)
		allMatchNeed = append(allMatchNeed, matchNeed...)
		allNeed = append(allNeed, allNeed...)
	}

	fmt.Println("\n[Checking https://wiki.chadnet.org/ray-peat]")
	handle(getProposals(catalog, getChadNetAssets()))

	fmt.Println("\n[Checking https://expulsia.com/health]")
	handle(getProposals(catalog, getExpulsiaAssets()))

	fmt.Println("\n[Checking https://github.com/0x2447196/raypeatarchive]")
	handle(getProposals(catalog, getRayPeatArchiveAssets()))

	fmt.Println("\n[Checking https://www.toxinless.com/peat/podcast.rss]")
	handle(getProposals(catalog, getToxinlessAssets()))

	fmt.Println("\n[Checking https://www.selftestable.com/ray-peat-stuff/sites]")
	handle(getProposals(catalog, getSelfTestableAssets()))

	fmt.Println("\n[Checking https://raypeatforum.com/community/forums/audio-interview-transcripts.73]")
	handle(getProposals(catalog, getRayPeatForumAssets()))

	// fmt.Println("\n[Checking https://www.functionalps.com/blog/tag/ray-peat]")

	var (
		choice string
	)

	for _, matched := range allMatchNeed {
		options := []huh.Option[string]{
			huh.NewOption("Modify match ./"+matched.Asset.Path, "modify-recommended").Selected(true),
			huh.NewOption("Create suggested asset name: "+matched.ScrapedAsset.ProposedAsset.Path, "create"),
			huh.NewOption("Create a new asset name (decided in next step)", "create-different"),
			huh.NewOption("Do nothing", "skip"),
		}

		for _, a := range catalog.Assets {
			if a.Path != matched.Asset.Path {
				options = append(options, huh.NewOption("Modify ./"+a.Path, "modify:"+a.Path))
			}
		}

		fmt.Printf("\n[%v]\n", matched.ScrapedAsset.DisplayText)
		fmt.Println("Listed:  " + matched.ScrapedAsset.ListedAtUrl)
		fmt.Println("Kind:    " + matched.Asset.FrontMatter.Source.Kind)
		fmt.Println("Series:  " + matched.Asset.FrontMatter.Source.Series)
		fmt.Println("Url:     https://raypeat.rodeo" + matched.Asset.UrlAbsPath)
		mirrors := matched.ScrapedAsset.ProposedAsset.FrontMatter.Source.Mirrors
		if len(mirrors) > 0 {
			fmt.Print("Mirrors: ")
			for i, m := range mirrors {
				if i > 0 {
					fmt.Print("\n         ")
				}
				fmt.Print(m)
			}
			fmt.Println()
		} else {
			fmt.Println("Mirrors: none")
		}

		fmt.Println("Found:   ./" + matched.ScrapedAsset.ProposedAsset.Path)
		fmt.Println("Match:   ./" + matched.Asset.Path)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Create Or Modify An Existing Asset?").
					Options(options...).
					Height(20).
					Value(&choice),
			),
		).WithTheme(huh.ThemeBase())
		form.Init()

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Choice:  ", choice)

	}

	return nil
}

type ScrapedAsset struct {
	ListedAtUrl   string
	DisplayText   string
	LinkHref      string
	ProposedAsset *rprCatalog.Asset
}

type ScrapedAssetMatch struct {
	Asset        *rprCatalog.Asset
	ScrapedAsset *ScrapedAsset
}
