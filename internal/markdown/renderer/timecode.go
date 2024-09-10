package renderer

import (
	"fmt"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type TimecodeHTMLRenderer struct{}

func NewTimecodeHTMLRenderer() renderer.NodeRenderer {
	return &TimecodeHTMLRenderer{}
}

func (r *TimecodeHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindTimecode, r.renderTimecode)
}

func (t *TimecodeHTMLRenderer) renderTimecode(w util.BufWriter, source []byte, node gmAst.Node, entering bool) (gmAst.WalkStatus, error) {
	timecode := node.(*ast.Timecode)

	if entering {
		w.WriteString(fmt.Sprintf(`<rpr-timecode external-url="%v" time="%v" primary="%v">`, timecode.ExternalURL, timecode.Terse(), ast.IsPrimarySpeaker(timecode)))
	} else {
		w.WriteString("</rpr-timecode>")
	}

	return gmAst.WalkContinue, nil
}
