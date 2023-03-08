package citations

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

var KindCitationNode = ast.NewNodeKind("Citation")

type CitationNode struct {
	ast.BaseInline
	renderEntering func(b util.BufWriter)
	renderExiting  func(b util.BufWriter)
}

func (t *CitationNode) Dump(source []byte, level int) {
	ast.DumpHelper(t, source, level, nil, nil)
}

func (n *CitationNode) Kind() ast.NodeKind {
	return KindCitationNode
}

func NewCitationNode(renderEntering, renderExiting func(b util.BufWriter)) *CitationNode {
	return &CitationNode{
		BaseInline:     ast.BaseInline{},
		renderEntering: renderEntering,
		renderExiting:  renderExiting,
	}
}
