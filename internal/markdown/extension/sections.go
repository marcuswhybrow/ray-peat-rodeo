package extension

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/renderer"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/transformer"
	"github.com/yuin/goldmark"
	gmParser "github.com/yuin/goldmark/parser"
	gmRenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type sections struct{}

var Sections = &sections{}

func NewSections() goldmark.Extender {
	return &sections{}
}

func (s *sections) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		gmParser.WithASTTransformers(
			util.Prioritized(transformer.SectionTransformer, 100),
		),
	)
	m.Renderer().AddOptions(
		gmRenderer.WithNodeRenderers(
			util.Prioritized(renderer.NewSectionHTMLRenderer(), 100),
		),
	)
}
