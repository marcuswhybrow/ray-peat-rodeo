package ast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
)

type MentionablePart struct {
	Cardinal string
	Prefix   string
	URLTitle string
}

func NewMentionablePart(c, p string, httpCache *cache.HTTPCache) MentionablePart {
	m := MentionablePart{
		Cardinal: strings.Trim(c, " "),
		Prefix:   strings.Trim(p, " "),
	}
	m.URLTitle = m.urlTitle(httpCache)
	return m
}

func (m *MentionablePart) HasPrefix() bool {
	return len(m.Prefix) > 0
}

func (m *MentionablePart) ID() string {
	id := m.Cardinal
	if m.HasPrefix() {
		id = m.Prefix + "-" + m.Cardinal
	}
	id = strings.ToLower(id)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "'", "")
	id = url.QueryEscape(id)
	return id
}

func (m *MentionablePart) IsURL() bool {
	if m.HasPrefix() {
		return false
	}

	u, err := url.Parse(m.Cardinal)
	if err != nil || !u.IsAbs() {
		return false
	}
	return true
}

// If MentionablePart is a URL, returns a sensible title. Results are cached
// and persisted under source control.
func (m MentionablePart) urlTitle(httpCache *cache.HTTPCache) string {
	if m.HasPrefix() {
		return ""
	}

	url, err := url.Parse(m.Cardinal)
	if err != nil || !url.IsAbs() {
		return ""
	}

	isDOIDotOrg := url.Hostname() == "doi.org"

	// All DOIs start with the number 10 followed by a period
	isDOIPath := strings.HasPrefix(url.Path, "/10.")

	if isDOIDotOrg && isDOIPath {
		title := <-httpCache.GetJSON(m.Cardinal, "title", func(res *http.Response) string {

			body, err := io.ReadAll(res.Body)
			if err != nil {
				panic(fmt.Sprintf("Failed to read body of HTTP response for url '%v': %v", m.Cardinal, err))
			}

			data := DOIData{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				panic(fmt.Sprintf("Failed to unmarshal JSON response for url '%v': %v", m.Cardinal, err))
			}

			return data.Title
		})
		return title
	} else {
		// For normal URLs get the first H1 text or, failing that, the html title
		title := <-httpCache.Get(m.Cardinal, "title", cache.GetH1OrTitle)
		return title
	}
}

func (p *MentionablePart) IsEmpty() bool {
	return len(p.Cardinal) == 0 && len(p.Prefix) == 0
}

func (p *MentionablePart) PrefixFirst() string {
	return strings.Trim(fmt.Sprintf("%v %v", p.Prefix, p.Cardinal), " ")
}

func (p *MentionablePart) CardinalFirst() string {
	result := p.Cardinal
	if len(p.Prefix) > 0 {
		result += ", " + p.Prefix
	}
	return result
}

func (p *MentionablePart) ParseUrl() (*url.URL, error) {
	if len(p.Prefix) > 0 {
		return nil, nil
	}
	return url.Parse(p.Cardinal)
}

type DOIData struct {
	Title string
}
