package parser

import (
	"bytes"
	"slices"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var prevSpeakersKey = parser.NewContextKey()
var isRaySpeakingKey = parser.NewContextKey()

type SpeakerParser struct {
}

func NewSpeakerParser() parser.BlockParser {
	return &SpeakerParser{}
}

func (s *SpeakerParser) Trigger() []byte {
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

func (s *SpeakerParser) Open(parent gmAst.Node, reader text.Reader, pc parser.Context) (gmAst.Node, parser.State) {
	line, _ := reader.PeekLine()
	shortName, bytesConsumed := getShortName(line)

	speaker := ast.NewSpeaker()
	speaker.ShortName = string(shortName)

	pc.Set(isRaySpeakingKey, speaker.IsRay())

	if len(speaker.ShortName) <= 0 {
		return nil, parser.HasChildren | parser.Continue
	}

	var prevSpeakers []*ast.Speaker
	prevSpeakersVal := pc.Get(prevSpeakersKey)
	if prevSpeakersVal == nil {
		prevSpeakers = []*ast.Speaker{}
	} else {
		prevSpeakers = *prevSpeakersVal.(*[]*ast.Speaker)
	}

	var prevSpeaker *ast.Speaker
	if len(prevSpeakers) >= 1 {
		prevSpeaker = prevSpeakers[len(prevSpeakers)-1]
	}

	var penultimateSpeaker *ast.Speaker
	if len(prevSpeakers) >= 2 {
		penultimateSpeaker = prevSpeakers[len(prevSpeakers)-2]
	}

	speaker.IsHello = !slices.ContainsFunc(prevSpeakers, func(s *ast.Speaker) bool {
		return s.ShortName == speaker.ShortName
	})

	if speaker.IsHello {
		speaker.CanRetort = false
	} else if speaker.IsRay() && penultimateSpeaker.ShortName == speaker.ShortName {
		speaker.CanRetort = true
	} else if !speaker.IsRay() && prevSpeaker.IsRay() {
		speaker.CanRetort = true
	} else {
		speaker.CanRetort = false
	}

	prevSpeakers = append(prevSpeakers, speaker)
	pc.Set(prevSpeakersKey, &prevSpeakers)

	reader.Advance(bytesConsumed)
	return speaker, parser.HasChildren
}

func (s *SpeakerParser) Continue(node gmAst.Node, reader text.Reader, pc parser.Context) parser.State {
	speaker := node.(*ast.Speaker)
	line, _ := reader.PeekLine()
	shortName, bytesConsumed := getShortName(line)

	if shortName != nil && string(shortName) != speaker.ShortName {
		return parser.Close
	}

	reader.Advance(bytesConsumed)
	return parser.HasChildren | parser.Continue
}

func (s *SpeakerParser) Close(node gmAst.Node, reader text.Reader, pc parser.Context) {
}

func (s *SpeakerParser) CanInterruptParagraph() bool {
	return false
}

func (s *SpeakerParser) CanAcceptIndentedLine() bool {
	return false
}
