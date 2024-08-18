package check

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	rprCatalog "github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
)

// Checks Ray Peat Rodeo has a copy of all Ray Peat assets listed on specific
// other web pages.
//
// Each page's HTML is fetched over the network, and relevant anchor tags are
// extracted and their "href" attributes added to a running list. Each href in
// that list is sought in Ray Peat Rodeo's catalog of assets: if the href is
// found to be an existing asset's source.url or in it's list of source.mirrors
// then the href is considered "known". The remainder are considered "missing",
// and printed to stdout for the user to rectify.
func Check() {
	catalog := rprCatalog.NewCatalog("./assets")

	knownExternalAssets := []string{}
	for _, asset := range catalog.Assets {
		knownExternalAssets = append(knownExternalAssets, asset.FrontMatter.Source.Url)
		knownExternalAssets = append(knownExternalAssets, asset.FrontMatter.Source.Mirrors...)
	}

	slices.Sort(knownExternalAssets)
	knownExternalAssets = slices.Compact(knownExternalAssets)

	unescapedKnownExternalAssets := []string{}
	for _, knownHref := range knownExternalAssets {
		unescaped, err := url.QueryUnescape(knownHref)
		if err != nil {
			log.Panicf("Failed to unescape href '%v': %v", knownHref, err)
		}
		unescapedKnownExternalAssets = append(unescapedKnownExternalAssets, unescaped)
	}

	fmt.Printf("Ray Peat Rodeo has %v assets distilled from %v locations.\n", len(catalog.Assets), len(knownExternalAssets))

	externalAssets := []string{}
	missingExternalAssets := []string{}
	check := func(hrefs []string) {
		for _, href := range hrefs {
			externalAssets = append(externalAssets, href)
			unescaped, err := url.QueryUnescape(href)
			if err != nil {
				log.Fatalf("Failed to unescape href '%v': %v", href, err)
			}

			if !slices.Contains(unescapedKnownExternalAssets, unescaped) {
				missingExternalAssets = append(missingExternalAssets, href)
			}
		}
	}

	check(scrapeLinks("http://raypeat.com/articles", ".posted a", []string{
		"http://raypeat.com/articles/articles/spanish-alzheimers.shtml",
		"http://raypeat.com/articles/articles/italian-salt.shtml",
	}))
	check(scrapeLinks("https://wiki.chadnet.org/ray-peat", "ul li a", []string{
		"https://wiki.chadnet.org/ray-peat-interviews",
	}))
	check(scrapeLinks("https://wiki.chadnet.org/ray-peat-interviews", "ul li a", []string{
		"https://bioenergetic.life/",
		"https://wiki.chadnet.org/polsci-080918-reductionist-science.mp3", // a clip, not a full interview
	}))
	check(scrapeLinks("https://www.westernbotanicalmedicine.com/pages/ask-your-herb-doctor", "table a", []string{
		"http://askyourherbdoctor.com/audio/Moles%20&%20Melanoma%20Clips%20All.mp3",
	}))
	check(scrapeLinks("https://askyourherbdoctor.com/media.html", "table a", []string{
		"https://askyourherbdoctor.com/audio/Moles%20&%20Melanoma%20Clips%20All.mp3",
	}))
	check(scrapeLinks("https://expulsia.com/health/interviews", "ul li a", []string{
		"https://expulsia.com/health/interviews/thyroidinterview.html",
		"https://expulsia.com/health/interviews/7-2012.html",
		"https://expulsia.com/health/interviews/9-2012.html",
		"https://expulsia.com/health/interviews/organizingthepanic.html",
		"https://expulsia.com/health/interviews/western.html",
		"https://expulsia.com/health/interviews/thyroid.html",
		"https://expulsia.com/health/interviews/Negation.html",
		"https://expulsia.com/health/interviews/onculture.html",
		"https://www.toxinless.com/polsci-080918-reductionist-science.mp3", // clip from longer episode, which I cannot as of yet identify
	}))
	check(scrapeLinks("https://expulsia.com/health/peat-index", "ul li a", nil))

	rayPeatArchive := []string{}
	rayPeatArchive = append(rayPeatArchive, scrapeLinks("https://github.com/0x2447196/raypeatarchive/tree/main/transcripts", "a.Link--primary", []string{
		"https://github.com/0x2447196/raypeatarchive/blob/main/transcripts/oxidation-pufa-roddy-peat.vtt",          // short clip
		"https://github.com/0x2447196/raypeatarchive/blob/main/transcripts/polsci-080918-reductionist-science.vtt", // short clip
		"https://github.com/0x2447196/raypeatarchive/blob/main/transcripts/.!82696!.DS_Store",
		"https://github.com/0x2447196/raypeatarchive/blob/main/transcripts/.DS_Store",
	})...)
	rayPeatArchive = append(rayPeatArchive, scrapeLinks("https://github.com/0x2447196/raypeatarchive/tree/main/documents/books", "a", nil)...)
	rayPeatArchive = append(rayPeatArchive, scrapeLinks("https://github.com/0x2447196/raypeatarchive/tree/main/documents/newsletters", "a", nil)...)
	rayPeatArchive = append(rayPeatArchive, scrapeLinks("https://github.com/0x2447196/raypeatarchive/tree/main/documents/raypeat.com", "a", []string{
		"https://github.com/0x2447196/raypeatarchive/blob/main/documents/raypeat.com/italian-salt.md",
		"https://github.com/0x2447196/raypeatarchive/blob/main/documents/raypeat.com/s.sh",
	})...)
	rayPeatArchive = append(rayPeatArchive, scrapeLinks("https://github.com/0x2447196/raypeatarchive/tree/main/documents/raypeatinsight.wordpress.com", "a", nil)...)

	check(hasPrefix(rayPeatArchive, "https://github.com/0x2447196/raypeatarchive/blob/"))

	check(scrapeLinks("https://www.selftestable.com/ray-peat-stuff/sites", "body > div.container > ul li a", []string{
		"https://www.biochemnordic.com/",
		"https://www.facebook.com/groups/134080236950737/",
		"https://www.facebook.com/groups/1490515417882745/",
		"https://www.facebook.com/groups/1551028545157939/",
		"https://www.facebook.com/groups/1671185569825129/",
		"https://www.facebook.com/groups/218721281534328/",
		"https://www.facebook.com/groups/252581431550065/",
		"https://www.facebook.com/groups/332104186934942/",
		"https://www.facebook.com/groups/417987225060250/",
		"https://www.facebook.com/groups/biochemnordic/",
		"https://www.pureenergypdx.com/",
		"https://www.raypeatforum.com/forum/viewtopic.php?f=68&t=1035",
		"https://www.toxinless.com/peat/",
		"https://www.toxinless.com/peat/search",
		"https://www.absolutelypure.com",
		"https://www.alexfergus.com/blog/",
		"https://www.amazon.com/Something-Sweet-Youll-Feel-Better/dp/B08KHPHH16/",
		"http://rayslight.wix.com/home",
		"http://seanbissell.com/blog/",
		"http://slimbirdy.com/",
		"http://valtsus.blogspot.fi/",
		"http://vitalityincorporated.com/?cat=10",
		"http://180degreehealth.com/",
		"http://50kzone.blogspot.com.br/search/label/name%20-%20Ray%20Peat%20Phd",
		"http://blog.arkofwellness.com/",
		"http://butternutrition.com/",
		"http://cowseatgrass.org/",
		"http://doctorsaredangerous.com/articles/dont_be_conned_by_the_resveratrol_scam.htm",
		"http://fatiguerecovery.co",
		"http://haidut.me/",
		"http://katedeering.com/blog",
		"http://losingcreekfarm.blogspot.com.br/",
		"http://onibasu.com/cgi-bin/search.cgi?query=ray+peat&submit=Go%21&idxname=cl&sort=score&max=20",
		"http://web.archive.org/web/20130508011036/http://www.systemicvitality.com/",
		"http://web.archive.org/web/20131213021923/http://www.paleogo.com/",
		"http://web.archive.org/web/20140619164428/http://www.workoutmaster.com/?tag=ray-peat",
		"https://www.toxinless.com/polsci-080918-reductionist-science.mp3", // 5 minute clip
		"http://www.danielstrassmann.de/",
		"http://www.dannyroddy.com/",
		"https://web.archive.org/web/20140105062345/http://vvfitness.wordpress.com/",
		"https://web.archive.org/web/20140620035830/http://www.roguewellness.com/",
		"https://web.archive.org/web/20150918164408/http://www.thenutritionwhisperer.com/blog1/",
		"https://web.archive.org/web/20160803113616/http://litalee.com/shopdisplayproducts.asp?id=30",
		"https://web.archive.org/web/20160820185955/http://www.andrewkimblog.com:80/",
		"https://web.archive.org/web/20161030104416/https://raypeatinsight.com/",
		"https://web.archive.org/web/20180808223029/http://peatarian.com/",
		"https://web.archive.org/web/20180829115504/http://www.visionandacceptance.com/",
		"https://www.youtube.com/user/joshrubineastwest/videos",
		"http://www.functionalps.com/blog/",
		"http://www.functionalps.com/blog/2012/12/07/collection-of-ray-peat-quote-blogs-by-fps/",
		"http://www.generativeenergy.com/",
		"http://www.kinesyne.com/blogue/",
		"http://www.ncbi.nlm.nih.gov/pubmed/?term=%22Peat+R%22[Author]",
		"http://www.nutritionbynature.com.au/",
		"http://www.perceivethinkact.com/",
		"http://www.raypeatforum.de/",
		"http://www.raypeatforums.org/",
		"http://www.resonantfm.com",
		"http://www.scottschlegel.net/category/health/",
		"http://www.thenutritioncoach.com.au/blog/",
		"http://www.thyroid-info.com/articles/ray-peat.htm",
		"http://www.tidesoflife.com/essential.htm",
		"http://www.yourownhealthandfitness.org/?page_id=483",
		"https://web.archive.org/web/20171213104513/http://joeylott.com/nutrition/guides/ray-peat-guide/",
		"https://web.archive.org/web/20171228030211/http://peatarian.com/peatexchanges",
		"https://jdperryhealth.tumblr.com/tagged/ray-peat",
		"https://l-i-g-h-t.com/",
		"https://peatarianreviews.blogspot.com/",
		"https://pinterest.com/brittany8832/foods-and-nutrition/",
		"https://pinterest.com/shaadoe/peatatarian/",
		"https://raypeatforum.com/community/",
		"https://raypeatforum.com/community/threads/stress-and-water.1261/",
		"https://raypeatforum.com/wiki/index.php/Ray_Peat_Email_Exchanges",
		"https://theurazoflazhealthblog.blogspot.com/",
		"https://vashinvetala.wordpress.com/",
		"https://co2factor.blogspot.com/",
		"https://fosteryourhealth.wordpress.com/",
		"https://github.com/Ray-Peat/interview/wiki/_pages",
		"https://gregorytaper.com/blog/",
		"https://groups.yahoo.com/neo/groups/AV-Skeptics/info",
		"https://healthandhappinezz.blogspot.com/2013/04/ray-peat.html",
		"https://ichooseicecream.wordpress.com/",
		"https://jayfeldmanwellness.com/articles/",
		"http://www.naturodoc.com/library/hormones/estrogen_pollution.htm",
		"http://www.naturodoc.com/library/nutrition/coconut_oil.htm",
		"https://cfsmethylation.blogspot.com/2014/05/sweet-sleep.html",
		"https://cfsmethylation.blogspot.com/2014/06/sweet-sleep-part-2.html",
		"https://archive.org/details/@fosteryourhealth",
		"https://astrotas.wordpress.com/bblog/",
		"https://bloqdnotas.blogspot.com/",
	}))

	// raypeatforum.com anti-bot proection makes conventional scraping impossible
	// functionalps.com is not a list of Ray Peat assets

	numMissing := len(missingExternalAssets)
	numAssets := len(externalAssets)
	fmt.Printf("A scan of other known websites that list Ray Peat assets finds %v locations of which %v are not known to Ray Peat Rodeo.\n", numAssets, numMissing)
	if numMissing > 0 {
		fmt.Println("Consider adding the following URLs to an asset's source.url or source.mirrors frontmatter to make this complaint go away:")
		for _, href := range missingExternalAssets {
			fmt.Println(href)
		}
	}
}

func scrapeLinks(targetUrl string, selector string, blacklist []string) []string {
	baseUrl, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatalf("Invalid URL for '%v': %v", targetUrl, err)
	}

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Panicln("Invalid GET request for ", targetUrl, ": ", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Panicln("Failed to GET URL", targetUrl, ": ", err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln("Failed to parse HTML for URL", targetUrl, ": ", err)
	}

	// allHTML, err := doc.Html()
	// fmt.Println(allHTML)

	links := doc.Find(selector)
	if links.Length() == 0 {
		log.Panicln("Failed to find any links in HTML for URL", targetUrl)
	}

	hrefs := []string{}
	links.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		absUrl, err := baseUrl.Parse(href)
		if err != nil {
			log.Printf("Ignoring invalid href '%v' for url '%v': %v", href, targetUrl, err)
			return
		}

		absUrlStr := absUrl.String()
		if !slices.Contains(blacklist, absUrlStr) {
			hrefs = append(hrefs, absUrlStr)
		}
	})

	slices.Sort(hrefs)
	return slices.Compact(hrefs)
}

func hasPrefix(hrefs []string, prefix string) []string {
	results := []string{}
	for _, h := range hrefs {
		unescaped, err := url.QueryUnescape(h)
		if err != nil {
			log.Fatalf("Failed to unescape href '%v': %v", h, err)
		}
		if strings.HasPrefix(unescaped, prefix) {
			results = append(results, h)
		}
	}
	return results
}
