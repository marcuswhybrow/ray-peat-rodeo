package extension

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/parser"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/renderer"
	"github.com/yuin/goldmark"
	gparser "github.com/yuin/goldmark/parser"
	grenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type sidenotes struct{}

// Sidenotes is an extension for Goldmark that replaces [00:00:00] with a link
var Sidenotes = &sidenotes{}

func NewSidenotes() goldmark.Extender {
	return &sidenotes{}
}

func (e *sidenotes) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(gparser.WithInlineParsers(
		util.Prioritized(parser.NewSidenotesParser(), 1),
	))
	m.Renderer().AddOptions(grenderer.WithNodeRenderers(
		util.Prioritized(renderer.NewSidenoteHTMLRenderer(), 1),
	))
}
