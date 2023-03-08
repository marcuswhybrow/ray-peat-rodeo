package citations

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type AggregatedCitations struct {
	Count         int
	People        map[string]Person
	Books         map[string]Book
	SciencePapers map[string]SciencePaper
	ExternalLinks map[string]ExternalLink
}

func Get(pc parser.Context) AggregatedCitations {
	return AggregatedCitations{
		People:        pc.Get(peopleContextKey).(map[string]Person),
		Books:         pc.Get(booksContextKey).(map[string]Book),
		SciencePapers: pc.Get(sciencePapersContextKey).(map[string]SciencePaper),
		ExternalLinks: pc.Get(externalLinksContextKey).(map[string]ExternalLink),
	}
}

type citations struct {
}

// Citations is an extension for Goldmark that replaces [[tags]] with links
var Citations = &citations{}

func New() goldmark.Extender {
	return &citations{}
}

func (e *citations) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(NewParser(), 100)),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(util.Prioritized(&CitationRenderer{}, 100)),
	)
}
