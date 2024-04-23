package parser

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strconv"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type TimecodeParser struct{}

func NewTimecodeParser() *TimecodeParser {
	return &TimecodeParser{}
}

func (w *TimecodeParser) Trigger() []byte {
	return []byte{'['}
}

func (w *TimecodeParser) Parse(parent gmAst.Node, block text.Reader, pc parser.Context) gmAst.Node {
	line, _ := block.PeekLine()

	i := bytes.Index(line, []byte{']'})
	if i < 2 {
		return nil
	}

	consumed := i + 1

	sections := bytes.Split(line[1:i], []byte{':'})
	timecode := ast.NewTimecode()
	timecode.Source = string(line[:i])

	n := len(sections)
	if n > 3 || n < 2 {
		return nil
	}

	seconds, err := strconv.Atoi(string(sections[n-1]))
	if err != nil {
		return nil
	}
	timecode.Seconds = seconds

	minutes, err := strconv.Atoi(string(sections[n-2]))
	if err != nil {
		return nil
	}
	timecode.Minutes = minutes

	var hours int
	if n < 3 {
		hours = 0
	} else {
		hours, err = strconv.Atoi(string(sections[n-3]))
		if err != nil {
			return nil
		}
	}
	timecode.Hours = hours

	asset := ast.GetAsset(pc)

	sourceURLStr := asset.GetSourceURL()
	timecode.ExternalURL = externalUrl(sourceURLStr, timecode.Seconds, timecode.Minutes, timecode.Hours)

	block.Advance(consumed)

	return timecode
}

func externalUrl(sourceURLStr string, seconds, minutes, hours int) string {
	sourceUrl, err := url.Parse(sourceURLStr)
	if err != nil {
		log.Panicf("Failed to parse frontmatter source url: %v", err)
	}

	newUrl := *sourceUrl

	var timecode string
	isYouTube := slices.Contains([]string{
		"www.youtube.com",
		"youtube.com",
		"youtu.be",
	}, sourceUrl.Hostname())

	if isYouTube {
		// Youtube timecodes: ?t=1h12m32s
		if hours == 0 && minutes == 0 {
			timecode = fmt.Sprintf("%ds", seconds)
		} else if hours == 0 {
			timecode = fmt.Sprintf("%dm%ds", minutes, seconds)
		} else {
			timecode = fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
		}

		query := newUrl.Query()
		query.Del("t")
		query.Add("t", timecode)
		newUrl.RawQuery = query.Encode()
	} else {
		// Everyone else: #t=01:12:32
		newUrl.Fragment = fmt.Sprintf("t=%02d:%02d:%02d", hours, minutes, seconds)
	}

	return newUrl.String()
}
