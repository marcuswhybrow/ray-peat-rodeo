package transformer

import (
	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type utteranceLinkTransformer struct{}

var UtteranceLinkTransformer = &utteranceLinkTransformer{}

func NewUteranceLinkTransformer() gparser.ASTTransformer {
	return &utteranceLinkTransformer{}
}

func (t *utteranceLinkTransformer) Transform(document *gast.Document, reader text.Reader, pc gparser.Context) {
	gast.Walk(document, func(node gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}

		if link, ok := node.(*gast.Link); ok {
			link.SetAttribute([]byte("class"), []uint8("border-b hover:border-b-2 border-gray-400"))
		}

		return gast.WalkContinue, nil
	})
}
