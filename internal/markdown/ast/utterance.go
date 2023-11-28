package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

type Utterance struct {
	BaseBlock
	FileNode

	Speaker      Speaker
	IsNewSpeaker bool
}

func (s *Utterance) Dump(source []byte, level int) {
	gast.DumpHelper(s, source, level, nil, nil)
}

var KindUtterance = gast.NewNodeKind("Utterance")

func (s *Utterance) Kind() gast.NodeKind {
	return KindUtterance
}

func (u *Utterance) IsSpeaker(id string) bool {
	return id == u.Speaker.GetID()
}

func (u *Utterance) IsSandwichedBetween(speakerId string) bool {
	prev, next := u.Prev(), u.Next()
	if prev == nil || next == nil {
		return false
	}
	return next.IsSpeaker(speakerId) && prev.IsSpeaker(speakerId)
}

func (u *Utterance) PrevAndNextIsSameSpeaker() bool {
	prev, next := u.Prev(), u.Next()
	if prev == nil || next == nil {
		return false
	}
	return next.IsSpeaker(prev.Speaker.GetID())
}

func (u *Utterance) IsSandwichingPrevious() bool {
	prev := u.Prev()
	if prev == nil {
		return false
	}
	return prev.IsSandwichedBetween(u.Speaker.GetID())
}

func (u *Utterance) Prev() *Utterance {
	prev, _ := u.PreviousSibling().(*Utterance)
	return prev
}

func (u *Utterance) Next() *Utterance {
	next, _ := u.NextSibling().(*Utterance)
	return next
}

func (u *Utterance) PrevIsPrimarySpeaker() bool {
	prev := u.Prev()
	if prev == nil {
		return false
	}
	return prev.Speaker.GetIsPrimarySpeaker()
}

func IsPrimarySpeaker(node gast.Node) bool {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		utterance, ok := parent.(*Utterance)
		if ok {
			return utterance.Speaker.GetIsPrimarySpeaker()
		}
	}

	return false
}
