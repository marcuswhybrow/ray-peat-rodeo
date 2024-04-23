package parser

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/cache"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/global"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type githubIssueParser struct{}

func NewGitHubIssueParser() gparser.InlineParser {
	return &githubIssueParser{}
}

func (p *githubIssueParser) Trigger() []byte {
	return []byte{'{'}
}

func (p *githubIssueParser) Parse(parent gast.Node, block text.Reader, pc gparser.Context) gast.Node {
	line, _ := block.PeekLine()

	if len(line) < 4 {
		return nil
	}

	inside, _, foundEnd := strings.Cut(string(line[1:]), "}")
	if !foundEnd {
		return nil
	}

	trimmed := strings.Trim(inside, " ")

	if trimmed[0] != '#' {
		return nil
	}

	id, err := strconv.Atoi(trimmed[1:])
	if err != nil {
		return nil
	}

	httpCache := pc.Get(ast.HTTPCacheKey).(*cache.HTTPCache)
	url := global.GitHubIssueLink(id)
	key := "title"
	handler := func(res *http.Response) string {
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse HTTP response for url: %v", res.Request.URL.String()))
		}
		return doc.Find("h1 bdi").Text()
	}
	title := <-httpCache.Get(url, key, handler)

	asset := ast.GetAsset(pc)
	asset.RegisterIssue(id)

	block.Advance(len(inside) + 2)
	return ast.NewGitHubIssue(id, title)
}
