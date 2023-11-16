package parser

import (
	"bytes"
	"slices"
	"strings"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var prevSpeakersKey = parser.NewContextKey()
var isRaySpeakingKey = parser.NewContextKey()

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

func getShortName(line []byte) ([]byte, int) {
	colon := bytes.Index(line, []byte{':'})
	if colon < 0 {
		return nil, 0
	}

	for _, byte := range line[:colon] {
		if !util.IsAlphaNumeric(byte) {
			return nil, 0
		}
	}

	return line[:colon], colon + 1
}

func (s *UtteranceParser) Open(parent gmAst.Node, reader text.Reader, pc parser.Context) (gmAst.Node, parser.State) {
	line, _ := reader.PeekLine()
	shortName, bytesConsumed := getShortName(line)

	speaker := ast.NewSpeaker()
	speaker.SpeakerID = strings.Trim(string(shortName), " ")

	pc.Set(isRaySpeakingKey, speaker.IsRay())

	if len(speaker.SpeakerID) <= 0 {
		return nil, parser.HasChildren | parser.Continue
	}

	var prevSpeakers []*ast.Utterance
	prevSpeakersVal := pc.Get(prevSpeakersKey)
	if prevSpeakersVal == nil {
		prevSpeakers = []*ast.Utterance{}
	} else {
		prevSpeakers = *prevSpeakersVal.(*[]*ast.Utterance)
	}

	speaker.IsNewSpeaker = !slices.ContainsFunc(prevSpeakers, func(s *ast.Utterance) bool {
		return s.SpeakerID == speaker.SpeakerID
	})

	prevSpeakers = append(prevSpeakers, speaker)
	pc.Set(prevSpeakersKey, &prevSpeakers)

	reader.Advance(bytesConsumed)
	return speaker, parser.HasChildren
}

func (s *UtteranceParser) Continue(node gmAst.Node, reader text.Reader, pc parser.Context) parser.State {
	speaker := node.(*ast.Utterance)
	line, _ := reader.PeekLine()
	shortName, bytesConsumed := getShortName(line)

	if shortName != nil && string(shortName) != speaker.SpeakerID {
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
