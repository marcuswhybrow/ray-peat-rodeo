package blog

import (
	"context"
	"path/filepath"
	"strings"

	rprCatalog "github.com/marcuswhybrow/ray-peat-rodeo/internal/catalog"
	"github.com/marcuswhybrow/ray-peat-rodeo/internal/utils"
)

func Write(catalog *rprCatalog.Catalog) ([]*BlogPost, error) {
	postPaths := utils.Files(".", "assets/blog", func(filePath string) (*string, error) {
		ext := filepath.Ext(filePath)
		if strings.ToLower(ext) != ".md" {
			return nil, nil
		}
		return &filePath, nil
	})

	blogPosts := utils.Parallel(postPaths, func(filePath string) *BlogPost {
		blogPost := NewBlogPost(filePath, catalog.AvatarPaths)
		blogPost.Write()
		return blogPost
	})

	blogPage, _ := utils.MakePage("blog")
	component := BlogArchive(blogPosts)
	component.Render(context.Background(), blogPage)

	return blogPosts, nil
}
