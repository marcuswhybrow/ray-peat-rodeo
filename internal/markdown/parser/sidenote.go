package parser

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var sidenoteCountKey = gparser.NewContextKey()

type sidenoteParser struct{}

func NewSidenotesParser() gparser.InlineParser {
	return &sidenoteParser{}
}

func (p *sidenoteParser) Trigger() []byte {
	return []byte{'{', '}'}
}

func (p *sidenoteParser) Parse(parent gast.Node, block text.Reader, pc gparser.Context) gast.Node {
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()

	// ScanDelimiter honors nesting, ensuring closing tab is a sibling
	node := gparser.ScanDelimiter(line, before, 1, &sidenotesDelimiterProcessor{pc})
	if node == nil {
		return nil
	}

	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)

	block.Advance((node.OriginalLength))
	pc.PushDelimiter(node)

	return node
}

type sidenotesDelimiterProcessor struct {
	context gparser.Context
}

func (p *sidenotesDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '{' || b == '}'
}

func (p *sidenotesDelimiterProcessor) CanOpenCloser(opener, closer *gparser.Delimiter) bool {
	return opener.Char == '{' && closer.Char == '}'
}

func (p *sidenotesDelimiterProcessor) OnMatch(consumes int) gast.Node {
	sidenodeCount := func() int {
		if extantCount := p.context.Get(sidenoteCountKey); extantCount != nil {
			return extantCount.(int)
		}
		return 0
	}()
	sidenodeCount += 1
	p.context.Set(sidenoteCountKey, sidenodeCount)
	return ast.NewSidenote(sidenodeCount)
}
