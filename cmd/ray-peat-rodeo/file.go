package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gernest/front"
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark"
	gparser "github.com/yuin/goldmark/parser"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/global"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
)

// Markdown input file
type File struct {
	FrontMatter   FrontMatter
	Path          string
	ID            string
	OutPath       string
	Date          string
	Markdown      []byte
	Html          []byte
	EditPermalink string
	RawPermalink  string
	Permalink     string
	Mentions      Mentions
	Mentionables  ByMentionable[Mentions]
	Issues        []int
	Speakers      []*Speaker
}

type FrontMatter struct {
	Source struct {
		Series   string
		Title    string
		Url      string
		Kind     string
		Duration string
	}
	Speakers        map[string]string
	PrimarySpeakers []string
	Transcription   struct {
		Url    string
		Kind   string
		Date   string
		Author string
	}
	Added struct {
		Date   string
		Author string
	}
	Completion struct {
		Content            bool
		ContentVerified    bool `mapstructure:"content-verified"`
		SpeakersIdentified bool `mapstructure:"speakers-identified"`
		Mentions           bool
		Issues             bool
		Notes              bool
		Timestamps         bool
	}
}

func NewFile(filePath string, markdownParser goldmark.Markdown, httpCache *cache.HTTPCache, avatarPaths *AvatarPaths) (*File, error) {
	fileName := filepath.Base(filePath)
	fileStem := strings.TrimSuffix(fileName, filepath.Ext(filePath))

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}

	// 🔗 Details

	id := fileStem[11:]
	permalink := "/" + id
	editPermalink := global.GITHUB_LINK + path.Join("/edit/main", filePath)
	rawPermalink := global.GITHUB_LINK + path.Join("/raw/main", filePath)
	outPath := path.Join(id, "index.html")

	// 📄 FrontMatter

	matter := front.NewMatter()
	matter.Handle("---", front.YAMLHandler)
	rawFMatter, _, err := matter.Parse(strings.NewReader(string(fileBytes)))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse frontmatter: %v", err)
	}

	frontMatter := FrontMatter{}
	err = mapstructure.Decode(rawFMatter, &frontMatter)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode YAML frontmatter: %v", err)
	}

	// 👨👱 Speakers & Avatars

	speakers := []*Speaker{}
	for id, name := range frontMatter.Speakers {
		var isPrimarySpeaker bool
		if len(frontMatter.PrimarySpeakers) == 0 {
			isPrimarySpeaker = id == "RP"
		} else {
			isPrimarySpeaker = slices.Contains(frontMatter.PrimarySpeakers, id)
		}

		speakers = append(speakers, &Speaker{
			ID:               id,
			Name:             name,
			AvatarPath:       avatarPaths.Get(name),
			IsPrimarySpeaker: isPrimarySpeaker,
		})
	}

	file := &File{
		ID:            id,
		Path:          filePath,
		OutPath:       outPath,
		Date:          fileStem[:10],
		Permalink:     permalink,
		EditPermalink: editPermalink,
		RawPermalink:  rawPermalink,
		FrontMatter:   frontMatter,
		Markdown:      fileBytes,
		Mentions:      Mentions{},
		Mentionables:  ByMentionable[Mentions]{},
		Speakers:      speakers,
	}

	// 🖥 HTML

	parserContext := gparser.NewContext()
	parserContext.Set(ast.FileKey, file)
	parserContext.Set(ast.HTTPCacheKey, httpCache)

	var html bytes.Buffer
	err = markdownParser.Convert(file.Markdown, &html, gparser.WithContext(parserContext))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse markdown: %v", err)
	}
	file.Html = html.Bytes()
	return file, nil
}

// True if all known frontmatter `completion` fields are true
func (f *File) IsComplete() bool {
	c := f.FrontMatter.Completion
	return c.Content && c.ContentVerified && c.SpeakersIdentified &&
		c.Mentions && c.Issues && c.Notes && c.Timestamps
}

// Writes file to f.outPath
func (f *File) Render() error {
	buildFile, _ := MakeFile(f.OutPath)

	err := RenderChat(f).Render(context.Background(), buildFile)
	if err != nil {
		return fmt.Errorf("Failed to render template: %v", err)
	}

	return nil
}

// Implement ast.File interface

// Returns the raw source markdown (without any file frontmatter)
func (f *File) GetMarkdown() []byte {
	return f.Markdown
}

// Returns the source file path
func (f *File) GetPath() string {
	return f.Path
}

func (f *File) RegisterMention(mention *ast.Mention) {
	f.Mentions = append(f.Mentions, mention)
	mention.Position = len(f.Mentions)
	f.Mentionables[mention.Mentionable] = append(f.Mentionables[mention.Mentionable], mention)
}

func (f *File) GetSpeakers() []ast.Speaker {
	speakers := make([]ast.Speaker, len(f.Speakers))
	for i, s := range f.Speakers {
		speakers[i] = s
	}
	return speakers
}

func (f *File) GetSourceURL() string {
	return f.FrontMatter.Source.Url
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

// Other

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

	slices.SortFunc(results, mostMentioned)
	return results
}

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

	slices.SortFunc(results, mostMentionedPrimary)
	return results
}

func (f *File) IssueCount() int {
	return len(f.Issues)
}

func (f *File) HasIssues() bool {
	return f.IssueCount() > 0
}

func (f *File) TopSpeakers() []*Speaker {
	speakers := make([]*Speaker, len(f.Speakers))
	copy(speakers, f.Speakers)

	slices.SortFunc(speakers, func(a, b *Speaker) int {
		// Prefer speakers with avatars
		aScore := 0
		bScore := 0

		if len(a.GetAvatarPath()) > 0 {
			aScore += 1
		}
		if len(b.GetAvatarPath()) > 0 {
			bScore += 1
		}

		if !strings.Contains(a.Name, "(") {
			aScore += 1
		}
		if !strings.Contains(b.Name, "(") {
			bScore += 1
		}

		if a.GetIsPrimarySpeaker() {
			aScore += 1
		}
		if b.GetIsPrimarySpeaker() {
			bScore += 1
		}

		if aScore > bScore {
			return -1
		} else if aScore == bScore {
			return 0
		} else {
			return 1
		}
	})
	return speakers
}

func (f *File) IsAutoGenerated() bool {
	fm := f.FrontMatter
	return fm.Completion.Content && fm.Transcription.Kind == "auto-generated"
}

// type Speaker struct {
// 	Key    string
// 	Name   string
// 	Avatar string
// }

type MentionCount struct {
	Mention *ast.Mention
	Count   int
}

func mostMentioned(a, b MentionCount) int {
	if a.Count > b.Count {
		return -1
	} else if a.Count == b.Count {
		aCardinal := a.Mention.Mentionable.Ultimate().Cardinal
		bCardinal := b.Mention.Mentionable.Ultimate().Cardinal
		if aCardinal == bCardinal {
			return 0
		} else if len(aCardinal) > len(bCardinal) {
			return -1
		} else {
			return 1
		}
	}
	return 1
}

// Struct for counting mentionable parts
type MentionablePartCount struct {
	MentionablePart ast.MentionablePart
	Count           int
}

// Comparison function used to order MentionablePartCount's
// Orders first the MPart with the highest count.
// In a tie, the one with a shorter cardinal is prefered.
func mostMentionedPrimary(a, b MentionablePartCount) int {
	if a.Count > b.Count {
		return -1
	} else if a.Count == b.Count {
		aCardinal := a.MentionablePart.Cardinal
		bCardinal := b.MentionablePart.Cardinal
		if len(aCardinal) < len(bCardinal) {
			return -1
		} else if aCardinal == bCardinal {
			return 0
		} else {
			return 1
		}
	}
	return 1
}

func filesByDate(a *File, b *File) int {
	if a.Date > b.Date {
		return -1
	} else if a.Date < b.Date {
		return 1
	} else {
		return 0
	}
}

func filesByDateAdded(a *File, b *File) int {
	if a.FrontMatter.Added.Date > b.FrontMatter.Added.Date {
		return -1
	} else if a.FrontMatter.Added.Date < b.FrontMatter.Added.Date {
		return 1
	} else {
		return 0
	}
}

func unencode(filePath string) string {
	str, err := url.QueryUnescape(filePath)
	if err != nil {
		log.Panicf("Failed to unescape path '%v': %v", filePath, err)
	}
	return str
}

// For a given speaker's full name, returns path to avatar image.
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
