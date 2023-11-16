package extension

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/parser"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/renderer"
	"github.com/yuin/goldmark"
	gmParser "github.com/yuin/goldmark/parser"
	gmRenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type speakers struct {
}

// Speakers is an extension for Goldmark that converts initialled markdown
// blocks into HTML divs
var Speakers = &speakers{}

func NewSpeakers() goldmark.Extender {
	return &speakers{}
}

func (s *speakers) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(gmParser.WithBlockParsers(
		util.Prioritized(parser.NewSpeakerParser(), 100),
	))
	m.Renderer().AddOptions(gmRenderer.WithNodeRenderers(
		util.Prioritized(renderer.NewUtteranceHTMLRenderer(), 100),
	))
}
