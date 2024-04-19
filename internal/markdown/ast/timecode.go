package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

type Timecode struct {
	BaseInline
	FileNode

	Hours       int
	Minutes     int
	Seconds     int
	Source      string
	ExternalURL string
	Permalink   string
}

func (t *Timecode) Terse() string {
	if t.Hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", t.Hours, t.Minutes, t.Seconds)
	}

	return fmt.Sprintf("%02d:%02d", t.Minutes, t.Seconds)
}

func (t *Timecode) Dump(source []byte, level int) {
	gast.DumpHelper(t, source, level, nil, nil)
}

var KindTimecode = gast.NewNodeKind("Timecode")

func (n *Timecode) Kind() gast.NodeKind {
	return KindTimecode
}

func NewTimecode() *Timecode {
	return &Timecode{
		BaseInline: BaseInline{},
	}
}
