package check

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getRayPeatArchiveAssets() []*ScrapedAsset {
	res, err := http.Get("https://api.github.com/repos/0x2447196/raypeatarchive/git/trees/main")
	if err != nil {
		log.Fatal("Failed to HTTP GET git tree via GitHub API: ", err)
	}

	var mainTree GitTreeData
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read response body of GitHub API for git tree:", err)
	}
	json.Unmarshal(bodyBytes, &mainTree)

	transcriptsEndpoint := ""
	documentsEndpoint := ""

	for _, path := range mainTree.Tree {
		switch path.Path {
		case "transcripts":
			transcriptsEndpoint = path.Url
		case "documents":
			documentsEndpoint = path.Url
		default:
			// ignore
		}
	}

	if transcriptsEndpoint == "" {
		log.Fatal("Failed to find ./transcripts in git tree")
	}
	if documentsEndpoint == "" {
		log.Fatal("Failed to find ./documents in git tree")
	}

	res, err = http.Get(transcriptsEndpoint)
	if err != nil {
		log.Fatal("Failed to HTTP GET ./transcripts via GitHub API:", err)
	}
	bodyBytes, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read response body of GitHub API for ./transcripts git tree:", err)
	}
	var transcriptsTree GitTreeData
	json.Unmarshal(bodyBytes, &transcriptsTree)

	assets := []*ScrapedAsset{}

	for _, entry := range transcriptsTree.Tree {
		filename := entry.Path
		listingUrl := "https://github.com/0x2447196/raypeatarchive/tree/main/transcripts"
		assetUrl := listingUrl + "/" + filename
		displayText := strings.TrimSuffix(filename, filepath.Ext(filename))
		date, title, series := extractDetails(displayText)

		slug := slug.Make(title)
		title = cases.Title(language.English, cases.Compact).String(title)

		kind := "unknown"
		switch series {
		case "Generative Energy":
			kind = "video"
		default:
			kind = "audio"
		}

		path := fmt.Sprintf("asset/%v-%v.md", date, slug)

		assets = append(assets, &ScrapedAsset{
			LinkHref:    assetUrl,
			ListedAtUrl: listingUrl,
			DisplayText: displayText,
			ProposedAsset: &catalog.Asset{
				FrontMatter: catalog.AssetFrontMatter{
					Source: catalog.AssetFrontMatterSource{
						Title:  title,
						Kind:   kind,
						Series: series,
						Mirrors: []string{
							assetUrl,
						},
					},
				},
				Date: date,
				Slug: slug,
				Path: path,
			},
		})
	}

	res, err = http.Get(documentsEndpoint)
	if err != nil {
		log.Fatal("Failed to HTTP GET ./documents via GitHub API:", err)
	}
	bodyBytes, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read response body of GitHub API for ./documents git tree:", err)
	}
	var documentsTree GitTreeData
	json.Unmarshal(bodyBytes, &documentsTree)

	for _, entry := range documentsTree.Tree {
		dirName := entry.Path
		if entry.Type != "tree" {
			log.Fatal("Found non-directory item in ./documents: ", dirName)
		}
		endpoint := entry.Url
		res, err = http.Get(endpoint)
		if err != nil {
			log.Fatalf("Failed to HTTP GET ./documents/%v via GitHub API: %v\n", dirName, err)
		}
		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Failed to read response body of GitHub API for ./documents/%v git tree: %v", dirName, err)
		}
		var tree GitTreeData
		json.Unmarshal(bodyBytes, &tree)

		for _, e := range tree.Tree {
			filename := e.Path
			title := strings.TrimSuffix(filename, filepath.Ext(filename))
			series := "unknown"
			kind := "unknown"

			switch dirName {
			case "books":
				series = "Raymond Peat"
				kind = "book"
			case "newsletters":
				series = "Townsend Letter for Doctors & Patients" // just a quess
				kind = "newsletter"
			case "raypeat.com":
				series = "Ray Peat's Website"
				kind = "article"
			case "raypeatinsight.wordpress.com":
				series = "Ray Peat Insight"
				kind = "text" // text and article are ambiguous, yes. See #55
			default:
				log.Fatal("Unexpected directory in ./documents: ", dirName)
			}

			listingUrl := "https://github.com/0x2447196/raypeatarchive/tree/main/documents"
			assetUrl := fmt.Sprintf("%v/%v/%v", listingUrl, dirName, filename)
			slug := slug.Make(title)
			date := "0000-00-00"
			path := fmt.Sprintf("assets/%v-%v.md", date, slug)

			assets = append(assets, &ScrapedAsset{
				ListedAtUrl: listingUrl,
				DisplayText: filename,
				LinkHref:    assetUrl,
				ProposedAsset: &catalog.Asset{
					Date: date,
					Path: path,
					Slug: slug,
					FrontMatter: catalog.AssetFrontMatter{
						Source: catalog.AssetFrontMatterSource{
							Kind:   kind,
							Title:  title,
							Series: series,
							Mirrors: []string{
								assetUrl,
							},
						},
					},
				},
			})
		}
	}

	return assets
}

func extractDetails(filename string) (string, string, string) {
	// e.g. 01.17.22 Peat Ray [1198355269].vtt
	if len(filename) > 8 {
		t, err := time.Parse("01.02.06", filename[:len("01.02.14")+1]) // assume Jan 01, 2006
		if err == nil {
			date := t.Format("2006-01-02")
			series := "Patrick Timpone's One Radio Network"
			title := filename[len("01.02.06")+2:]
			return date, title, series
		}
	}

	// e.g. (2005-10) Ray Peat - Nervous System Protect & Restore [mdLHWFJI2y0].vtt
	t, err := time.Parse("(2006-01)", strings.Split(filename, " ")[0])
	if err == nil {
		date := t.Format("2006-01-00")
		title := strings.Join(strings.Split(filename, " ")[1:], " ")
		series := "unknown"
		return date, title, series
	}

	// e.g. eluv-080918-fats.vtt
	if strings.Contains(filename, "-") {
		words := strings.Split(filename, "-")
		t, err = time.Parse("060102", words[1])
		if err == nil {
			date := t.Format("2006-01-02")
			title := strings.Join(words[2:], " ")
			series := "unknown"
			switch words[0] {
			case "blp":
				series = "Butter Living Podcast"
			case "eluv":
				series = "Eluv"
			case "ewh":
				series = "East West Healing"
			case "garynull":
				series = "Gary Null"
			case "jf":
				series = "Jodellefit"
			case "kkvv":
				series = "Hope For Health"
			case "kmud":
				series = "Ask Your Herb Doctor"
			case "orn":
				series = "Patrick Timpone's One Radio Network"
			case "polsci":
				series = "Politics & Science"
			case "rainmaking":
				series = "Rainmaking Time"
			case "sourcenutritional":
				series = "Source Nutritional Show"
			case "voiceofamerica":
				series = "Sharon Kleyne Hour"
			case "wp":
				series = "World Puja"
			case "yohaf":
				series = "Your Own Health And Fitness"
			}
			return date, title, series
		}
	}

	if strings.HasPrefix(filename, "#") {
		date := "0000-00-00"
		title := filename
		series := "Generative Energy"
		return date, title, series
	}

	return "0000-00-00", "untitled", "unknown"
}

type GitTreeData struct {
	Tree []struct {
		Path string
		Mode string
		Type string
		Sha  string
		Size int
		Url  string
	}
}
