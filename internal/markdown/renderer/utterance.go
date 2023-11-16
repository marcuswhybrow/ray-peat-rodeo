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
	reg.Register(ast.KindSpeaker, r.renderSpeaker)
}

var SpeakerAttributeFilter = html.GlobalAttributeFilter

func (s *UtteranceHTMLRenderer) renderSpeaker(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		utterance := node.(*ast.Utterance)

		// It's pleasant to visually squish many small replies together if it's
		// obvious who's speaking by the context alone.
		isShort := func() bool {
			if utterance.IsNewSpeaker {
				return false
			}

			if !utterance.PrevAndNextIsSameSpeaker() {
				return false
			}

			if len(utterance.Text(source)) > 50 {
				return false
			}

			// Utterances by other speakers may only be squished if Ray was
			// previously speaking. That's because all utterances by other speakers
			// are visually indistinguisable when squished.
			if !utterance.IsRay() {
				return utterance.PrevIsRay()
			}

			// However Ray's utterances are always visually distinguishable from any
			// utterance by another speaker.
			return true
		}()

		t, err := template.New("openSpeaker").Parse(`
      <div 
        class="
          font-sans relative first:mt-0

          {{ if .Utterance.IsNewSpeaker }} 
            hello
          {{ end }}

          {{ if .Utterance.IsRay }}
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
              text-sm mt-8 mb-4 block

              {{ if .Utterance.IsRay }} 
                text-gray-400
              {{ else }} 
                text-sky-400
              {{ end }}
            "
          >{{ .Utterance.SpeakerName }}</div>
        {{ end }}

        <div 
          class="
            p-8 rounded shadow [&>*]:mb-6 [&>*:last-child]:mb-0 [&>blockquote]:pl-4 [&>blockquote]:text-sm 

            {{ if .IsShort }} 
              inline-block
            {{ else }}
              block
            {{ end }}

            {{ if .Utterance.IsRay }} 
              text-gray-900 bg-gray-100
            {{ else }} 
              text-sky-900 bg-gradient-to-br from-sky-100 to-blue-200
            {{ end }}
          "
        >
    `)

		if err != nil {
			return gast.WalkStop, fmt.Errorf("Failed to parse speaker html/template: %v", err)
		}

		t.Execute(w, map[string]interface{}{
			"IsShort":   isShort,
			"Utterance": utterance,
		})
	} else {
		_, _ = w.WriteString("</div></div>")
	}

	return gast.WalkContinue, nil
}
