package blog

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var markdownParser = goldmark.New(
	goldmark.WithRendererOptions(html.WithUnsafe()),
	goldmark.WithExtensions(
		meta.New(meta.WithStoresInDocument()),
		extension.Typographer,
	),
)

type BlogPost struct {
	Path             string
	OutPath          string
	ID               string
	Permalink        string
	Date             string
	Title            string
	Author           string
	AuthorAvatarPath string
	Markdown         []byte
	HTML             []byte
}

type frontMatter struct {
	Author string
	Title  string
}

func NewBlogPost(filePath string, avatarPaths *catalog.AvatarPaths) *BlogPost {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("Failed to open blog post '%v': %v", filePath, err)
	}

	var html bytes.Buffer
	document := markdownParser.Parser().Parse(text.NewReader(fileBytes))
	meta := document.OwnerDocument().Meta()

	markdownParser.Renderer().Render(&html, fileBytes, document)

	frontMatter := frontMatter{}
	err = mapstructure.Decode(meta, &frontMatter)
	if err != nil {
		log.Panicf("Failed to decode front matter in blog post '%v': %v", filePath, err)
	}

	base := filepath.Base(filePath)
	name, _ := strings.CutSuffix(base, filepath.Ext(filePath))
	date := name[:10]
	id := name[11:]

	cleanID := strings.ToLower(id)
	cleanID = url.QueryEscape(cleanID)
	if cleanID != id {
		log.Panicf("Bag blog post file name '%v': filename contains uppercase characters or characters that must be escaped to be URL safe into '%v'", filePath, cleanID)
	}

	authorAvatarPath := avatarPaths.Get(frontMatter.Author)

	return &BlogPost{
		Path:             filePath,
		OutPath:          path.Join("blog", id, "index.html"),
		ID:               id,
		Permalink:        "/blog/" + id,
		Date:             date,
		Author:           frontMatter.Author,
		AuthorAvatarPath: authorAvatarPath,
		Title:            frontMatter.Title,
		Markdown:         fileBytes,
		HTML:             html.Bytes(),
	}
}

func (b *BlogPost) Write() error {
	buildFile, _ := utils.MakeFile(b.OutPath)

	err := RenderBlogPost(b).Render(context.Background(), buildFile)
	if err != nil {
		return fmt.Errorf("Failed to render template: %v", err)
	}

	return nil

}
