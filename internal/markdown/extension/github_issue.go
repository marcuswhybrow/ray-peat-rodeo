package extension

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/parser"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/renderer"
	"github.com/yuin/goldmark"
	gparser "github.com/yuin/goldmark/parser"
	grenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type githubIssues struct{}

// GitHubIssues is an extension for Goldmark that replaces {#7} with an aside
var GitHubIssues = &githubIssues{}

func NewGitHubIssues() goldmark.Extender {
	return &githubIssues{}
}

func (e *githubIssues) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(gparser.WithInlineParsers(
		util.Prioritized(parser.NewGitHubIssueParser(), 1),
	))
	m.Renderer().AddOptions(grenderer.WithNodeRenderers(
		util.Prioritized(renderer.NewGitHubIssueHTMLRenderer(), 1),
	))
}
