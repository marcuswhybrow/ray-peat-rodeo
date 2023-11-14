package parser

import (
	"bytes"
	"strings"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var mentionCountKey = gparser.NewContextKey()
var mentionsKey = gparser.NewContextKey()

func GetMentions(pc gparser.Context) []*ast.Mention {
	return pc.ComputeIfAbsent(mentionsKey, func() interface{} {
		var mentions []*ast.Mention
		return mentions
	}).([]*ast.Mention)
}

type mentionParser struct{}

func NewMentionParser() gparser.InlineParser {
	return &mentionParser{}
}

func (p *mentionParser) Trigger() []byte {
	return []byte{'['}
}

// Parses the mention tag which has several parts, some of which are optional.
// [[Primary Mention, Prefix > Secondary Mention, Prefix | Display Text]]
func (p *mentionParser) Parse(parent gast.Node, block text.Reader, pc gparser.Context) gast.Node {
	line, _ := block.PeekLine()
	if !bytes.HasPrefix(line, []byte{'[', '['}) {
		return nil
	}

	inside, _, foundCloser := strings.Cut(string(line[2:]), "]]")
	if !foundCloser || len(inside) == 0 {
		return nil
	}

	signature, label, _ := strings.Cut(inside, "|")

	primary, secondary, _ := strings.Cut(signature, ">")

	pCardinal, pPrefix, _ := strings.Cut(primary, ",")
	sCardinal, sPrefix, _ := strings.Cut(secondary, ",")

	primaryPart := ast.MentionPart{
		Cardinal: strings.Trim(pCardinal, " "),
		Prefix:   strings.Trim(pPrefix, " "),
	}
	secondaryPart := ast.MentionPart{
		Cardinal: strings.Trim(sCardinal, " "),
		Prefix:   strings.Trim(sPrefix, " "),
	}

	mention := ast.NewMention(pc, primaryPart, secondaryPart, label)

	mentions := GetMentions(pc)
	mentions = append(mentions, mention)
	pc.Set(mentionsKey, mentions)

	block.Advance(4 + len(inside))

	return mention
}

func (p *mentionParser) CloseBlock(parent gast.Node, pc gparser.Context) {
	// do nothing
}

// Part of a string that may or may not have been surrounded by double quotes (")
type segment struct {

	// True if this string was surrounded by quotations
	Quoted bool

	// The string, without any (optional) double quotation marks
	String string
}

// Parses a string that contains quoted segments into a slice of structs
// representing that state. This is usefull for ignoring control symbols within
// quoted segments of mention signatures.
//
//	assertEq(quotedSegments(`my "quoted" string`), []segment{
//	  segment{ Quoted: false, String: "my " },
//	  segment{ Quoted: true, String: "quoted" },
//	  segment{ Quoted: false, String: " string" },
//	})
func quotedSegments(s string) []segment {
	var results []segment

	for pos := 0; pos < len(s); {
		if opener := strings.IndexByte(s[pos:], '"'); opener < 0 {
			results = append(results, segment{
				Quoted: false,
				String: s[pos:],
			})
			pos = len(s)
		} else {
			if closer := strings.IndexByte(s[pos+opener+1:], '"'); closer > 0 {
				if opener > pos {
					results = append(results, segment{
						Quoted: false,
						String: s[pos:opener],
					})
				}
				results = append(results, segment{
					Quoted: true,
					String: s[pos+opener+1 : pos+opener+1+closer],
				})
				pos += opener + closer + 2
			} else {
				results = append(results, segment{
					Quoted: false,
					String: s[pos:],
				})
				pos = len(s)
			}
		}
	}

	return results
}

type parserPhase int

const (
	mentionCardinalPhase parserPhase = iota
	mentionPrefixPhase
	subMentionCardinalPhase
	subMentionPrefixPhase
	displayTextPhase
)
