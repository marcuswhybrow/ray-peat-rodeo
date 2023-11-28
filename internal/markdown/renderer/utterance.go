package renderer

import (
	"fmt"
	"html/template"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

const RETORT_MAX_LEN = 50

type UtteranceHTMLRenderer struct {
	html.Config
}

func NewUtteranceHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &UtteranceHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

func (r *UtteranceHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindUtterance, r.renderSpeaker)
}

var SpeakerAttributeFilter = html.GlobalAttributeFilter

func isShort(utterance *ast.Utterance, source []byte) bool {
	speakerId := utterance.Speaker.GetID()

	if utterance.IsNewSpeaker {
		return false
	}

	if len(utterance.Text(source)) > 50 {
		return false
	}

	if !utterance.Prev().IsSandwichedBetween(speakerId) {
		return false
	}

	return utterance.Speaker.GetIsPrimarySpeaker() != utterance.PrevIsPrimarySpeaker()
}

func showAvatar(utterance *ast.Utterance, source []byte) bool {
	if utterance.IsNewSpeaker {
		return true
	}

	return len(utterance.Text(source)) > 50
}

func (s *UtteranceHTMLRenderer) renderSpeaker(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		utterance := node.(*ast.Utterance)

		t, err := template.New("openUtterance").Parse(`
      <div 
        title="{{ .Utterance.Speaker.GetName }}"
        class="
          utterance
          font-sans relative first:mt-0 relative

          {{ if .Utterance.Speaker.GetIsPrimarySpeaker }}
            is-ray ml-1 mr-16
          {{ else }} 
            ml-16 mr-1
          {{ end }}

          {{ if .IsShort }}
            retort -mt-4 [.is-short+&>.name]
          {{ else }} 
            mt-4
          {{ end }}
        "
      >
        {{ if not .IsShort }}
          <div 
            class="
              speaker-name
              text-sm mt-8 mb-4 block

              {{ if .Utterance.Speaker.GetIsPrimarySpeaker }} 
                text-gray-400
              {{ else }} 
                text-sky-400
              {{ end }}
            "
          >{{ .Utterance.Speaker.GetName }}</div>
        {{ end }}

        <div 
          class="
            utterance-body
            p-8 rounded shadow [&>p]:mb-6 [&>*:last-child]:mb-0 [&>blockquote]:pl-4 [&>blockquote]:text-sm 

            {{ if .IsShort }} 
              inline-block
            {{ else }}
              {{ if and .ShowAvatar .Utterance.Speaker.GetAvatarPath }}
                min-h-[9rem]
              {{ end }}
              block
            {{ end }}

            {{ if .Utterance.Speaker.GetIsPrimarySpeaker }} 
              text-gray-900 bg-gray-100
            {{ else }} 
              text-sky-900 bg-gradient-to-br from-sky-100 to-blue-200
            {{ end }}
          "
        >
          {{ if not .IsShort }}
            {{ if .ShowAvatar }}
              {{ if .Utterance.Speaker.GetAvatarPath }}
                <div class="speaker-avatar w-16 h-20 rounded-lg inline-bock shadow float-left mr-4 mb-0 overflow-hidden">
                  <div class="w-[9999px]">
                    <img
                      class="h-20"
                      src="{{ .Utterance.Speaker.GetAvatarPath }}"
                      alt="{{ .Utterance.Speaker.GetName }}"
                    />
                  </div>
                </div>
              {{ end }}
            {{ end }}
          {{ end }}
    `)

		if err != nil {
			return gast.WalkStop, fmt.Errorf("Failed to parse speaker html/template: %v", err)
		}

		t.Execute(w, map[string]interface{}{
			"IsShort":    isShort(utterance, source),
			"ShowAvatar": showAvatar(utterance, source),
			"Utterance":  utterance,
		})
	} else {
		_, _ = w.WriteString("</div></div>")
	}

	return gast.WalkContinue, nil
}
