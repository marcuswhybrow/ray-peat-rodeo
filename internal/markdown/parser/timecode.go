package parser

import (
	"bytes"
	"log"
	"strconv"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/markdown/ast"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type TimecodeParser struct{}

func NewTimecodeParser() *TimecodeParser {
	return &TimecodeParser{}
}

func (w *TimecodeParser) Trigger() []byte {
	return []byte{'['}
}

func (w *TimecodeParser) Parse(parent gmAst.Node, block text.Reader, pc parser.Context) gmAst.Node {
	line, _ := block.PeekLine()

	i := bytes.Index(line, []byte{']'})
	if i < 2 {
		return nil
	}

	consumed := i + 1

	sections := bytes.Split(line[1:i], []byte{':'})
	timecode := ast.NewTimecode()
	timecode.Source = string(line[:i])

	n := len(sections)
	if n > 3 || n < 2 {
		return nil
	}

	seconds, err := strconv.Atoi(string(sections[n-1]))
	if err != nil {
		return nil
	}
	timecode.Seconds = seconds

	minutes, err := strconv.Atoi(string(sections[n-2]))
	if err != nil {
		return nil
	}
	timecode.Minutes = minutes

	var hours int
	if n < 3 {
		hours = 0
	} else {
		hours, err = strconv.Atoi(string(sections[n-3]))
		if err != nil {
			return nil
		}
	}
	timecode.Hours = hours
	timecode.IsRaySpeaking = isRay(pc)

	block.Advance(consumed)

	return timecode
}

func isRay(pc parser.Context) bool {
	for _, block := range pc.OpenedBlocks() {
		log.Printf("Inside %v", block.Node.Type())
		speaker, ok := block.Node.(*ast.Speaker)
		if ok && speaker.IsRay() {
			return true
		}
	}
	return false
}
