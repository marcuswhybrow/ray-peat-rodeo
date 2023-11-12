package renderer

import (
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
)

func IsRaySpeaking(node gast.Node) bool {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		speaker, ok := parent.(*ast.Speaker)
		if ok {
			return speaker.IsRay()
		}
	}

	return false
}
