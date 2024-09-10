package renderer

import (
	"fmt"
	"html/template"

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

		t, err := template.New("openGitHubIssue").Parse(`
      <rpr-issue issue-id="{{ .id }}" url="{{ .url }}" title="{{ .title }}">
    `)
		// t, err := template.New("openGitHubIssue").Parse(`
		// <a
		// id="issue-{{ .id }}"
		// href="{{ .url }}"
		// class="z-10 block transition-all m-2 p-4 hover:translate-y-1 shadow-xl hover:shadow-2xl shadow-yellow-800/20 hover:shadow-yellow-600/40 rounded-md bg-gradient-to-br from-yellow-200 from-10% to-amber-200 hover:from-yellow-100 hover:from-70% hover:to-amber-200 xl:block w-2/5 mr-[-5%] md:mr-[-10%] lg:mr-[-20%] float-right clear-right text-sm relative leading-5 tracking-tight">

		// <span class="text-yellow-900 font-bold mr-0.5">
		// <img src="/assets/images/github-mark.svg" class="h-4 w-4 inline-block relative top-[-1px] mr-0.5" /> #{{ .id }}
		// </span>

		// <span class="text-yellow-800">{{ .title }}</span>
		// `)
		if err != nil {
			return gast.WalkStop, fmt.Errorf("Failed to parse GitHub issue html/template: %v", err)
		}

		t.Execute(w, map[string]string{
			"url":   github_issue.Url(),
			"id":    fmt.Sprint(github_issue.Id),
			"title": github_issue.Title,
		})

	} else {
		_, _ = w.WriteString("</rpr-issue>")
	}

	return gast.WalkContinue, nil
}
