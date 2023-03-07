package speakers

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type speakerParser struct {
}

func NewParser() parser.BlockParser {
	return &speakerParser{}
}

func (s *speakerParser) Trigger() []byte {
	return []byte{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}
}

func getShortName(line []byte) ([]byte, int) {
	i := 0
	for ; i < len(line); i++ {
		if line[i] == ':' {
			return line[0:i], i + 1
		}
		if !util.IsAlphaNumeric(line[i]) {
			return nil, 0
		}
	}
	return nil, 0
}

func (s *speakerParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {

	shortName, bytesConsumed := func() ([]byte, int) {
		line, _ := reader.PeekLine()
		return getShortName(line)
	}()
	if len(shortName) <= 0 {
		return nil, parser.HasChildren | parser.Continue
	}

	speakerParagraph := func() ast.Node {
		node := NewSpeakerParagraph()
		node.shortName = string(shortName)
		return node
	}()
	reader.Advance(bytesConsumed)
	return speakerParagraph, parser.HasChildren | parser.Continue
}

func (s *speakerParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	shortName, _ := getShortName(line)
	if shortName != nil {
		return parser.Close
	}
	return parser.HasChildren | parser.Continue
}

func (s *speakerParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {

}

func (s *speakerParser) CanInterruptParagraph() bool {
	return false
}

func (s *speakerParser) CanAcceptIndentedLine() bool {
	return false
}

type SpeakerParagraph struct {
	ast.BaseBlock
	shortName string
	longName  string
}

func (n *SpeakerParagraph) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindSpeakerParagraph = ast.NewNodeKind("Paragraph")

func (n *SpeakerParagraph) Kind() ast.NodeKind {
	return KindSpeakerParagraph
}

func NewSpeakerParagraph() *SpeakerParagraph {
	return &SpeakerParagraph{
		BaseBlock: ast.BaseBlock{},
	}
}

type SpeakerParagraphHTMLRenderer struct {
	html.Config
}

func NewSpeakerParagraphHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &SpeakerParagraphHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

func (r *SpeakerParagraphHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindSpeakerParagraph, r.renderSpeakerParagraph)
}

var SpeakerParagraphAttributeFilter = html.GlobalAttributeFilter

func (s *SpeakerParagraphHTMLRenderer) renderSpeakerParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		sp := node.(*SpeakerParagraph)
		longName := func() string {
			if sp.shortName == "RP" {
				return "Ray Peat"
			} else if sp.shortName == "Host" {
				return "Host"
			}
			speakers, speakersOk := sp.OwnerDocument().Meta()["speakers"].(map[interface{}]interface{})
			if speakersOk {
				longName, longNameOk := speakers[sp.shortName].(string)
				if longNameOk {
					return longName
				}
			}
			panic(fmt.Sprintf("'speakers.%s' not defined in frontmatter", sp.shortName))
		}()
		w.WriteString(fmt.Sprintf(`
			<div data-shortname="%s" data-longname="%s" class="speaker">
			<span class="speaker-name">%s:</span>
		`, sp.shortName, longName, longName))
	} else {
		_, _ = w.WriteString("</div>\n")
	}
	return ast.WalkContinue, nil
}

type speakers struct {
}

// Speakers is an extension for Goldmark that prefixes paragraphs with exclamation marks
var Speakers = &speakers{}

func New() goldmark.Extender {
	return &speakers{}
}

func (e *speakers) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewParser(), 100),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewSpeakerParagraphHTMLRenderer(), 100),
	))
}
