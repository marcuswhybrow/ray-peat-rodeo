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
      <rpr-utterance 
        {{ if .Utterance.Speaker.GetName }}by="{{ .Utterance.Speaker.GetName }}"{{ end }}
        {{ if .Utterance.Speaker.GetAvatarPath }}avatar="{{ .Utterance.Speaker.GetAvatarPath }}"{{ end }}
        {{ if .Utterance.Speaker.GetIsPrimarySpeaker }}primary="true"{{ end }}
        {{ if .IsShort }}short="true"{{ end }}
      >
    `)

		if err != nil {
			return gast.WalkStop, fmt.Errorf("Failed to parse speaker html/template: %v", err)
		}

		t.Execute(w, map[string]interface{}{
			"IsShort":   isShort(utterance, source),
			"Utterance": utterance,
		})
	} else {
		// _, _ = w.WriteString("</div></div>")
		_, _ = w.WriteString("</rpr-utterance>")
	}

	return gast.WalkContinue, nil
}
