package catalog

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path"
	"slices"
	"strings"
	"sync"
	"text/template"

	// "github.com/gosimple/slug"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/global"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/extension"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/utils"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	gmExtension "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Mentions = []*ast.Mention
type ByAsset[T any] map[*Asset]T
type ByPart[T any] map[ast.MentionablePart]T
type ByMentionable[T any] map[ast.Mentionable]T

type Catalog struct {
	MarkdownParser    goldmark.Markdown
	HttpCache         *cache.HTTPCache
	AvatarPaths       *AvatarPaths
	ByMentionable     ByMentionable[ByAsset[Mentions]]
	ByMentionablePart ByPart[ByPart[ByAsset[Mentions]]]
	BySeries          map[string][]*Asset
	Assets            []*Asset
	Mutex             sync.RWMutex
}

func NewCatalog(assetsPath string) *Catalog {

	// 📶 Get HTTP Cache from file

	cacheData, err := cache.DataFromYAMLFile(global.CACHE_PATH)
	if err != nil {
		log.Panicf("Failed to read cache file '%v': %v", global.CACHE_PATH, err)
	}

	httpCache := cache.NewHTTPCache(cacheData)

	// 👨👱 Speaker Avatars

	avatarPaths := NewAvatarPaths()

	// 🖹🧩 Markdown + Extensions

	markdownParser := goldmark.New(
		goldmark.WithRendererOptions(html.WithUnsafe()),
		goldmark.WithExtensions(
			meta.New(meta.WithStoresInDocument()),
			gmExtension.Typographer,
			extension.Mentions,
			extension.Timecodes,
			extension.Speakers,
			extension.Sidenotes,
			extension.GitHubIssues,
		),
	)

	filePaths := utils.Files(assetsPath, ".", func(filePath string) (*string, error) {
		base := path.Base(filePath)
		baseLower := strings.ToLower(base)

		if baseLower == "readme.md" {
			return nil, nil
		}

		outPath := path.Join(assetsPath, filePath)
		return &outPath, nil
	})

	catalog := &Catalog{
		MarkdownParser:    markdownParser,
		HttpCache:         httpCache,
		AvatarPaths:       avatarPaths,
		ByMentionable:     ByMentionable[ByAsset[Mentions]]{},
		ByMentionablePart: ByPart[ByPart[ByAsset[Mentions]]]{},
		BySeries:          map[string][]*Asset{},
		Assets:            []*Asset{},
	}

	utils.Parallel(filePaths, func(filePath string) error {
		_, err := catalog.NewAsset(filePath)
		if err != nil {
			return fmt.Errorf("Failed to retrieve file '%v': %v", filePath, err)
		}
		return nil
	})

	slices.SortFunc(catalog.Assets, SortAssetsByDate)

	return catalog
}

func (c *Catalog) NewAsset(filePath string) (*Asset, error) {
	asset, err := NewAsset(filePath, c.MarkdownParser, c.HttpCache, c.AvatarPaths)
	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate file: %v", err)
	}

	// Catalog is designed to handle asynchronous workloads
	// This requires locking and unlocking using the Mutex technique
	c.Mutex.Lock()

	c.Assets = append(c.Assets, asset)
	c.BySeries[asset.GetSeriesSlug()] = append(c.BySeries[asset.GetSeriesSlug()], asset)

	for mentionable, mentions := range asset.Mentionables {
		for existingMentionable, existingByFile := range c.ByMentionable {
			if mentionable.IsDuplicate(existingMentionable) {
				if mentionable.IsMoreComplex(existingMentionable) {
					suspect := anyValue(existingByFile)[0]
					suggestion := mentions[0]
					collisionPanic(suspect, suggestion)
				} else {
					suspect := mentions[0]
					suggestion := anyValue(existingByFile)[0]
					collisionPanic(suspect, suggestion)
				}
			}
		}

		if c.ByMentionable[mentionable] == nil {
			c.ByMentionable[mentionable] = ByAsset[Mentions]{}
		}
		c.ByMentionable[mentionable][asset] = mentions

		if c.ByMentionablePart[mentionable.Primary] == nil {
			c.ByMentionablePart[mentionable.Primary] = ByPart[ByAsset[Mentions]]{}
		}
		if c.ByMentionablePart[mentionable.Primary][mentionable.Secondary] == nil {
			c.ByMentionablePart[mentionable.Primary][mentionable.Secondary] = ByAsset[Mentions]{}
		}
		c.ByMentionablePart[mentionable.Primary][mentionable.Secondary][asset] = mentions
	}

	c.Mutex.Unlock()

	return asset, nil
}

func (c *Catalog) WriteMentionPages() error {
	for primary, mentionsByFileBySecondary := range c.ByMentionablePart {
		primaries := mentionsByFileBySecondary[ast.EmptyMentionablePart]
		delete(mentionsByFileBySecondary, ast.EmptyMentionablePart)

		file, _ := utils.MakePage(primary.ID())
		component := MentionPage(primary, primaries, mentionsByFileBySecondary, c.HttpCache)
		err := component.Render(context.Background(), file)
		if err != nil {
			return fmt.Errorf("Failed to render template: %v", err)
		}
		//log.Printf("Wrote %v", outPath)
	}
	return nil
}

func (c *Catalog) WritePopups() error {
	for mentionable, mentionsByFile := range c.ByMentionable {
		location, _ := strings.CutPrefix(mentionable.PopupPermalink(), "/")

		otherMentionables := []ast.Mentionable{}
		for m := range c.ByMentionable {
			samePrimary := m.Primary == mentionable.Primary
			identical := m.Secondary == mentionable.Secondary
			if samePrimary && !identical {
				otherMentionables = append(otherMentionables, m)
			}
		}

		f, _ := utils.MakePage(location)
		component := MentionablePopup(mentionable, mentionsByFile, otherMentionables)
		err := component.Render(context.Background(), f)
		if err != nil {
			return fmt.Errorf("Failed to render template: %v", err)
		}
		//log.Printf("Wrote %v", outPath)
	}

	return nil
}

func (c *Catalog) WriteSeriesPages() error {
	for series, assets := range c.BySeries {
		f, _ := utils.MakePage(series)
		component := SeriesPage(assets[0].FrontMatter.Source.Series, assets)
		err := component.Render(context.Background(), f)
		if err != nil {
			return fmt.Errorf("Failed to render template for asset series page '%v': %v", series, err)
		}
	}
	return nil
}

func (c *Catalog) GetCompletedAssets() []*Asset {
	var assets []*Asset
	for _, asset := range c.Assets {
		if asset.IsComplete() {
			assets = append(assets, asset)
		}
	}
	return assets
}

func collisionPanic(suspect, suggestion *ast.Mention) {
	t, err := template.New("AmbiguousMentionable").Parse(`
Found mention collision for...
  {{ .SuspectTag }}
{{ .Suspect.Source.Row }}:{{ .Suspect.Source.Col }} {{ .Suspect.File.GetPath }}

It has the same logical name as...
  {{ .SuggestionTag }}
{{ .Suggestion.Source.Row }}:{{ .Suggestion.Source.Col }} {{ .Suggestion.File.GetPath }}

Two mentions collide if they have different prefixes but identical conclusions.
For example, these three have the same conclusion "P. G. Wodehouse":

  [[P. G. Wodehouse]]
  [[Wodehouse, P. G.]]
  [[G. Wodehouse, P.]]

If the collision indeed refers to two different entities, consider making one 
more specific. For example:

  [[Wodehouse (1881), P. G.]]

`)
	if err != nil {
		panic("Failed to render collision panic template.")
	}

	var b bytes.Buffer
	t.Execute(&b, map[string]interface{}{
		"SuspectTag":    string(suspect.Source.Segment.Value(suspect.Asset.GetMarkdown())),
		"Suspect":       suspect,
		"SuggestionTag": string(suggestion.Source.Segment.Value(suggestion.Asset.GetMarkdown())),
		"Suggestion":    suggestion,
	})
	panic(b.String())

}

func anyValue[Key comparable, Val any](m map[Key]Val) Val {
	for _, v := range m {
		return v
	}
	panic("Failed to find any entries in map")
}

func (c *Catalog) MatchAsset(a *Asset) (*Asset, bool) {
	perfectMatch := true
	partialMatch := false

	var (
		dateMatch *Asset
		kindMatch *Asset
		slugMatch *Asset
	)

	for _, b := range c.Assets {
		// Compare URLs & Mirrors
		if a.FrontMatter.Source.Url == b.FrontMatter.Source.Url {
			return b, perfectMatch
		}
		if slices.Contains(b.FrontMatter.Source.Mirrors, a.FrontMatter.Source.Url) {
			return b, perfectMatch
		}
		for _, m := range a.FrontMatter.Source.Mirrors {
			if slices.Contains(b.FrontMatter.Source.Mirrors, m) {
				return b, perfectMatch
			}
		}

		// Date & source kind
		if a.Date != "0000-00-00" && b.Date == a.Date {
			if b.FrontMatter.Source.Kind == a.FrontMatter.Source.Kind {
				if b.Slug == a.Slug {
					slugMatch = b
				} else {
					kindMatch = b
				}
			} else {
				dateMatch = b
			}
		}
	}

	if slugMatch != nil {
		return slugMatch, partialMatch
	} else if kindMatch != nil {
		return kindMatch, partialMatch
	} else if dateMatch != nil {
		return dateMatch, partialMatch
	}

	return nil, perfectMatch
}
