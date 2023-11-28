package main

import (
	"io/fs"
	"os"
	"path"
	"strings"
)

type AvatarPaths struct {
	Paths map[string]string
}

func NewAvatarPaths() *AvatarPaths {
	paths := map[string]string{}
	fs.WalkDir(os.DirFS("./internal"), "assets/images/avatars", func(filePath string, entry fs.DirEntry, err error) error {
		fileStem := path.Base(filePath)
		ext := path.Ext(fileStem)
		fileName, _ := strings.CutSuffix(fileStem, ext)

		if !entry.IsDir() {
			paths[fileName] = filePath
		}

		return nil
	})
	return &AvatarPaths{
		Paths: paths,
	}
}

func (a *AvatarPaths) Get(speakerName string) string {
	lowerName := strings.ToLower(speakerName)
	kebabName := strings.ReplaceAll(lowerName, " ", "-")

	aPath := a.Paths[kebabName]
	if len(aPath) == 0 {
		return ""
	}

	return "/" + aPath
}
