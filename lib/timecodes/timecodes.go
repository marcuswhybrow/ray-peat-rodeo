package timecodes

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type timecodesParser struct {
}

func NewParser() parser.InlineParser {
	return &timecodesParser{}
}

func (w *timecodesParser) Trigger() []byte {
	return []byte{'['}
}

func (w *timecodesParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	sections, bytesConsumed := func() ([][]byte, int) {
		i := 1
		results := [][]byte{}
		colonCount := 0
		head := i
		for ; i < len(line); i++ {
			if line[i] == ':' {
				if colonCount >= 2 {
					break
				}
				results = append(results, line[head:i])
				head = i + 1
			} else if line[i] == ']' {
				if i == head {
					break
				}
				results = append(results, line[head:i])
				break
			} else if !util.IsNumeric(line[i]) {
				break
			}
		}
		return results, i + 1
	}()

	sectionCount := len(sections)
	if sectionCount == 0 || sectionCount > 3 {
		return nil
	}

	parseInt := func(bytes []byte) int {
		i, err := strconv.Atoi(string(bytes))
		if err != nil {
			panic(fmt.Sprintf("Timecode contains non-numeric characters:\n%s", line))
		}
		return i
	}

	timecode := func() *Timecode {
		t := NewTimecode()
		t.seconds = parseInt(sections[sectionCount-1])
		t.minutes = func() int {
			if sectionCount >= 2 {
				return parseInt(sections[sectionCount-2])
			} else {
				return 0
			}
		}()
		t.hours = func() int {
			if sectionCount >= 3 {
				return parseInt(sections[sectionCount-3])
			} else {
				return 0
			}
		}()
		return t
	}()
	block.Advance(bytesConsumed)
	return timecode
}

type Timecode struct {
	ast.BaseInline
	hours   int
	minutes int
	seconds int
	url     string
}

func (t *Timecode) Dump(source []byte, level int) {
	ast.DumpHelper(t, source, level, nil, nil)
}

var KindTimecode = ast.NewNodeKind("Timecode")

func (n *Timecode) Kind() ast.NodeKind {
	return KindTimecode
}

func NewTimecode() *Timecode {
	return &Timecode{
		BaseInline: ast.BaseInline{},
	}
}

type TimecodeHTMLRenderer struct {
}

func NewTimecodeHTMLRenderer() renderer.NodeRenderer {
	return &TimecodeHTMLRenderer{}
}

func (r *TimecodeHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindTimecode, r.renderTimecode)
}

func (t *TimecodeHTMLRenderer) renderTimecode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		timecode := node.(*Timecode)
		sourceUrl := func() *url.URL {
			source, sourceOk := timecode.OwnerDocument().Meta()["source"].(string)
			if !sourceOk {
				// TODO Include filename in which error occured
				panic("Timecode defined without 'source' defined in frontmatter")
			}
			u, err := url.Parse(source)
			if err != nil {
				panic(fmt.Sprintf("'source' defined in frontmatter cannot be parsed as a URL: %s", source))
			}
			return u
		}()

		timecodeUrl := func() string {
			hostname := sourceUrl.Hostname()
			if strings.Contains(hostname, "youtube.com") || strings.Contains(hostname, "youtu.be") {
				sourceUrl.Fragment = fmt.Sprintf("t=%02dh%02dm%02ds", timecode.hours, timecode.minutes, timecode.seconds)
			} else {
				sourceUrl.Fragment = fmt.Sprintf("t=%02d:%02d:%02d", timecode.hours, timecode.minutes, timecode.seconds)
			}
			return sourceUrl.String()
		}()

		text := func() string {
			if timecode.hours > 0 {
				return fmt.Sprintf("%02d:%02d:%02d", timecode.hours, timecode.minutes, timecode.seconds)
			} else {
				return fmt.Sprintf("%02d:%02d", timecode.minutes, timecode.seconds)
			}
		}()
		w.WriteString(fmt.Sprintf(`
			<span class="timecode">
				<a class="internal" href="#t=%s">
					<svg width="16" height="16" version="1.1" viewBox="0 0 383.028 383.027">
						<path d="M361.213,244.172l-71.073-71.073c-16.042-16.042-37.632-23.216-58.648-21.562c1.653-21.019-5.521-42.609-21.563-58.651
							l-71.073-71.073c-29.084-29.084-76.408-29.083-105.492,0L21.814,33.361c-29.084,29.084-29.084,76.408,0,105.493l71.073,71.073
							c16.042,16.042,37.632,23.217,58.651,21.563c-1.654,21.02,5.52,42.607,21.563,58.65l71.073,71.073
							c29.084,29.084,76.408,29.083,105.492,0l11.548-11.548C390.297,320.58,390.297,273.256,361.213,244.172z M136.174,161.292
							l29.458,29.458c-14.997,8.932-34.734,6.955-47.629-5.94l-71.073-71.073c-15.233-15.234-15.233-40.022,0-55.258l11.549-11.548
							c15.235-15.235,40.023-15.235,55.259,0l71.072,71.073c12.896,12.895,14.873,32.632,5.94,47.63l-29.458-29.458
							c-6.937-6.937-18.181-6.937-25.117,0S129.238,154.354,136.174,161.292z M336.095,324.547l-11.548,11.548
							c-15.234,15.235-40.022,15.234-55.258,0l-71.073-71.073c-12.895-12.895-14.873-32.632-5.938-47.629l29.458,29.458
							c6.936,6.938,18.181,6.938,25.116,0c6.937-6.937,6.938-18.181,0-25.115l-29.458-29.459c14.998-8.934,34.735-6.956,47.631,5.939
							l71.072,71.073C351.331,284.523,351.331,309.312,336.095,324.547z"/>
					</svg>
				</a>
				<a class="external" target="_blank" href="%s" id="t=%s">
				%s
		`, text, timecodeUrl, text, text))
	} else {
		_, _ = w.WriteString("</a></span>")
	}
	return ast.WalkContinue, nil
}

type timecodes struct {
}

// Timecodes is an extension for Goldmark that replaces [00:00:00] with links
var Timecodes = &timecodes{}

func New() goldmark.Extender {
	return &timecodes{}
}

func (e *timecodes) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(NewParser(), 100)),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(util.Prioritized(NewTimecodeHTMLRenderer(), 100)),
	)
}
