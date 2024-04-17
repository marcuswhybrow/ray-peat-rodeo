package ast

import (
	"fmt"
	"net/url"
	"slices"

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

func (t *Timecode) ExternalUrl() (*url.URL, error) {
	sourceUrl, err := url.Parse(t.FrontMatter().Source.Url)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse frontmatter source url: %v", err)
	}

	newUrl := *sourceUrl

	var timecode string
	if slices.Contains([]string{
		"www.youtube.com",
		"youtube.com",
		"youtu.be",
	}, sourceUrl.Hostname()) {

		// Youtube timecodes: ?t=1h12m32s
		if t.Hours == 0 && t.Minutes == 0 {
			timecode = fmt.Sprintf("%ds", t.Seconds)
		} else if t.Hours == 0 {
			timecode = fmt.Sprintf("%dm%ds", t.Minutes, t.Seconds)
		} else {
			timecode = fmt.Sprintf("%dh%dm%ds", t.Hours, t.Minutes, t.Seconds)
		}

		query := newUrl.Query()
		query.Del("t")
		query.Add("t", timecode)
		newUrl.RawQuery = query.Encode()

	} else {
		// Everyone else: #t=01:12:32
		newUrl.Fragment = fmt.Sprintf("t=%02d:%02d:%02d", t.Hours, t.Minutes, t.Seconds)
	}

	return &newUrl, nil
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
