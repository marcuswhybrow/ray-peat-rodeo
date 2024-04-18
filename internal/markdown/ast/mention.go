package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"golang.org/x/net/html"
)

var countKey = gparser.NewContextKey()

// Represents a mention in the markdown parser's Abstract Syntax Tree.
type Mention struct {
	BaseInline
	FileNode

	// The thing being mentioned
	Mentionable Mentionable

	// An option label to override the full name of this mention
	Label string

	// The number of mentions referring to the this mentionable before this, + 1
	Occurance int

	// The number of mentions before this, + 1
	Position int

	// The file this mention is in
	File File

	// The markdown source which created this mention
	Source Source

	// The markdown source which defines the mentionable
	MentionableSource Source
}

// A location in a source markdown file
type Source struct {
	Row     int
	Col     int
	Segment text.Segment
}

// Creates a new Mention with the correct context
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

// Computes a ID for this mention that's unique to this document.
// Because a single mentionable maybe mentioned multiple times,
// the local ID takes the number of occurances of a particular
// mentionable into account, in order to compute a unique local ID.
func (m *Mention) LocalID() string {
	id := m.Mentionable.ID()
	if m.Occurance > 1 {
		id += "-" + fmt.Sprint(m.Occurance)
	}
	return id
}

// A globally unique ID for this mention in this file.
func (m *Mention) ID() string {
	return m.LocalID() + "@" + m.File.GetID()
}

// A URL that links directly to this mention using an HTTP fragment.
func (m *Mention) Permalink() string {
	return m.File.GetPermalink() + "#" + m.LocalID()
}

// Returns a mention in plain text (formatting removed), surrounded by a
// specified amount of adjacent, contextual text.
// The mentioned text is wraped in a <mark> HTML tag with baked in styling.
func (m *Mention) VignetteHTML(source []byte, radius int) string {
	for p := m.Parent(); p != nil; p = p.Parent() {
		if p.Kind() == gast.KindParagraph {
			before, after, _ := cutText(source, p, m, false)

			rBefore := []rune(before)
			rAfter := []rune(after)

			var result []rune
			if len(rBefore) > radius {
				result = append(result, []rune("... ")...)
			}
			result = append(result, rBefore[max(0, len(rBefore)-radius):]...)

			result = append(result, []rune(fmt.Sprintf(
				`<mark class="p-px text-yellow-900 hover:text-yellow-950 bg-amber-100 hover:bg-yellow-300 rounded"><a id="%v" href="%v" class="">%s</a></mark>`,
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
	panic("Failed to find parent utterance block for mention node")
}

// Required to be used as a node in the Markdown AST
func (m *Mention) Dump(source []byte, level int) {
	gast.DumpHelper(m, source, level, nil, nil)
}

// Returns the text the human reader will see.
func (m *Mention) Text(source []byte) []byte {
	if len(m.Label) > 0 {
		return []byte(m.Label)
	} else {
		return []byte(m.Mentionable.Ultimate().PrefixFirst())
	}
}

// The Markdown Abstract Syntax Tree Node Kind is a way of
// hooking our Mention class into the Markdown parser
var KindMention = gast.NewNodeKind("Mention")

// Tells the Markdown parser which Node Kind a Mention is
func (m *Mention) Kind() gast.NodeKind {
	return KindMention
}

// Recursive search of Abstracy Syntax Tree for a target Node.
// returns the Node text before and after the target.
//
// The source argument is the original markdown input string. Each AST Node
// stores a reference to a segment from this source, i.e. it doesn't have it's
// own copy. Hence we need the original source as an argument to determine an
// arbitary node's final text: the bit the a human reads.
//
// The root argument is any AST node that has children to search.
//
// The target argument will be equated against each child to find a match.
//
// The found argument should be passed as `false` in the initial call.
// This is a recursive function. Once the target is located, found is
// set to true so later iterations can act appropriately.
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
