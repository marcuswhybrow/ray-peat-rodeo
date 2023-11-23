package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

type UtteranceLink struct {
	gast.BaseInline
}

func (u *UtteranceLink) Dump(source []byte, level int) {
	gast.DumpHelper(u, source, level, nil, nil)
}

var KindUtteranceLink = gast.NewNodeKind("UtteranceLink")

func (u *UtteranceLink) Kind() gast.NodeKind {
	return KindUtteranceLink
}

func NewUtteranceLink() *UtteranceLink {
	ul := &UtteranceLink{
		BaseInline: gast.BaseInline{},
	}
	return ul
}
