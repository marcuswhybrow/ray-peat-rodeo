package parser

import (
	"bytes"
	"strings"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
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

	httpCache := pc.Get(ast.HTTPCache).(*cache.HTTPCache)
	primaryPart := ast.NewMentionPart(pCardinal, pPrefix, httpCache)
	secondaryPart := ast.NewMentionPart(sCardinal, sPrefix, httpCache)

	mention := ast.NewMention(pc, *primaryPart, *secondaryPart, label)

	mentions := GetMentions(pc)
	mentions = append(mentions, mention)
	pc.Set(mentionsKey, mentions)

	block.Advance(4 + len(inside))

	return mention
}

func (p *mentionParser) CloseBlock(parent gast.Node, pc gparser.Context) {
	// do nothing
}
