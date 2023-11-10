package renderer

import (
	"fmt"

	meta "github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	"github.com/mitchellh/mapstructure"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

const RETORT_MAX_LEN = 50

type SpeakerHTMLRenderer struct {
	html.Config
}

func NewSpeakerHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &SpeakerHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

func (r *SpeakerHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindSpeaker, r.renderSpeaker)
}

var SpeakerAttributeFilter = html.GlobalAttributeFilter

func (s *SpeakerHTMLRenderer) renderSpeaker(w util.BufWriter, source []byte, node gmAst.Node, entering bool) (gmAst.WalkStatus, error) {
	if entering {
		speaker := node.(*ast.Speaker)

		var frontMatter meta.FrontMatter
		mapstructure.Decode(speaker.OwnerDocument().Meta(), &frontMatter)

		longName := frontMatter.Speakers[speaker.ShortName]

		wrapperClass := "font-sans relative first:mt-0"

		if speaker.IsHello {
			wrapperClass += " hello"
		}

		if speaker.IsRay() {
			wrapperClass += " is-ray ml-1 mr-16" // tailwind CSS classes
		} else {
			wrapperClass += " ml-16 mr-1"
		}

		if speaker.IsRetorting(source) {
			wrapperClass += " retort mt-4"
		} else {
			wrapperClass += " -mt-4 [.is-short+&>.name]"
		}

		w.WriteString(fmt.Sprintf(`<div class="%v">`, wrapperClass))

		if !speaker.IsRetorting(source) {
			nameClass := "text-sm mt-8 mb-4 block"

			if speaker.IsRay() {
				nameClass += " text-gray-400"
			} else {
				nameClass += " text-sky-400"
			}

			w.WriteString(fmt.Sprintf(`<div class="%v">%v</div>`, nameClass, longName))
		}

		innerClass := "p-8 rounded shadow [&>*]:mb-6 [&>*:last-child]:mb-0 [&>blockquote]:pl-4 [&>blockquote]:text-sm"

		if speaker.IsRetorting(source) {
			innerClass += " inline-block"
		} else {
			innerClass += " block"
		}

		if speaker.IsRay() {
			innerClass += " text-gray-900 bg-gray-100"
		} else {
			innerClass += " text-sky-900 bg-gradient-to-br from-sky-100 to-blue-200"
		}

		w.WriteString(fmt.Sprintf(`<div class="name %v">`, innerClass))
	} else {
		_, _ = w.WriteString("</div></div>")
	}

	return gmAst.WalkContinue, nil
}
