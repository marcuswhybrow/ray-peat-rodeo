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
var mentionablesKey = gparser.NewContextKey()

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
	line, segment := block.PeekLine()
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

	httpCache := pc.Get(ast.HTTPCacheKey).(*cache.HTTPCache)
	primaryPart := ast.NewMentionablePart(pCardinal, pPrefix, httpCache)
	secondaryPart := ast.NewMentionablePart(sCardinal, sPrefix, httpCache)

	mentionable := ast.Mentionable{
		Primary:   primaryPart,
		Secondary: secondaryPart,
	}

	mention := ast.NewMention(pc, mentionable, label)

	row, _ := block.Position()
	mentionSegment := segment.WithStop(segment.Start + len(inside) + 4)
	mention.Source = ast.Source{
		Row:     row,
		Col:     block.LineOffset(),
		Segment: mentionSegment,
	}

	block.Advance(4 + len(inside))

	mention.Asset = ast.GetAsset(pc)
	mention.Asset.RegisterMention(mention)

	return mention
}

func (p *mentionParser) CloseBlock(parent gast.Node, pc gparser.Context) {
	// do nothing
}
