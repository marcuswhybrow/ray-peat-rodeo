package ast

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	gparser "github.com/yuin/goldmark/parser"
)

var KindMention = gast.NewNodeKind("Mention")

type Mention struct {
	gast.BaseInline
	Primary   MentionPart
	Secondary MentionPart
	Label     string
	Permalink string
	Occurance int
	ID        string
	FileID    string
}

func (m *Mention) Title() string {
	if len(m.Label) > 0 {
		return m.Label
	} else {
		return m.Ultimate().PrefixFirst()
	}
}

func (m *Mention) CatalogPermalink() string {
	return "/" + m.Primary.ID() + "#" + m.FileID + "-" + m.ID
}

func (m *Mention) VignetteHTML(source []byte, radius int) string {
	for p := m.Parent(); p != nil; p = p.Parent() {
		if p.Kind() == KindSpeaker {
			before, after, _ := CutText(source, p, m, false)

			rBefore := []rune(before)
			rAfter := []rune(after)

			var result []rune
			if len(rBefore) > radius {
				result = append(result, []rune("... ")...)
			}
			result = append(result, rBefore[max(0, len(rBefore)-radius):]...)

			result = append(result, []rune(fmt.Sprintf(`<mark><a id="%v" href="%v" class="underline">%s</a></mark>`, m.FileID+"-"+m.ID, m.Permalink, m.Text(source)))...)

			result = append(result, rAfter[:min(len(rAfter), radius)]...)

			if len(rAfter) > radius {
				result = append(result, []rune(" ...")...)
			}

			return string(result)
		}
	}
	panic("Failed to find parent speaker block for mention node")
}

// Recursive search of tree for a target Node that returns the Node text before
// and after the target.
func CutText(source []byte, root gast.Node, target gast.Node, found bool) (string, string, bool) {

	if root.HasChildren() {
		left, right := "", ""
		for child := root.FirstChild(); child != nil; child = child.NextSibling() {
			if child == target {
				found = true
			} else {
				l, r, f := CutText(source, child, target, found)
				left += l
				right += r
				found = f
			}
		}

		// Block level elements (paragraphs) don't have any spaces when inline
		isBlock := root.Type() == gast.TypeBlock
		if isBlock {
			if found {
				right += " "
			} else {
				left += " "
			}
		}

		return left, right, found
	} else {
		text := html.UnescapeString(string(root.Text(source)))
		if found {
			return "", text, found
		} else {
			return text, "", found
		}
	}
}

type MentionPart struct {
	Cardinal string
	Prefix   string
	UrlTitle string
	IsURL    bool
}

func NewMentionPart(cardinal, prefix string, httpCache *cache.HTTPCache) *MentionPart {
	mp := &MentionPart{
		Cardinal: strings.Trim(cardinal, " "),
		Prefix:   strings.Trim(prefix, " "),
	}

	if len(mp.Prefix) == 0 {
		url, err := url.Parse(mp.Cardinal)
		if err == nil && url.IsAbs() {
			mp.IsURL = true

			doiHost := url.Hostname() == "doi.org"

			// All DOIs start with the number 10 followed by a period
			doiPath := strings.HasPrefix(url.Path, "/10.")

			if doiHost && doiPath {
				mp.UrlTitle = <-httpCache.GetJSON(mp.Cardinal, "title", func(res *http.Response) string {

					body, err := io.ReadAll(res.Body)
					if err != nil {
						panic(fmt.Sprintf("Failed to read body of HTTP response for url '%v': %v", mp.Cardinal, err))
					}

					data := DOIData{}
					err = json.Unmarshal(body, &data)
					if err != nil {

						panic(fmt.Sprintf("Failed to unmarshal JSON response for url '%v': %v", mp.Cardinal, err))
					}

					return data.Title
				})
			} else {
				mp.UrlTitle = <-httpCache.Get(mp.Cardinal, "title", cache.H1OrTitleHandler)
			}
		}
	}

	return mp
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

func (p *MentionPart) ID() string {
	id := p.CardinalFirst()
	id = strings.ToLower(id)
	id = strings.ReplaceAll(id, " ", "-")
	return id
}

func (m *Mention) Dump(source []byte, level int) {
	gast.DumpHelper(m, source, level, nil, nil)
}

func (m *Mention) Text(source []byte) []byte {
	if len(m.Label) > 0 {
		return []byte(m.Label)
	} else {
		return []byte(m.Ultimate().PrefixFirst())
	}
}

func (m *Mention) Kind() gast.NodeKind {
	return KindMention
}

var perMentionCountKey = parser.NewContextKey()

func NewMention(pc gparser.Context, primary, secondary MentionPart, displayText string) *Mention {
	counts := pc.ComputeIfAbsent(perMentionCountKey, func() interface{} {
		return map[MentionPart]int{}
	}).(map[MentionPart]int)
	count := counts[primary]
	count += 1
	counts[primary] = count
	pc.Set(perMentionCountKey, counts)

	filePermalink := pc.Get(markdown.PermalinkKey).(string)
	id := primary.ID()
	if count > 1 {
		id += "-" + fmt.Sprint(count)
	}

	permalink := filePermalink + "#" + id

	return &Mention{
		BaseInline: gast.BaseInline{},
		Primary:    primary,
		Secondary:  secondary,
		Label:      displayText,
		ID:         id,
		FileID:     pc.Get(markdown.IDKey).(string),
		Permalink:  permalink,
		Occurance:  count,
	}
}

func (m *Mention) Ultimate() *MentionPart {
	if len(m.Secondary.Cardinal) > 0 || len(m.Secondary.Prefix) > 0 {
		return &m.Secondary
	} else {
		return &m.Primary
	}
}

type DOIData struct {
	Title string
}
