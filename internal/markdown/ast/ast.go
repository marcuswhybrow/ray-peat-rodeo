package ast

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

var FileKey = parser.NewContextKey()
var SourceKey = parser.NewContextKey()
var HTTPCacheKey = parser.NewContextKey()
var AvatarsKey = parser.NewContextKey()

type File interface {
	GetMarkdown() []byte
	GetPath() string
	RegisterMention(mention *Mention)
	RegisterIssue(id int)
	GetID() string
	GetPermalink() string
	GetSpeakers() []Speaker
	GetSourceURL() string
}

type Speaker interface {
	GetID() string
	GetName() string
	GetAvatarPath() string
	GetIsPrimarySpeaker() bool
}

type Speakers []Speaker

type FileNode struct {
	gast.BaseNode
}

func GetFile(pc parser.Context) File {
	file, ok := pc.Get(FileKey).(File)
	if !ok {
		panic("Failed to coerce FileKey in parser context to File interface")
	}
	return file
}

type BaseBlock struct {
	gast.BaseBlock
}

type BaseInline struct {
	gast.BaseInline
}
