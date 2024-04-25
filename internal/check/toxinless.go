package check

import (
	"fmt"
	"log"
	"strings"

	"github.com/gosimple/slug"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
	"github.com/mmcdole/gofeed"
)

func getToxinlessAssets() []*ScrapedAsset {
	feed, err := gofeed.NewParser().ParseURL("https://www.toxinless.com/peat/podcast.rss")
	if err != nil {
		log.Fatal("Failed to parse RSS feed for toxinless.com:", err)
	}

	assets := []*ScrapedAsset{}

	for _, item := range feed.Items {
		rssTitle := item.Title

		series := "unknown"
		title := "untitled"
		kind := "audio"

		colonIndex := strings.Index(rssTitle, ":")
		if colonIndex >= 0 {
			rssSeries := rssTitle[:colonIndex]
			switch rssSeries {
			case "Ask the Herb Doctor":
				series = "Ask Your Herb Doctor"
			case "Politics & Science", "Butter Living Podcast", "Jodellefit",
				"World Puja", "Source Nutritional Show", "Hope For Health", "Gary Null":
				series = rssSeries
			case "[RELEASED 2019] Politics & Science":
				series = "Politics & Science"
			case "One Radio Network":
				series = "Patrick Timpone's One Radio Network"
			case "It's Rainmaking Time":
				series = "Rainmaking Time"
			case "ELUV":
				series = "Eluv"
			case "Voice of America":
				series = "Sharon Kleyne Hour"
			case "EastWest Healing":
				series = "East West Healing"
			default:
				log.Fatal("Unexpected rss item series:", rssSeries)
			}
			title = rssTitle[colonIndex+2:]
		} else {
			title = rssTitle
		}

		assetSlug := slug.Make(title)
		date := item.PublishedParsed.Format("2006-01-02")
		path := fmt.Sprintf("assets/%v-%v.md", date, assetSlug)

		assets = append(assets, &ScrapedAsset{
			ProposedAsset: &catalog.Asset{
				Date: date,
				Slug: assetSlug,
				Path: path,
				FrontMatter: catalog.AssetFrontMatter{
					Source: catalog.AssetFrontMatterSource{
						Title:  title,
						Kind:   kind,
						Series: series,
						Mirrors: []string{
							item.Link,
						},
					},
				},
			},
		})
	}

	return assets
}
