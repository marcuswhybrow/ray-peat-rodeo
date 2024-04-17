package parser

import (
	"bytes"
	"log"
	"slices"
	"strings"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var prevSpeakersKey = parser.NewContextKey()

type UtteranceParser struct{}

func NewSpeakerParser() parser.BlockParser {
	return &UtteranceParser{}
}

func (s *UtteranceParser) Trigger() []byte {
	return []byte{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}
}

func getSpeakerID(line []byte) ([]byte, int) {
	colonPos := bytes.Index(line, []byte{':'})
	if colonPos < 0 {
		return nil, 0
	}

	for _, byte := range line[:colonPos] {
		if !util.IsAlphaNumeric(byte) {
			return nil, 0
		}
	}

	return line[:colonPos], colonPos + 1
}

func (s *UtteranceParser) Open(parent gmAst.Node, reader text.Reader, pc parser.Context) (gmAst.Node, parser.State) {
	line, _ := reader.PeekLine()
	lineNumber, _ := reader.Position()
	lineNumber += 1

	speakerID, bytesConsumed := getSpeakerID(line)
	speakerIDStr := strings.Trim(string(speakerID), " ")

	if len(speakerIDStr) <= 0 {
		return nil, parser.Continue
	}

	file := ast.GetFile(pc)
	speakers := file.GetSpeakers()

	i := slices.IndexFunc(speakers, func(s ast.Speaker) bool {
		return s.GetID() == speakerIDStr
	})
	if i < 0 {
		log.Panicf("Failed to find speaker ID '%s' (line %v) in frontmatter of %v:\n\n%s\n", speakerIDStr, lineNumber, file.GetPath(), line)
	}
	speaker := speakers[i]

	prevSpeakers := pc.ComputeIfAbsent(prevSpeakersKey, func() interface{} {
		return &[]*ast.Utterance{}
	}).(*[]*ast.Utterance)

	isNewSpeaker := true
	for _, prev := range *prevSpeakers {
		if prev.Speaker.GetID() == speakerIDStr {
			isNewSpeaker = false
			break
		}
	}

	utterance := &ast.Utterance{
		Speaker:      speaker,
		IsNewSpeaker: isNewSpeaker,
	}

	*prevSpeakers = append(*prevSpeakers, utterance)
	pc.Set(prevSpeakersKey, prevSpeakers)

	reader.Advance(bytesConsumed)
	return utterance, parser.HasChildren
}

func (s *UtteranceParser) Continue(node gmAst.Node, reader text.Reader, pc parser.Context) parser.State {
	speaker := node.(*ast.Utterance)
	line, _ := reader.PeekLine()
	shortName, bytesConsumed := getSpeakerID(line)

	if shortName != nil && string(shortName) != speaker.Speaker.GetID() {
		return parser.Close
	}

	reader.Advance(bytesConsumed)
	return parser.HasChildren | parser.Continue
}

func (s *UtteranceParser) Close(node gmAst.Node, reader text.Reader, pc parser.Context) {
}

func (s *UtteranceParser) CanInterruptParagraph() bool {
	return false
}

func (s *UtteranceParser) CanAcceptIndentedLine() bool {
	return false
}
