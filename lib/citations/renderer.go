package citations

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type CitationRenderer struct {
}

func (r *CitationRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCitationNode, r.renderCitation)
}

func (t *CitationRenderer) renderCitation(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	citation := node.(*CitationNode)
	if entering {
		citation.renderEntering(w)
	} else {
		citation.renderExiting(w)
	}
	return ast.WalkContinue, nil
}

func withLinkRenderer(getAttrs func() map[string]string) (func(util.BufWriter), func(util.BufWriter)) {
	entering := func(b util.BufWriter) {
		attrs := getAttrs()
		text := attrs["text"]
		delete(attrs, "text")
		attrs["target"] = "_blank"
		attrsStr := func() string {
			attrsStr := ""
			for key, val := range attrs {
				attrsStr += fmt.Sprintf(` %s="%s"`, key, val)
			}
			return attrsStr
		}()
		b.WriteString(fmt.Sprintf(`<a%s>`, attrsStr))
		b.WriteString(text)
	}
	exiting := func(b util.BufWriter) {
		b.WriteString("</a>")
	}
	return entering, exiting
}
