package ast

import (
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

type Speaker struct {
	gast.BaseBlock

	// The short name of a speaker, e.g. "RP", used in markdown.
	ShortName string

	// Is speaking for the first time in this chat.
	IsHello bool

	// A short reponse for which it isn't worth redeclaring the speakers name if
	// the previous was directly interrupting this speaker.
	IsRetorting bool
}

func NewSpeaker() *Speaker {
	return &Speaker{
		BaseBlock: gast.BaseBlock{},
	}
}

func (s *Speaker) Dump(source []byte, level int) {
	gast.DumpHelper(s, source, level, nil, nil)
}

var KindSpeaker = gast.NewNodeKind("Speaker")

func (s *Speaker) Kind() gast.NodeKind {
	return KindSpeaker
}

func (s *Speaker) IsRay() bool {
	return strings.Trim(s.ShortName, " ") == "RP"
}

func IsRaySpeaking(node gast.Node) bool {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		speaker, ok := parent.(*Speaker)
		if ok {
			return speaker.IsRay()
		}
	}

	return false
}
