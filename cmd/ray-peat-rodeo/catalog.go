package main

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"text/template"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	"github.com/yuin/goldmark"
)

type Mentions = []*ast.Mention
type ByFile[T any] map[*File]T
type ByPart[T any] map[ast.MentionablePart]T
type ByMentionable[T any] map[ast.Mentionable]T

type Catalog struct {
	MarkdownParser    goldmark.Markdown
	HttpCache         *cache.HTTPCache
	AvatarPaths       *AvatarPaths
	ByMentionable     ByMentionable[ByFile[Mentions]]
	ByMentionablePart ByPart[ByPart[ByFile[Mentions]]]
	Files             []*File
	Mutex             sync.RWMutex
}

func (c *Catalog) NewFile(filePath string) error {
	file, err := NewFile(filePath, c.MarkdownParser, c.HttpCache, c.AvatarPaths)
	if err != nil {
		return fmt.Errorf("Failed to instantiate file: %v", err)
	}

	c.Mutex.Lock()
	c.Files = append(c.Files, file)

	for mentionable, mentions := range file.Mentionables {
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
			c.ByMentionable[mentionable] = ByFile[Mentions]{}
		}
		c.ByMentionable[mentionable][file] = mentions

		if c.ByMentionablePart[mentionable.Primary] == nil {
			c.ByMentionablePart[mentionable.Primary] = ByPart[ByFile[Mentions]]{}
		}
		if c.ByMentionablePart[mentionable.Primary][mentionable.Secondary] == nil {
			c.ByMentionablePart[mentionable.Primary][mentionable.Secondary] = ByFile[Mentions]{}
		}
		c.ByMentionablePart[mentionable.Primary][mentionable.Secondary][file] = mentions
	}
	c.Mutex.Unlock()

	return nil
}

func (c *Catalog) RenderMentionPages() error {
	for primary, mentionsByFileBySecondary := range c.ByMentionablePart {
		primaries := mentionsByFileBySecondary[ast.EmptyMentionablePart]
		delete(mentionsByFileBySecondary, ast.EmptyMentionablePart)

		file, _ := makePage(unencode(primary.ID()))
		component := MentionPage(primary, primaries, mentionsByFileBySecondary, c.HttpCache)
		err := component.Render(context.Background(), file)
		if err != nil {
			return fmt.Errorf("Failed to render template: %v", err)
		}
		//log.Printf("Wrote %v", outPath)
	}
	return nil
}

func (c *Catalog) RenderPopups() error {
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

		f, _ := makePage(unencode(location))
		component := MentionablePopup(mentionable, mentionsByFile, otherMentionables)
		err := component.Render(context.Background(), f)
		if err != nil {
			return fmt.Errorf("Failed to render template: %v", err)
		}
		//log.Printf("Wrote %v", outPath)
	}

	return nil
}

func (c *Catalog) SortFilesByDate() {
	slices.SortFunc(c.Files, filesByDate)
}

func (c *Catalog) CompletedFiles() []*File {
	var files []*File
	for _, file := range c.Files {
		if file.IsComplete() {
			files = append(files, file)
		}
	}
	return files
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
		"SuspectTag":    string(suspect.Source.Segment.Value(suspect.File.GetMarkdown())),
		"Suspect":       suspect,
		"SuggestionTag": string(suggestion.Source.Segment.Value(suggestion.File.GetMarkdown())),
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
