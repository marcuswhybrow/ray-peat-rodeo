package check

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
)

func getExpulsiaAssets() []*ScrapedAsset {
	baseUrl, err := url.Parse("https://expulsia.com/health/peat-index")
	if err != nil {
		log.Fatalf("Failed to parse base url: %v", err)
	}
	expectedHash := "f7f9d1ed3fe1e69b2f31e2da9d11e7aa8aabd3b659053c3e1b3bf2465c30e8b9"

	body, err := fetchUrl(baseUrl.String(), expectedHash)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to fetch URL '%v': %v", baseUrl.String(), err))
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Panicln("Failed to parse Chadnet's HTML.")
	}

	assets := []*ScrapedAsset{}

	doc.Find("ul li a").Each(func(i int, anchor *goquery.Selection) {
		href := anchor.AttrOr("href", "")
		hrefUrl, err := baseUrl.Parse(href)
		if err != nil {
			log.Fatalf("Failed to parse asset href '%v': %v", href, err)
		}
		prevH2Text := anchor.ParentsFiltered("div").First().PrevAllFiltered("h2").First().Text()

		kind := "unknown"
		title := "untitled"
		assetDate := "0000-00-00"
		series := "unknown"

		switch filepath.Ext(href) {
		case ".pdf":
			switch prevH2Text {
			case "Books":
				kind = "book"
				title = anchor.Text()[len("[PDF] "):]
				series = "Raymond Peat"
			default:
				title = anchor.Text()[len("[PDF] "):]
				kind = "newsletter"
				series = "Townsend Letter for Doctors & Patients" // probable guess
				words := strings.Split(anchor.Text(), " ")
				month := words[0]
				year := words[1]
				if strings.HasPrefix(month, "Q-") {
					switch month[2] {
					case '1':
						month = "March"
					case '2':
						month = "June"
					case '3':
						month = "September"
					case '4':
						month = "December"
					}
				}

				t, err := time.Parse("January 2006", fmt.Sprint(month, year))
				if err == nil {
					assetDate = fmt.Sprintf("%v-%v-%v", t.Year(), t.Month(), "00")
				}
			}
		default:
			kind = "article"
			title = anchor.Text()
			series = "Ray Peat's Website"
		}

		slug := cleanSlug(title)
		path := fmt.Sprintf("assets/%v-%v.md", assetDate, slug)

		assets = append(assets, &ScrapedAsset{
			ListedAtUrl: baseUrl.String(),
			DisplayText: anchor.Text(),
			LinkHref:    hrefUrl.String(),
			ProposedAsset: &catalog.Asset{
				FrontMatter: catalog.AssetFrontMatter{
					Source: catalog.AssetFrontMatterSource{
						Kind:    kind,
						Title:   title,
						Series:  series,
						Mirrors: []string{hrefUrl.String()},
					},
				},
				Path: path,
				Date: assetDate,
			},
		})
	})

	baseUrl, err = url.Parse("https://expulsia.com/health/interviews")
	if err != nil {
		log.Fatalf("Failed to parse base url: %v", err)
	}
	expectedHash = "02ce84f8146cf5027cc3d9efc7dbbd9e5935f5fc55fa1b32f1821f33b9b2230a"

	body, err = fetchUrl(baseUrl.String(), expectedHash)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to fetch URL '%v': %v", baseUrl.String(), err))
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Panicln("Failed to parse Chadnet's HTML.")
	}

	articlesAndPapers := doc.Find("ul li a")
	articlesAndPapers.Each(func(i int, anchor *goquery.Selection) {
		href := anchor.AttrOr("href", "")
		hrefUrl, err := baseUrl.Parse(href)
		if err != nil {
			log.Fatalf("Failed to parse asset href '%v': %v", href, err)
		}

		kind := "unknown"
		title := "untitled"
		assetDate := "0000-00-00"
		series := "unknown"

		anchorText := anchor.Text()

		hrefUrlPath, _ := strings.CutSuffix(hrefUrl.String(), "?"+hrefUrl.RawQuery)

		ext := filepath.Ext(hrefUrlPath)
		switch ext {
		case ".html":
			kind = "text"
			title = anchor.Text()
		case ".mp3":
			kind = "audio"
			colonIndex := strings.Index(anchorText, ":")
			if colonIndex >= 0 {
				series = anchorText[:colonIndex]
				title = anchorText[colonIndex+2:]
			} else {
				title = anchorText
			}
		case ".doc", ".docx", ".php", ".pdf":
			// All .docs are transcription's we can ignore
		default:
			if slices.Contains(
				[]string{
					"www.youtube.com",
					"youtube.com",
					"youtu.be",
				},
				hrefUrl.Hostname(),
			) {
				kind = "video"
				title = anchor.Text()
			} else {
				log.Fatalf("Unexpect asset link extension: ", ext)
			}
		}

		slug := cleanSlug(title)
		path := fmt.Sprintf("assets/%v-%v.md", assetDate, slug)

		assets = append(assets, &ScrapedAsset{
			ListedAtUrl: baseUrl.String(),
			DisplayText: anchor.Text(),
			LinkHref:    hrefUrl.String(),
			ProposedAsset: &catalog.Asset{
				FrontMatter: catalog.AssetFrontMatter{
					Source: catalog.AssetFrontMatterSource{
						Kind:    kind,
						Title:   title,
						Series:  series,
						Mirrors: []string{hrefUrl.String()},
					},
				},
				Path: path,
				Date: assetDate,
			},
		})
	})

	return assets
}
