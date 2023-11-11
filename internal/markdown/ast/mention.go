package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

var KindMention = gast.NewNodeKind("Mention")

type Mention struct {
	gast.BaseInline
	IsRaySpeaking bool
	Inside        string
}

func (t *Mention) Dump(source []byte, level int) {
	gast.DumpHelper(t, source, level, nil, nil)
}

func (n *Mention) Kind() gast.NodeKind {
	return KindMention
}

func NewMention(inside string) *Mention {
	return &Mention{
		BaseInline:    gast.BaseInline{},
		IsRaySpeaking: false,
		Inside:        inside,
	}
}
