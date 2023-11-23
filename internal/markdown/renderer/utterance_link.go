package renderer

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	grenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type UtteranceLinkRenderer struct{}

func NewUtteranceLinkRenderer() grenderer.NodeRenderer {
	return &UtteranceLinkRenderer{}
}

func (r *UtteranceLinkRenderer) RegisterFuncs(reg grenderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindUtteranceLink, r.renderUtteranceLink)
}

func (t *UtteranceLinkRenderer) renderUtteranceLink(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(`<span class="border-b hover:border-b-2 border-gray-400">`)
	} else {
		_, _ = w.WriteString("</span>")
	}

	return gast.WalkContinue, nil
}
