package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gernest/front"
	"github.com/mitchellh/mapstructure"
)

func main() {
	fmt.Println("[Checking chadnet.org]")

	// 🌐🎰 Download & Parse chadnet.org/ray-peat

	res, err := http.Get("https://wiki.chadnet.org/ray-peat")
	if err != nil {
		log.Panicln("Failed to get Chadnet's Ray Peat page:", err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln("Failed to parse Chadnet's HTML.")
	}

	// ✉️  Identify Newsletters

	selection := doc.Find("#newsletters + ul > li > a")
	if selection.Length() == 0 {
		log.Panicln("Failed to find Chadnet's newsletters.")
	}

	assets := []*asset{}
	selection.Each(func(i int, s *goquery.Selection) {
		asset := anchorToAsset(s, "newsletter", "Townsend Letter for Doctors & Patients")
		assets = append(assets, asset)
	})
	createAssets(existingUrls, assets, "newsletters")

	// 📰 Identify Articles

	assets = []*asset{}
	selection = doc.Find("#articles + ul > li > a")
	if selection.Length() == 0 {
		log.Panicln("Failed to find Chadnet's articles.")
	}
	selection.Each(func(i int, s *goquery.Selection) {
		asset := anchorToAsset(s, "article", "TO_BE_DETERMINED")
		assets = append(assets, asset)
	})
	createAssets(existingUrls, assets, "articles")
}

type asset struct {
	Filename string
	Year     string
	Month    string
	Day      string
	Slug     string
	Title    string
	Url      string
	Series   string
	Kind     string
}

func cleanSlug(title string) string {
	title = strings.ReplaceAll(title, "\"", "\\\"")

	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	slug = strings.ReplaceAll(slug, ",", "")
	slug = strings.ReplaceAll(slug, ".", "")
	slug = strings.ReplaceAll(slug, "(", "")
	slug = strings.ReplaceAll(slug, ")", "")
	slug = strings.ReplaceAll(slug, "?", "")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, ":", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\\", "")
	return slug
}

func anchorToDetails(anchorText string) (string, string, string, string) {
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

type FrontMatter struct {
	Source struct {
		Url     string
		Mirrors []string
	}
}

func anchorToAsset(s *goquery.Selection, kind string, series string) *asset {
	href, ok := s.Attr("href")
	if !ok {
		fmt.Println("Skipping asset with missing link:", s.Text())
		return nil
	}
	year, month, day, title := anchorToDetails(s.Text())

	url := "https://wiki.chadnet.org/" + href
	if href[len(href)-4:] == ".pdf" {
		url = "https://wiki.chadnet.org/files/" + href
	}

	slug := cleanSlug(title)

	return &asset{
		Filename: fmt.Sprintf("%v-%v-%v-%v.md", year, month, day, slug),
		Year:     year,
		Month:    month,
		Day:      day,
		Title:    title,
		Url:      url,
		Series:   series,
		Kind:     kind,
	}
}

func createAssets(existingUrls []string, assets []*asset, kind string) {
	need := []*asset{}
	got := []*asset{}

	for _, asset := range assets {
		if slices.Contains(existingUrls, asset.Url) {
			got = append(got, asset)
		} else {
			need = append(need, asset)
		}
	}

	if len(need) == 0 {
		fmt.Printf("✅ All %v %v on chadnet.org are on Ray Peat Rodeo.\n", len(assets), kind)
		return
	}

	fmt.Printf("❌ Found %v %v on chadnet.org, of which %v were not on Ray Peat Rodeo.\n", len(assets), kind, len(need))

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nCreate %v automatically [y/N]: \n", kind)
	answer, err := reader.ReadString('\n')
	if err != nil {
		log.Panic("Failed to read stdin response to question: ", err)
	}
	if strings.ToLower(answer) == "y" {
		for _, asset := range need {
			path := "./assets/" + asset.Filename
			body := "---\n" +
				"source:\n" +
				"  url: " + asset.Url + "\n" +
				"  title: \"" + asset.Title + "\"\n" +
				"  kind: " + asset.Kind + "\n" +
				"  series: " + asset.Series + "\n" +
				"---\n"
			os.WriteFile(path, []byte(body), 0644)
			fmt.Println("Wrote ", asset.Url)
		}
	}
}
