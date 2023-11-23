package main

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/global"
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
	EditPermalink string
	Permalink     string
	Mentions      Mentions
	Mentionables  ByMentionable[Mentions]
	Issues        []int
}

func NewFile(filePath string) *File {
	fileName := filepath.Base(filePath)
	fileStem := strings.TrimSuffix(fileName, filepath.Ext(filePath))

	markdownBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("Failed to read markdown file '%v': %v", filePath, err)
	}

	id := fileStem[11:]
	permalink := "/" + id
	editPermalink := global.GITHUB_LINK + path.Join("/edit/main", filePath)
	outPath := path.Join(BUILD, id, "index.html")
	parentPath := path.Dir(filePath)
	parentName := path.Base(parentPath)

	return &File{
		ID:            id,
		Path:          filePath,
		OutPath:       outPath,
		Date:          fileStem[:10],
		Permalink:     permalink,
		EditPermalink: editPermalink,
		IsTodo:        parentName == "todo",
		Markdown:      markdownBytes,
		Mentions:      Mentions{},
		Mentionables:  ByMentionable[Mentions]{},
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
	mention.Position = len(f.Mentions)
	f.Mentionables[mention.Mentionable] = append(f.Mentionables[mention.Mentionable], mention)
}

type MentionCount struct {
	Mention *ast.Mention
	Count   int
}

type ByMostMentioned []MentionCount

func (m ByMostMentioned) Len() int { return len(m) }

func (m ByMostMentioned) Less(i, j int) bool {
	if m[i].Count > m[j].Count {
		return true
	} else if m[i].Count == m[j].Count {
		iCardinal := m[i].Mention.Mentionable.Ultimate().Cardinal
		jCardinal := m[j].Mention.Mentionable.Ultimate().Cardinal
		return len(iCardinal) < len(jCardinal)
	}
	return false
}

func (m ByMostMentioned) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

func (f *File) TopMentions() []MentionCount {
	results := []MentionCount{}

	for _, m := range f.Mentions {
		i := slices.IndexFunc(results, func(ms MentionCount) bool {
			return ms.Mention.Mentionable == m.Mentionable
		})
		if i >= 0 {
			results[i].Count += 1
			if m.Occurance == 1 {
				results[i].Mention = m
			}
		} else {
			results = append(results, MentionCount{
				Mention: m,
				Count:   1,
			})
		}
	}

	sort.Sort(ByMostMentioned(results))
	return results
}

type MentionablePartCount struct {
	MentionablePart ast.MentionablePart
	Count           int
}

type ByMostMentionedPrimary []MentionablePartCount

func (m ByMostMentionedPrimary) Len() int { return len(m) }

func (m ByMostMentionedPrimary) Less(i, j int) bool {
	if m[i].Count > m[j].Count {
		return true
	} else if m[i].Count == m[j].Count {
		iCardinal := m[i].MentionablePart.Cardinal
		jCardinal := m[j].MentionablePart.Cardinal
		return len(iCardinal) < len(jCardinal)
	}
	return false
}

func (m ByMostMentionedPrimary) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

func (f *File) TopPrimaryMentionables() []MentionablePartCount {
	results := []MentionablePartCount{}

	for _, m := range f.Mentions {
		i := slices.IndexFunc(results, func(ms MentionablePartCount) bool {
			return ms.MentionablePart == m.Mentionable.Primary
		})
		if i >= 0 {
			results[i].Count += 1
		} else {
			results = append(results, MentionablePartCount{
				MentionablePart: m.Mentionable.Primary,
				Count:           1,
			})
		}
	}

	sort.Sort(ByMostMentionedPrimary(results))
	return results
}

func (f *File) IssueCount() int {
	return len(f.Issues)
}

func (f *File) HasIssues() bool {
	return f.IssueCount() > 0
}

func (f *File) RegisterIssue(id int) {
	f.Issues = append(f.Issues, id)
}

func (f *File) GetID() string {
	return f.ID
}

func (f *File) GetPermalink() string {
	return f.Permalink
}

func (f *File) TopSpeakers() []Speaker {
	speakers := []Speaker{}
	for key, name := range f.FrontMatter.Speakers {
		avatar, _ := SpeakerAvatar(name)
		speakers = append(speakers, Speaker{
			Key:    key,
			Name:   name,
			Avatar: avatar,
		})
	}
	slices.SortFunc(speakers, func(a, b Speaker) int {
		// Prefer speakers with avatars
		if len(a.Avatar) > 0 && len(b.Avatar) == 0 {
			return -1
		}
		if len(b.Avatar) > 0 && len(a.Avatar) == 0 {
			return 1
		}

		// Prefer speakers without parenthesis: "Audience Member (Male)"
		if strings.Contains(a.Name, "(") && !strings.Contains(b.Name, "(") {
			return -1
		}
		if strings.Contains(b.Name, "(") && !strings.Contains(a.Name, "(") {
			return 1
		}

		return 0
	})
	return speakers
}

type Speaker struct {
	Key    string
	Name   string
	Avatar string
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

type ByAddedDate []*File

func (f ByAddedDate) Len() int { return len(f) }

func (f ByAddedDate) Less(i, j int) bool {
	return f[i].FrontMatter.Added.Date > f[j].FrontMatter.Added.Date
}

func (f ByAddedDate) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

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

func AtMost[T any](ts []T, i int) []T {
	if len(ts) > i {
		return ts[:i]
	}
	return ts
}

func SpeakerAvatar(speakerName string) (string, bool) {
	speakerName = strings.ToLower(speakerName)
	speakerName = strings.ReplaceAll(speakerName, " ", "-")
	found := ""

	fs.WalkDir(os.DirFS("./internal"), "assets/images/avatars", func(filePath string, entry fs.DirEntry, err error) error {
		fileStem := path.Base(filePath)
		ext := path.Ext(fileStem)
		fileName, _ := strings.CutSuffix(fileStem, ext)

		if speakerName == fileName {
			found = "/" + filePath
		}
		return nil
	})
	return found, len(found) > 0
}
