package ast

import (
	"fmt"
	"net/url"
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

var KindMention = gast.NewNodeKind("Mention")

type Mention struct {
	gast.BaseInline
	Primary     MentionPart
	Secondary   MentionPart
	DisplayText string
}

func (m *Mention) Title() string {
	if len(m.DisplayText) > 0 {
		return m.DisplayText
	} else {
		return m.Ultimate().PrefixFirst()
	}
}

type MentionPart struct {
	Cardinal string
	Prefix   string
}

func (p *MentionPart) PrefixFirst() string {
	return strings.Trim(fmt.Sprintf("%v %v", p.Prefix, p.Cardinal), " ")
}

func (p *MentionPart) CardinalFirst() string {
	result := p.Cardinal
	if len(p.Prefix) > 0 {
		result += ", " + p.Prefix
	}
	return result
}

func (p *MentionPart) ParseUrl() (*url.URL, error) {
	if len(p.Prefix) > 0 {
		return nil, nil
	}
	return url.Parse(p.Cardinal)
}

func (t *Mention) Dump(source []byte, level int) {
	gast.DumpHelper(t, source, level, nil, nil)
}

func (n *Mention) Kind() gast.NodeKind {
	return KindMention
}

func NewMention(primary, secondary MentionPart, displayText string) *Mention {
	return &Mention{
		BaseInline:  gast.BaseInline{},
		Primary:     primary,
		Secondary:   secondary,
		DisplayText: displayText,
	}
}

func (m *Mention) Ultimate() *MentionPart {
	if len(m.Secondary.Cardinal) > 0 || len(m.Secondary.Prefix) > 0 {
		return &m.Secondary
	} else {
		return &m.Primary
	}
}
