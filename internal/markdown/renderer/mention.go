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
        hx-trigger="mouseenter once"
        hx-target="find .popup"
        hx-get="{{ .Mention.Mentionable.PopupPermalink }}"
        hx-swap="innerHTML"
        hx-select=".hx-select"
        data-mention-id="{{ .Mention.ID }}"
        class="mention relative cursor-pointer"
        onclick="void(0)"
        _="
          on mousemove
            trigger closeAllPopups(exception: my @data-mention-id) on .open-popup
            trigger openPopup on .popup in me

          on closeAllPopups(exception)
            if my @data-mention-id is not the exception 
              trigger closePopup on .popup in me
            end

          on mouseleave
            wait for closeAllPopups(exception) or openPopup or 500ms
            if the result's type is 'closeAllPopups'
              if my @data-mention-id is not the exception
                trigger closePopup on .popup in me
              end
            else if the result's type is 'openPopup'
            else 
              trigger closePopup on .popup in me
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
          bg-white shadow-2xl block absolute 
          z-20 
          overflow-hidden
          overflow-y-auto 
          mb-4 
          w-[400px] h-[300px]
          left-[calc(50%-200px)]
          hidden
          top-[120%]
          scrollbar
          scrollbar-track-slate-100
          scrollbar-thumb-slate-200
        "
        data-mention-id="{{ .Mention.ID }}"
        data-pagefind-ignore
        _="
          on openPopup
            add .open-popup to me
            remove .hidden from me
            send checkPosition to me

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

          on closePopup
            remove .open-popup from me
            add .hidden to me
        "
      >
        <span class="text-center text-gray-400 block p-8">
          loading {{ .Mention.Mentionable.Ultimate.CardinalFirst }}...
        </span>
      </span>`,
		)

		if err != nil {
			panic(fmt.Sprintf("Failed to parse mention html/template: %v", err))
		}

		t.Execute(w, map[string]interface{}{
			"Mention":       mention,
			"IsRaySpeaking": ast.IsPrimarySpeaker(mention),
			"Label":         label,
		})
	} else {
		w.WriteString(`</span>`)
	}

	return gast.WalkContinue, nil
}
