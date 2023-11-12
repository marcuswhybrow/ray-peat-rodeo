package renderer

import (
	"fmt"
	"html/template"

	meta "github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	"github.com/mitchellh/mapstructure"
	gast "github.com/yuin/goldmark/ast"
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

func (s *SpeakerHTMLRenderer) renderSpeaker(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		speaker := node.(*ast.Speaker)

		var frontMatter meta.FrontMatter
		mapstructure.Decode(speaker.OwnerDocument().Meta(), &frontMatter)

		longName := frontMatter.Speakers[speaker.ShortName]

		t, err := template.New("openSpeaker").Parse(`
      <div 
        class="
          font-sans relative first:mt-0

          {{ if .IsHello }} 
            hello
          {{ end }}

          {{ if .IsRay }}
            is-ray ml-1 mr-16
          {{ else }} 
            ml-16 mr-1
          {{ end }}

          {{ if .IsRetorting }}
            retort mt-4
          {{ else }} 
            -mt-4 [.is-short+&>.name]
          {{ end }}
        "
      >
        {{ if not .IsRetorting }}
          <div 
            class="
              text-sm mt-8 mb-4 block

              {{ if .IsRay }} 
                text-gray-400
              {{ else }} 
                text-sky-400
              {{ end }}
            "
          >{{ .LongName }}</div>
        {{ end }}

        <div 
          class="
            p-8 rounded shadow [&>*]:mb-6 [&>*:last-child]:mb-0 [&>blockquote]:pl-4 [&>blockquote]:text-sm 

            {{ if .IsRetorting }} 
              inline-block
            {{ else }}
              block
            {{ end }}

            {{ if .IsRay }} 
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
			"LongName":    longName,
			"IsRetorting": speaker.IsRetorting(source),
			"IsRay":       speaker.IsRay(),
			"IsHello":     speaker.IsHello,
		})
	} else {
		_, _ = w.WriteString("</div></div>")
	}

	return gast.WalkContinue, nil
}
