package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
)

// Markdown input file
type File struct {
	FrontMatter   ast.FrontMatter
	IsTodo        bool
	Path          string
	ID            string
	OutPath       string
	Date          string
	Markdown      []byte
	Html          []byte
	IssueCount    int
	EditPermalink string
	Permalink     string
	Mentions      Mentions
	Mentionables  ByMentionable[Mentions]
}

func NewFile(filePath string) *File {
	fileName := filepath.Base(filePath)
	fileStem := strings.TrimSuffix(fileName, filepath.Ext(filePath))

	markdownBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("Failed to read markdown file '%v': %v", filePath, err)
	}

	id := fileStem
	permalink := "/" + id
	outPath := path.Join(BUILD, fileStem, "index.html")
	parentPath := path.Dir(filePath)
	parentName := path.Base(parentPath)

	return &File{
		ID:           id,
		Path:         filePath,
		OutPath:      outPath,
		Date:         fileStem[:11],
		Permalink:    permalink,
		IsTodo:       parentName == "todo",
		Markdown:     markdownBytes,
		Mentions:     Mentions{},
		Mentionables: ByMentionable[Mentions]{},
	}
}

// For ast.File interface

func (f *File) GetMarkdown() []byte {
	return f.Markdown
}

func (f *File) GetPath() string {
	return f.Path
}

func (f *File) RegisterMention(mention *ast.Mention) {
	f.Mentions = append(f.Mentions, mention)
	f.Mentionables[mention.Mentionable] = append(f.Mentionables[mention.Mentionable], mention)
}

func (f *File) GetID() string {
	return f.ID
}

func (f *File) GetPermalink() string {
	return f.Permalink
}

type ByDate []*File

func (f ByDate) Len() int { return len(f) }

func (f ByDate) Less(i, j int) bool {
	return f[i].Date > f[j].Date
}

func (f ByDate) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type ByTranscriptionDate []*File

func (f ByTranscriptionDate) Len() int { return len(f) }

func (f ByTranscriptionDate) Less(i, j int) bool {
	return f[i].FrontMatter.Transcription.Date > f[j].FrontMatter.Transcription.Date
}

func (f ByTranscriptionDate) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type Mentions = []*ast.Mention
type ByFile[T any] map[*File]T
type ByPart[T any] map[ast.MentionablePart]T
type ByMentionable[T any] map[ast.Mentionable]T

type Files struct {
	Popups  ByMentionable[ByFile[Mentions]]
	Catalog ByPart[ByPart[ByFile[Mentions]]]
	Files   []*File
}

func NewFiles() *Files {
	return &Files{
		Popups:  ByMentionable[ByFile[Mentions]]{},
		Catalog: ByPart[ByPart[ByFile[Mentions]]]{},
		Files:   []*File{},
	}
}

func (f *Files) Add(file *File) {
	f.Files = append(f.Files, file)

	for mentionable, mentions := range file.Mentionables {
		for existingMentionable, existingByFile := range f.Popups {
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

		if f.Popups[mentionable] == nil {
			f.Popups[mentionable] = ByFile[Mentions]{}
		}
		f.Popups[mentionable][file] = mentions

		if f.Catalog[mentionable.Primary] == nil {
			f.Catalog[mentionable.Primary] = ByPart[ByFile[Mentions]]{}
		}
		if f.Catalog[mentionable.Primary][mentionable.Secondary] == nil {
			f.Catalog[mentionable.Primary][mentionable.Secondary] = ByFile[Mentions]{}
		}
		f.Catalog[mentionable.Primary][mentionable.Secondary][file] = mentions
	}
}

func anyValue[Key comparable, Val any](m map[Key]Val) Val {
	for _, v := range m {
		return v
	}
	panic("Failed to find any entries in map")
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
