package transformer

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
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
			parent := link.Parent()
			utteranceLink := ast.NewUtteranceLink()
			parent.InsertAfter(parent, link, utteranceLink)
			utteranceLink.AppendChild(utteranceLink, link)
		}

		return gast.WalkContinue, nil
	})
}
