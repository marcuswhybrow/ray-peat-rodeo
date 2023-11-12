package parser

import (
	"strings"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var mentionCountKey = gparser.NewContextKey()

type mentionParser struct{}

func NewMentionParser() gparser.InlineParser {
	return &mentionParser{}
}

func (p *mentionParser) Trigger() []byte {
	return []byte{'['}
}

// Parses the mention tag which has several parts, some of which are optional.
// [[Primary Cardinal, Prefix > Secondary Cardinal, Prefix | Display Text]]
func (p *mentionParser) Parse(parent gast.Node, block text.Reader, pc gparser.Context) gast.Node {
	line, _ := block.PeekLine()
	if line[1] != '[' {
		return nil
	}

	inside, _, foundEnd := strings.Cut(string(line[2:]), "]]")
	if !foundEnd || len(inside) == 0 {
		return nil
	}

	segments := quotedSegments(inside)

	var primary ast.MentionPart
	var secondary ast.MentionPart
	var displayText string

	phase := mentionCardinalPhase
	target := func() *string {
		switch phase {
		case mentionCardinalPhase:
			return &primary.Cardinal
		case mentionPrefixPhase:
			return &primary.Prefix
		case subMentionCardinalPhase:
			return &secondary.Cardinal
		case subMentionPrefixPhase:
			return &secondary.Prefix
		case displayTextPhase:
			return &displayText
		default:
			return &primary.Cardinal
		}
	}

	for _, segment := range segments {
		if segment.Quoted {
			*target() += segment.String
		} else {
			for _, r := range segment.String {
				switch r {
				case ',':
					switch phase {
					case mentionCardinalPhase:
						phase = mentionPrefixPhase
					case subMentionCardinalPhase:
						phase = subMentionPrefixPhase
					}
				case '|':
					phase = displayTextPhase
				case '>':
					switch phase {
					case mentionCardinalPhase:
						phase = subMentionCardinalPhase
					case mentionPrefixPhase:
						phase = subMentionCardinalPhase
					}
				default:
					*target() += string(r)
				}
			}
		}
	}

	block.Advance(4 + len(inside))
	return ast.NewMention(primary, secondary, displayText)

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

// Parses a string that contains quoted segments into a list of structs
// representing that state. This is usefull for ignoring control symbols within
// quoted segments of mention signatures.
//
//	s := `my "quoted" string`
//	assertEq(quotedSegments(s), []segment{
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
