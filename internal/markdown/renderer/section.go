package renderer

import (
	"fmt"
	"html/template"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type SectionHTMLRenderer struct{}

func NewSectionHTMLRenderer() renderer.NodeRenderer {
	return &SectionHTMLRenderer{}
}

func (s *SectionHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindSection, s.renderSection)
}

func (s *SectionHTMLRenderer) renderSection(w util.BufWriter, source []byte, node gmAst.Node, entering bool) (gmAst.WalkStatus, error) {
	if entering {
		t, err := template.New("section").Parse(`
      <rpr-section 
        id="{{ .ID }}" 
        title="{{ .Title }}" 
        level="{{ .Level }}"
        prefix="{{ .PrefixString }}"
        {{ if .Timecode }}timecode="{{ .Timecode.Terse }}"{{ end }}
      >
    `)

		if err != nil {
			return gmAst.WalkStop, fmt.Errorf("Failed to parse section html/template: %v", err)
		}

		t.Execute(w, node.(*ast.Section))
	} else {
		w.WriteString("</rpr-section>")
	}
	return gmAst.WalkContinue, nil
}
