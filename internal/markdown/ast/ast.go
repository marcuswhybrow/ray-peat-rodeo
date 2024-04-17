package ast

import (
	"github.com/mitchellh/mapstructure"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

type FrontMatter struct {
	Source struct {
		Series   string
		Title    string
		Url      string
		Kind     string
		Duration string
	}
	Speakers      map[string]string
	Transcription struct {
		Url    string
		Kind   string
		Date   string
		Author string
	}
	Added struct {
		Date   string
		Author string
	}
}

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

func (n *FileNode) FrontMatter() FrontMatter {
	var frontMatter FrontMatter
	mapstructure.Decode(n.OwnerDocument().Meta(), &frontMatter)
	return frontMatter
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
