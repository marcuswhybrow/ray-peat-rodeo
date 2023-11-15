package markdown

import "github.com/yuin/goldmark/parser"

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
}

var PermalinkKey = parser.NewContextKey()
var IDKey = parser.NewContextKey()
var SourceKey = parser.NewContextKey()
var HTTPCache = parser.NewContextKey()
