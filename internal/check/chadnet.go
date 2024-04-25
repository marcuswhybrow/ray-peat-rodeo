package check

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
)

func getChadNetAssets() []*ScrapedAsset {
	res, err := http.Get("https://wiki.chadnet.org/ray-peat")
	if err != nil {
		log.Panicln("Failed to get Chadnet's Ray Peat page:", err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln("Failed to parse Chadnet's HTML.")
	}

	// Books

	bookLinks := doc.Find("#books + ul > li > a")
	if bookLinks.Length() == 0 {
		log.Panicln("Failed to find Chadnet's books.")
	}

	scrapedAssets := []*ScrapedAsset{}
	bookLinks.Each(func(i int, s *goquery.Selection) {
		scrapedAssets = append(scrapedAssets, chadnetAnchorToAsset(s, "book", "Raymond Peat"))
	})

	// Articles

	articleLinks := doc.Find("#articles + ul > li > a")
	if articleLinks.Length() == 0 {
		log.Panicln("Failed to find Chadnet's articles.")
	}

	articleLinks.Each(func(i int, s *goquery.Selection) {
		scrapedAssets = append(scrapedAssets, chadnetAnchorToAsset(s, "article", "Ray Peat's Website"))
	})

	// ✉️  Newsletters

	newsletterLinks := doc.Find("#newsletters + ul > li > a")
	if newsletterLinks.Length() == 0 {
		log.Panicln("Failed to find Chadnet's newsletters.")
	}

	newsletterLinks.Each(func(i int, s *goquery.Selection) {
		scrapedAssets = append(scrapedAssets, chadnetAnchorToAsset(s, "newsletter", "Townsend Letter for Doctors & Patients"))
	})

	// Miscellaneous

	miscLinks := doc.Find("#miscellaneous + ul > li > a")
	if miscLinks.Length() == 0 {
		log.Panicln("Failed to find Chadnet's miscellaneous assets.")
	}

	miscLinks.Each(func(i int, s *goquery.Selection) {
		scrapedAssets = append(scrapedAssets, chadnetAnchorToAsset(s, "unknown", "Raymond Peat"))
	})

	// Interviews

	interviewsRes, err := http.Get("https://wiki.chadnet.org/ray-peat-interviews")
	if err != nil {
		log.Panicln("Failed to get Chadnet's Ray Peat page:", err)
	}

	interviewsDoc, err := goquery.NewDocumentFromReader(interviewsRes.Body)
	if err != nil {
		log.Panicln("Failed to parse Chadnet's HTML.")
	}

	firstCategory := interviewsDoc.Find("body > h2")
	if firstCategory.Length() == 0 {
		log.Panicf("Failed to find Chadnet's interviews")
	}

	catagory := firstCategory.First()
	entries := firstCategory.Next().First()
	for true {
		entries.Find("li").Each(func(i int, s *goquery.Selection) {
			anchor := s.Find("a")
			anchorText := anchor.Text()
			liText := s.Text()
			href := anchor.AttrOr("href", "ewh-000000")
			title := "Untitled"
			series := "unknown"
			date := "0000-00-00"
			kind := "unknown"
			slug := "untitled"

			switch catagory.Text() {
			case "East West Healing":
				title = anchorText[len("East West: "):]
				series = "East West Healing"
				date, err = fromSixDate(href[len("ewh-"):])
				kind = "audio"
			case "Herb Doctors":
				title = anchorText[len("Herb Doctors: "):]
				series = "Ask Your Herb Doctor"
				date, err = fromSixDate(href[len("kmud-"):])
				kind = "audio"
			case "Politics & Science":
				title = anchorText[len("Politics & Science: "):]
				series = "Politics & Science"
				date, err = fromSixDate(href[len("polsci-"):])
				kind = "audio"
			case "One Radio Network (to be improved)":
				title = anchorText[len("14.01.01 "):]
				series = "One Radio Network"
				date, err = fromSixDotDate(anchorText)
				kind = "audio"
			case "Miscellaneous":
				colonIndex := strings.Index(anchorText, ":")
				if colonIndex >= 0 {
					title = anchorText[colonIndex+2:]
					series = anchorText[:colonIndex]
				} else {
					title = anchorText
				}
				kind = "audio"
				date, err = fromSixDate(href[strings.Index(href, "-")+1:])
			case "YouTube":
				title = anchorText
				parenIndex := strings.LastIndex(liText, "(")
				if parenIndex >= 0 {
					series = liText[parenIndex+1 : len(liText)-1]
				}
				kind = "video"
				date = "0000-00-00"

			default:
				title = anchorText
				series = catagory.Text()
			}

			slug = cleanSlug(title)

			if err != nil {
				log.Panicf("Failed to read date for interview '%v': %v", anchorText, err)
			}

			linkUrl, err := url.Parse(href)
			if err != nil {
				log.Panicf("Failed to parse link for asset '%v': %v", anchorText, err)
			}

			mirrors := []string{}
			if linkUrl.IsAbs() {
				if linkUrl.Hostname() == "wiki.chadnet.org" {
					mirrors = append(mirrors, "https://wiki.chadnet.org/"+linkUrl.Path)
					mirrors = append(mirrors, "https://wiki.chadnet.org/files/"+linkUrl.Path)
				} else {
					mirrors = append(mirrors, linkUrl.String())
				}
			} else {
				mirrors = append(mirrors, "https://wiki.chadnet.org/"+href)
				mirrors = append(mirrors, "https://wiki.chadnet.org/files/"+href)
			}

			scrapedAssets = append(scrapedAssets, &ScrapedAsset{
				ListedAtUrl: "https://wiki.chadnet.org/ray-peat-interviews",
				DisplayText: liText,
				LinkHref:    href,

				ProposedAsset: &catalog.Asset{
					FrontMatter: catalog.AssetFrontMatter{
						Source: catalog.AssetFrontMatterSource{
							Series:  series,
							Title:   title,
							Kind:    kind,
							Mirrors: mirrors,
						},
					},
					Date: date,
					Slug: slug,
					Path: fmt.Sprintf("assets/%v-%v.md", date, slug),
				},
			})
		})

		catagory = entries.Next()
		if catagory.Length() == 0 || !catagory.Is("h2") {
			break
		}
		entries = catagory.Next()
		if entries.Length() == 0 || !entries.Is("ul") {
			break
		}
	}

	return scrapedAssets
}

func chadnetAnchorToDetails(anchorText string) (string, string, string, string) {
	words := strings.Split(anchorText, " ")

	year := "0000"
	month := "00"
	day := "00"

	length := len(words)
	if length >= 1 {
		// 2024 ...
		_, err := strconv.Atoi(words[0])
		if err != nil {
			return year, month, day, strings.Join(words, " ")
		}
		year = words[0]
	}

	if length >= 2 {
		// 2024 07/08 ...
		if len(words[1]) == 5 && words[1][2] == '/' {
			month = words[1][3:5]
		} else if strings.ToLower(words[1]) == "q1" {
			month = "03"
		} else if strings.ToLower(words[1]) == "q2" {
			month = "06"
		} else if strings.ToLower(words[1]) == "q3" {
			month = "09"
		} else if strings.ToLower(words[1]) == "q4" {
			month = "12"
		} else {
			_, err := strconv.Atoi(words[1])
			if err != nil {
				return year, month, day, strings.Join(words[1:], " ")
			}
			month = words[1]
		}
	}

	if length >= 3 {
		return year, month, day, strings.Join(words[2:], " ")
	}

	return year, month, day, ""
}

func chadnetAnchorToAsset(s *goquery.Selection, kind string, series string) *ScrapedAsset {
	href, ok := s.Attr("href")
	if !ok {
		fmt.Println("Skipping asset with missing link:", s.Text())
		return nil
	}
	year, month, day, title := chadnetAnchorToDetails(s.Text())

	url := "https://wiki.chadnet.org/" + href
	if href[len(href)-4:] == ".pdf" {
		url = "https://wiki.chadnet.org/files/" + href
	}

	slug := cleanSlug(title)

	return &ScrapedAsset{
		ListedAtUrl: "https://wiki.chadnet.org/ray-peat",
		LinkHref:    s.AttrOr("href", ""),
		DisplayText: s.Text(),
		ProposedAsset: &catalog.Asset{
			FrontMatter: catalog.AssetFrontMatter{
				Source: catalog.AssetFrontMatterSource{
					Title: title,
					Mirrors: []string{
						url,
					},
					Series: series,
					Kind:   kind,
				},
			},
			Path:     fmt.Sprintf("assets/%v-%v-%v-%v.md", year, month, day, slug),
			Date:     fmt.Sprintf("%v-%v-%v", year, month, day),
			Slug:     slug,
			Markdown: []byte{},
		},
	}
}
