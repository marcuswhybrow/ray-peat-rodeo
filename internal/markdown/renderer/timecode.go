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

	timecodeExtUrl, err := timecode.ExternalUrl()
	if err != nil {
		return gmAst.WalkStop, fmt.Errorf("Failed to determine timecode external URL: %v", err)
	}

	linkClass := "text-sm px-2 py-1 rounded-md "
	if ast.IsRaySpeaking(node) {
		linkClass += "is-not-ray bg-gray-300 hover:bg-gray-500 text-gray-50"
	} else {
		linkClass += "is-ray bg-sky-300 hover:bg-sky-500 text-sky-50"
	}

	if entering {
		w.WriteString(`<span class="timecode text-right">`)
		w.WriteString(fmt.Sprintf(
			`<a href="%v" class="%v">`,
			timecodeExtUrl.String(),
			linkClass,
		))
		w.WriteString(timecode.Terse())
	} else {
		w.WriteString("</a>")
		w.WriteString("</span>")
	}

	return gmAst.WalkContinue, nil
}
