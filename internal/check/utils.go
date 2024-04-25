package check

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
)

func printResults(got []*ScrapedAsset, matchNeed []*ScrapedAssetMatch, need []*ScrapedAsset) {
	needed := len(matchNeed) + len(need)
	total := needed + len(got)
	fmt.Printf("Found %v assets of which %v are unreferenced, %v of which may match existing assets.\n", total, needed, len(matchNeed))
}

func getProposals(cat *catalog.Catalog, assets []*ScrapedAsset) ([]*ScrapedAsset, []*ScrapedAssetMatch, []*ScrapedAsset) {
	got := []*ScrapedAsset{}
	matchNeed := []*ScrapedAssetMatch{}
	need := []*ScrapedAsset{}

	for _, a := range assets {
		match, perfectMatch := cat.MatchAsset(a.ProposedAsset)
		if match != nil {
			if perfectMatch {
				got = append(got, a)
			} else {
				matchNeed = append(matchNeed, &ScrapedAssetMatch{
					Asset:        match,
					ScrapedAsset: a,
				})
			}
		} else {
			need = append(need, a)
		}
	}

	return got, matchNeed, need
}

// Converts "YYMMDD" into "YYYY-MM-DD"
func fromSixDate(sixDate string) (string, error) {
	year, err := fromTwoYear(sixDate[0:2])
	month := sixDate[2:4]
	day := sixDate[4:6]

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%v-%v-%v", year, month, day), nil
}

// Converts "YY.MM.DD" into "YYYY-MM-DD"
func fromSixDotDate(sixDotDate string) (string, error) {
	year, err := fromTwoYear(sixDotDate[0:2])
	month := sixDotDate[3:5]
	day := sixDotDate[6:8]

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%v-%v-%v", year, month, day), nil
}

func fromTwoYear(twoYear string) (string, error) {
	yearInt, err := strconv.Atoi(twoYear)
	if err != nil {
		return "", err
	}
	if 2000+yearInt > time.Now().Year() {
		return "19" + twoYear, nil
	} else {
		return "20" + twoYear, nil
	}
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

func fetchUrl(url string, sha256Hash string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to HTTP GET '%v': %v", url, err))
	}

	body := new(strings.Builder)
	io.Copy(body, res.Body)
	bodyString := body.String()

	h := sha256.New()
	h.Write([]byte(bodyString))
	pageHash := fmt.Sprintf("%x", h.Sum(nil))

	if pageHash != sha256Hash {
		return "", errors.New(
			"HTTP response did not match given hash. Perhaps the page has changed.\n" +
				"  Expected Hash: '" + sha256Hash + "'\n" +
				"  Actual Hash:   '" + pageHash + "'\n",
		)
	}
	return bodyString, nil
}
