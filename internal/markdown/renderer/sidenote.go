package renderer

import (
	"fmt"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	grenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type SidenoteHTMLRenderer struct{}

func NewSidenoteHTMLRenderer() grenderer.NodeRenderer {
	return &SidenoteHTMLRenderer{}
}

func (r *SidenoteHTMLRenderer) RegisterFuncs(reg grenderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindSidenote, r.renderSidenote)
}

func (t *SidenoteHTMLRenderer) renderSidenote(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		sidenote := node.(*ast.Sidenote)
		_, _ = w.WriteString(fmt.Sprintf(`<rpr-sidenote sidenote-id="%v">`, sidenote.Position))
	} else {
		_, _ = w.WriteString("</rpr-sidenote>")
	}

	return gast.WalkContinue, nil
}
