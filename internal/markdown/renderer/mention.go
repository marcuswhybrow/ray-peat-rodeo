package renderer

import (
	"fmt"
	"html/template"

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
		label := mention.Label
		if len(label) == 0 {
			label = mention.Mentionable.Ultimate().PrefixFirst()
		}

		t, err := template.New("OpenMention").Parse(
			`<span
        hx-trigger="mouseenter"
        hx-target="find .popup"
        hx-get="{{ .Mention.Mentionable.PopupPermalink }}"
        hx-swap="innerHTML"
        hx-select=".hx-select"
        data-mention-id="{{ .Mention.ID }}"
        class="relative"
        _="
          on mouseenter
            send open to .popup in me

          on mouseleave
            wait for open(elem) or 200ms
            if the result's type is not 'open'
              send close to .popup in me
            end
        "
      ><a 
        id="{{ .Mention.LocalID }}"
        href="{{ .Mention.Mentionable.Permalink }}"
        class="
          font-mono font-bold tracking-normal 
          drop-shadow-md box-decoration-clone 
          border-b hover:border-b-2

          {{ if .IsRaySpeaking }}
			      text-gray-700 hover:text-gray-900 shadow-orange-300 border-gray-700
          {{ else }}
			      text-sky-800 hover:text-sky-900 shadow-pink-300 border-sky-800
          {{ end }}
        "
      >{{ .Label }}</a><span
        class="
          popup 
          bg-white rounded shadow-lg block absolute 
          z-50 
          overflow-y-auto 
          mb-4 
          transition-[opacity]
          left-[calc(50%-200px)]
          top-[120%]
        "
        data-mention-id="{{ .Mention.ID }}"
        _="
          on open
            set my *width to 400px
            set my *height to 400px
            send checkPosition to me
            transition my *opacity to 1.0 over 50ms then settle

          on checkPosition
            measure my left then
            if left < 0 
              set my *left to 0
            else 
              measure my right then
              if right > max
                set my *right to 0
              end
            end

          on close 
            transition my *opacity to 0 over 50ms then settle
            set my *width to 0
            set my *height to 0
        "
      ></span>`,
		)

		if err != nil {
			panic(fmt.Sprintf("Failed to parse mention html/template: %v", err))
		}

		t.Execute(w, map[string]interface{}{
			"Mention":       mention,
			"IsRaySpeaking": ast.IsRaySpeaking(mention),
			"Label":         label,
		})
	} else {
		w.WriteString(`</span>`)
	}

	return gast.WalkContinue, nil
}
