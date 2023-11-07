package extension

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/parser"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/renderer"
	"github.com/yuin/goldmark"

	gmParser "github.com/yuin/goldmark/parser"
	gmRenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type timecodes struct{}

var Timecodes = &timecodes{}

func New() goldmark.Extender {
	return &timecodes{}
}

func (e *timecodes) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		gmParser.WithInlineParsers(
			util.Prioritized(parser.NewTimecodeParser(), 100),
		),
	)
	m.Renderer().AddOptions(
		gmRenderer.WithNodeRenderers(
			util.Prioritized(renderer.NewTimecodeHTMLRenderer(), 100),
		),
	)
}
