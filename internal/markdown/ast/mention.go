package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"golang.org/x/net/html"
)

var countKey = gparser.NewContextKey()

type Mention struct {
	BaseInline
	FileNode

	Mentionable Mentionable
	Label       string
	Occurance   int
	File        File

	Source            Source
	MentionableSource Source
}

type Source struct {
	Row     int
	Col     int
	Segment text.Segment
}

func NewMention(pc gparser.Context, mentionable Mentionable, label string) *Mention {
	count := pc.ComputeIfAbsent(countKey, func() interface{} {
		return map[Mentionable]int{}
	}).(map[Mentionable]int)

	count[mentionable] += 1
	pc.Set(countKey, count)

	file := GetFile(pc)

	return &Mention{
		BaseInline:  BaseInline{},
		Mentionable: mentionable,
		Label:       label,
		File:        file,
		Occurance:   count[mentionable],
	}
}

func (m *Mention) LocalID() string {
	id := m.Mentionable.ID()
	if m.Occurance > 1 {
		id += "-" + fmt.Sprint(m.Occurance)
	}
	return id
}

func (m *Mention) ID() string {
	return m.LocalID() + "@" + m.File.GetID()
}

func (m *Mention) Permalink() string {
	return m.File.GetPermalink() + "#" + m.LocalID()
}

func (m *Mention) Title() string {
	if len(m.Label) > 0 {
		return m.Label
	} else {
		return m.Mentionable.Ultimate().PrefixFirst()
	}
}

func (m *Mention) VignetteHTML(source []byte, radius int) string {
	for p := m.Parent(); p != nil; p = p.Parent() {
		if p.Kind() == KindSpeaker {
			before, after, _ := cutText(source, p, m, false)

			rBefore := []rune(before)
			rAfter := []rune(after)

			var result []rune
			if len(rBefore) > radius {
				result = append(result, []rune("... ")...)
			}
			result = append(result, rBefore[max(0, len(rBefore)-radius):]...)

			result = append(result, []rune(fmt.Sprintf(
				`<mark><a id="%v" href="%v" class="underline">%s</a></mark>`,
				m.ID(),
				m.Permalink(),
				m.Text(source),
			))...)

			result = append(result, rAfter[:min(len(rAfter), radius)]...)

			if len(rAfter) > radius {
				result = append(result, []rune(" ...")...)
			}

			return string(result)
		}
	}
	panic("Failed to find parent speaker block for mention node")
}

func (m *Mention) Dump(source []byte, level int) {
	gast.DumpHelper(m, source, level, nil, nil)
}

func (m *Mention) Text(source []byte) []byte {
	if len(m.Label) > 0 {
		return []byte(m.Label)
	} else {
		return []byte(m.Mentionable.Ultimate().PrefixFirst())
	}
}

var KindMention = gast.NewNodeKind("Mention")

func (m *Mention) Kind() gast.NodeKind {
	return KindMention
}

// Recursive search of tree for a target Node that returns the Node text before
// and after the target.
func cutText(source []byte, root gast.Node, target gast.Node, found bool) (string, string, bool) {

	if root.HasChildren() {
		left, right := "", ""
		for child := root.FirstChild(); child != nil; child = child.NextSibling() {
			if child == target {
				found = true
			} else {
				l, r, f := cutText(source, child, target, found)
				left += l
				right += r
				found = f
			}
		}

		// Block level elements (paragraphs) don't have any spaces when inline
		isBlock := root.Type() == gast.TypeBlock
		if isBlock {
			if found {
				right += " "
			} else {
				left += " "
			}
		}

		return left, right, found
	} else {
		text := html.UnescapeString(string(root.Text(source)))
		if found {
			return "", text, found
		} else {
			return text, "", found
		}
	}
}
