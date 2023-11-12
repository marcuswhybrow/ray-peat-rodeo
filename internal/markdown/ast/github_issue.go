package ast

import (
	"fmt"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/global"
	gast "github.com/yuin/goldmark/ast"
)

type GitHubIssue struct {
	gast.BaseInline
	Id    int
	Title string
}

func (g *GitHubIssue) Dump(source []byte, level int) {
	gast.DumpHelper(g, source, level, nil, nil)
}

func (g *GitHubIssue) Url() string {
	return fmt.Sprintf("%v/issues/%v", global.GITHUB_LINK, g.Id)
}

var KindGitHubIssue = gast.NewNodeKind("GitHubIssue")

func (g *GitHubIssue) Kind() gast.NodeKind {
	return KindGitHubIssue
}

func NewGitHubIssue(id int) *GitHubIssue {
	return &GitHubIssue{
		BaseInline: gast.BaseInline{},
		Id:         id,
		Title:      fmt.Sprintf("Issue #%v", id),
	}
}
