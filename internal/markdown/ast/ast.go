package ast

import (
	"github.com/mitchellh/mapstructure"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

type FrontMatter struct {
	Source struct {
		Series   string
		Title    string
		Url      string
		Kind     string
		Duration string
	}
	Speakers      map[string]string
	Transcription struct {
		Url    string
		Kind   string
		Date   string
		Author string
	}
}

var PermalinkKey = parser.NewContextKey()
var IDKey = parser.NewContextKey()
var SourceKey = parser.NewContextKey()
var HTTPCache = parser.NewContextKey()

type FrontMatterNode struct {
	gast.BaseNode
}

func (n *FrontMatterNode) FrontMatter() FrontMatter {
	var frontMatter FrontMatter
	mapstructure.Decode(n.OwnerDocument().Meta(), &frontMatter)
	return frontMatter
}

type ChatNode struct {
	FrontMatterNode
}

func (n *ChatNode) IsRaySpeaking() bool {
	for p := n.Parent(); p != nil; p = p.Parent() {
		utterance, ok := p.(*Utterance)
		if ok {
			return utterance.IsRay()
		}
	}
	return false
}

type BaseBlock struct {
	gast.BaseBlock
}

type BaseInline struct {
	gast.BaseInline
}
