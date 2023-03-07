package sidenotes

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var sidenoteCountKey = parser.NewContextKey()

type sidenotesDelimiterProcessor struct {
	context parser.Context
}

func (p *sidenotesDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '{' || b == '}'
}

func (p *sidenotesDelimiterProcessor) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == '{' && closer.Char == '}'
}

func (p *sidenotesDelimiterProcessor) OnMatch(consumes int) ast.Node {
	sidenodeCount := func() int {
		if extantCount := p.context.Get(sidenoteCountKey); extantCount != nil {
			return extantCount.(int)
		}
		return 0
	}()
	sidenodeCount += 1
	p.context.Set(sidenoteCountKey, sidenodeCount)
	return NewSidenote(sidenodeCount)
}

type sidenotesParser struct {
}

func NewSidenotesParser() parser.InlineParser {
	return &sidenotesParser{}
}

func (w *sidenotesParser) Trigger() []byte {
	return []byte{'{', '}'}
}

func (w *sidenotesParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()
	node := parser.ScanDelimiter(line, before, 1, &sidenotesDelimiterProcessor{pc})
	if node == nil {
		return nil
	}
	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)
	block.Advance((node.OriginalLength))
	pc.PushDelimiter(node)
	return node
}

type Sidenote struct {
	ast.BaseInline
	position int
}

func (t *Sidenote) Dump(source []byte, level int) {
	ast.DumpHelper(t, source, level, nil, nil)
}

var KindSidenote = ast.NewNodeKind("Sidenote")

func (n *Sidenote) Kind() ast.NodeKind {
	return KindSidenote
}

func NewSidenote(position int) *Sidenote {
	return &Sidenote{
		BaseInline: ast.BaseInline{},
		position:   position,
	}
}

type SidenoteHTMLRendereer struct {
}

func NewSidenoteHTMLRenderer() renderer.NodeRenderer {
	return &SidenoteHTMLRendereer{}
}

func (r *SidenoteHTMLRendereer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindSidenote, r.renderSidenote)
}

func (t *SidenoteHTMLRendereer) renderSidenote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		sidenote := node.(*Sidenote)
		position := fmt.Sprint(sidenote.position)
		w.WriteString(`
			<label for="sidenote-` + position + `" class="sidenote-toggle sidenote-number"></label>
			<input type="checkbox" id="sidenote-` + position + `" class="sidenote-toggle" />
			<span class="sidenote">
		`)
	} else {
		_, _ = w.WriteString("</span>")
	}
	return ast.WalkContinue, nil
}

type sidenotes struct {
}

// Sidenotes is an extension for Goldmark that replaces [00:00:00] with links
var Sidenotes = &sidenotes{}

func New() goldmark.Extender {
	return &sidenotes{}
}

func (e *sidenotes) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewSidenotesParser(), 1),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewSidenoteHTMLRenderer(), 1),
	))
}
