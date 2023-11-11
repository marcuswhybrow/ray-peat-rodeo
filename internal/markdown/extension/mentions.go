package extension

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/parser"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/renderer"
	"github.com/yuin/goldmark"
	gparser "github.com/yuin/goldmark/parser"
	grenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var mentionsContextKey = gparser.NewContextKey()

/*
type MentionSignature struct {
	Cardinal string
	Prefix   string
}

type Mention struct {
	Signature  MentionSignature
	SubMention SubMention
}

type SubMention struct {
	Signature MentionSignature
}

func Get(pc gparser.Context) []Mention {
	return pc.Get(mentionsContextKey).([]Mention)
}
*/

type mentions struct{}

// Mentions is an extension for Goldmark that replaces [[mention]] with a link
var Mentions = &mentions{}

func NewMentions() goldmark.Extender {
	return &mentions{}
}

func (e *mentions) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(gparser.WithInlineParsers(
		util.Prioritized(parser.NewMentionParser(), 1)),
	)
	m.Renderer().AddOptions(grenderer.WithNodeRenderers(
		util.Prioritized(renderer.NewMentionHTMLRenderer(), 1)),
	)
}
