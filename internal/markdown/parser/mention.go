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

func (p *mentionParser) Parse(parent gast.Node, block text.Reader, pc gparser.Context) gast.Node {
	line, _ := block.PeekLine()
	if line[1] != '[' {
		return nil
	}

	inside, _, foundEnd := strings.Cut(string(line[2:]), "]]")
	if !foundEnd || len(inside) == 0 {
		return nil
	}

	signature, displayText, foundEnd := strings.Cut(inside, "|")

	title := signature
	if displayText != "" {
		title = displayText
	}

	mention := ast.NewMention(title)

	block.Advance(4 + len(inside))
	return mention

}

func (p *mentionParser) CloseBlock(parent gast.Node, pc gparser.Context) {
	// do nothing
}
