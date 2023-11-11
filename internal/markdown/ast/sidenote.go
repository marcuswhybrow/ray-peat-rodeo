package ast

import gast "github.com/yuin/goldmark/ast"

type Sidenote struct {
	gast.BaseInline
	Position int
}

func (t *Sidenote) Dump(source []byte, level int) {
	gast.DumpHelper(t, source, level, nil, nil)
}

var KindSidenote = gast.NewNodeKind("Sidenote")

func (n *Sidenote) Kind() gast.NodeKind {
	return KindSidenote
}

func NewSidenote(position int) *Sidenote {
	return &Sidenote{
		BaseInline: gast.BaseInline{},
		Position:   position,
	}
}
