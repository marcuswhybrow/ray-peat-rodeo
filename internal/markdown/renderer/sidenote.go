package renderer

import (
	"fmt"
	"html/template"

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

		// Uses CSS counters. Requires "counter-reset:sidenote" on parent
		t, err := template.New("openSidenote").Parse(`
      <label 
        for="sidenote-{{.SidenoteId}}" 
        class="[counter-increment:sidenote] after:content-[counter(sidenote)] after:-top-1 after:left-0 after:align-baseline after:text-sm after:relative after:-top-1 font-serif after:bg-white after:rounded-md after:shadow after:text-gray-600 after:py-1 after:px-2"
      ></label><span
        id="sidenote-{{.SidenoteId}}"
        class="z-0 block bg-white rounded-md shadow w-1/2 mr-[-5%] sm:mr-[-10%] md:mr-[-15%] lg:mr-[-25%] float-right clear-right text-sm relative p-4 before:content-[counter(sidenote)_'.'] before:float-left m-2 before:mr-1 before:text-gray-500 leading-5 align-middle transition-all"
      >
    `)
		if err != nil {
			return gast.WalkStop, fmt.Errorf("Failed to parse sidenote html/template: %v", err)
		}

		t.Execute(w, map[string]string{
			"SidenoteId": fmt.Sprint(sidenote.Position),
		})

	} else {
		_, _ = w.WriteString("</span>")
	}

	return gast.WalkContinue, nil
}
