package ast

import (
	"strings"

	"github.com/gosimple/slug"
	gast "github.com/yuin/goldmark/ast"
)

type Section struct {
	BaseBlock
	FileNode

	Timecode *Timecode
	Prefix   []string
	Title    string
	Level    int
}

func (s *Section) PrefixString() string {
	if len(s.Prefix) > 0 {
		return strings.Join(s.Prefix, ".") + "."
	}
	return ""
}

func (s *Section) ID() string {
	return slug.Make(s.PrefixString() + " " + s.Title)
}

func (s *Section) Dump(source []byte, level int) {
	gast.DumpHelper(s, source, level, nil, nil)
}

var KindSection = gast.NewNodeKind("Section")

func (s *Section) Kind() gast.NodeKind {
	return KindSection
}

func NewSection(title string, prefix []string, level int, timecode *Timecode) *Section {
	return &Section{
		BaseBlock: BaseBlock{},
		Timecode:  timecode,
		Prefix:    prefix,
		Title:     title,
		Level:     level,
	}
}
