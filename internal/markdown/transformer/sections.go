package transformer

import (
	"html"
	"log"
	"strings"
	"unicode"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type sectionTransformer struct{}

var SectionTransformer = &sectionTransformer{}

func NewSectionTransformer() gparser.ASTTransformer {
	return &sectionTransformer{}
}

func (t *sectionTransformer) Transform(document *gast.Document, reader text.Reader, pc gparser.Context) {

	replacements := []Replacement{}

	gast.Walk(document, func(node gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}

		if heading, ok := node.(*gast.Heading); ok {
			asset := ast.GetAsset(pc)
			log.Printf("Heading %v#%v", asset.GetPath(), string(heading.Text(asset.GetMarkdown())))

			var ultimateTimecode *ast.Timecode
			for child := heading.FirstChild(); child != nil; child = child.NextSibling() {
				if timecode, ok := child.(*ast.Timecode); ok {
					ultimateTimecode = timecode
					heading.RemoveChild(heading, timecode)
				}
			}

			source := asset.GetMarkdown()
			title := html.UnescapeString(string(heading.Text(source)))

			prefix := []string{}
			var builder strings.Builder

			if unicode.IsDigit(rune(title[0])) {
				for pos, rune := range title {
					if unicode.IsDigit(rune) || len(title) > pos+1 && title[pos+1] == '.' {
						builder.WriteRune(rune)
					} else if rune == '.' {
						prefix = append(prefix, builder.String())
						builder.Reset()
					} else {
						title = title[pos:]
						break
					}
				}
			}

			section := ast.NewSection(title, prefix, heading.Level, ultimateTimecode)

			asset.RegisterSection(section)
			replacements = append(replacements, Replacement{
				Remove: heading,
				Add:    section,
			})

			return gast.WalkSkipChildren, nil
		}

		return gast.WalkContinue, nil
	})

	for _, replacement := range replacements {
		parent := replacement.Remove.Parent()
		parent.ReplaceChild(parent, replacement.Remove, replacement.Add)
	}
}

type Replacement struct {
	Remove *gast.Heading
	Add    *ast.Section
}
