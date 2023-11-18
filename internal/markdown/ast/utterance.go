package ast

import (
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

type Utterance struct {
	BaseBlock
	FileNode

	// Initials definied at the start of a paragraph used to lookup a full name.
	//
	// # Example
	//   ---
	//   speakers:
	//     RP: Ray Peat
	//   ---
	//   RP: Hi, Ray Here.
	//
	// RP is the SpeakerID, which refers to speakers.RP in the front matter.
	SpeakerID string

	// Is first definition of this speaker ID.
	IsNewSpeaker bool
}

func NewSpeaker() *Utterance {
	return &Utterance{}
}

func (s *Utterance) Dump(source []byte, level int) {
	gast.DumpHelper(s, source, level, nil, nil)
}

var KindSpeaker = gast.NewNodeKind("Speaker")

func (s *Utterance) Kind() gast.NodeKind {
	return KindSpeaker
}

func (u *Utterance) SpeakerName() string {
	return u.FrontMatter().Speakers[u.SpeakerID]
}

func (s *Utterance) IsRay() bool {
	return strings.Trim(s.SpeakerID, " ") == "RP"
}

func (u *Utterance) IsSpeaker(id string) bool {
	return id == u.SpeakerID
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
	return next.IsSpeaker(prev.SpeakerID)
}

func (u *Utterance) IsSendwichingPrevious() bool {
	prev := u.Prev()
	if prev == nil {
		return false
	}
	return prev.IsSandwichedBetween(u.SpeakerID)
}

func (u *Utterance) Prev() *Utterance {
	prev, _ := u.PreviousSibling().(*Utterance)
	return prev
}

func (u *Utterance) Next() *Utterance {
	next, _ := u.NextSibling().(*Utterance)
	return next
}

func (u *Utterance) PrevIsRay() bool {
	prev := u.Prev()
	if prev == nil {
		return false
	}
	return prev.IsRay()
}

func IsRaySpeaking(node gast.Node) bool {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		speaker, ok := parent.(*Utterance)
		if ok {
			return speaker.IsRay()
		}
	}

	return false
}
