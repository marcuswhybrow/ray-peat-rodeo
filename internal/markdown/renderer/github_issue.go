package renderer

import (
	"fmt"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	grenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type GitHubIssueHTMLRenderer struct{}

func NewGitHubIssueHTMLRenderer() grenderer.NodeRenderer {
	return &GitHubIssueHTMLRenderer{}
}

func (r *GitHubIssueHTMLRenderer) RegisterFuncs(reg grenderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindGitHubIssue, r.renderGitHubIssue)
}

func (t *GitHubIssueHTMLRenderer) renderGitHubIssue(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		github_issue := node.(*ast.GitHubIssue)
		_, _ = w.WriteString(fmt.Sprintf(`<rpr-issue issue-id="%v" url="%v" title="%v">`, github_issue.Id, github_issue.Url(), github_issue.Title))
	} else {
		_, _ = w.WriteString("</rpr-issue>")
	}

	return gast.WalkContinue, nil
}
