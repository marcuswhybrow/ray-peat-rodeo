package renderer

import (
	"fmt"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	grenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type MentionHTMLRenderer struct{}

func NewMentionHTMLRenderer() grenderer.NodeRenderer {
	return &MentionHTMLRenderer{}
}

func (r *MentionHTMLRenderer) RegisterFuncs(reg grenderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindMention, r.renderCitation)
}

func (t *MentionHTMLRenderer) renderCitation(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		mention := node.(*ast.Mention)

		var anchorClass string
		if IsRaySpeaking(node) {
			anchorClass = "text-gray-700 hover:text-gray-900 shadow-orange-300 border-gray-700"
		} else {
			anchorClass = "text-sky-800 hover:text-sky-900 shadow-pink-300 border-sky-800"
		}

		w.WriteString(fmt.Sprintf(`
      <span 
        data-signature="" 
        data-title="" 
        data-occurance=""
      ><a 
        href="#"
        class="font-mono font-bold drop-shadow-md tracking-normal box-decoration-clone border-b hover:border-b-2 %v">%v`, anchorClass, mention.Title()))
	} else {
		w.WriteString(`</a></span>`)
	}

	return gast.WalkContinue, nil
}

/*
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
*/
