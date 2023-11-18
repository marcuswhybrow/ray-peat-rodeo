package ast

type Mentionable struct {
	Primary   MentionablePart
	Secondary MentionablePart
}

func NewMentionable(primary, secondary MentionablePart) *Mentionable {
	return &Mentionable{
		Primary:   primary,
		Secondary: secondary,
	}
}

func (m *Mentionable) ID() string {
	id := m.Primary.ID()
	if !m.Secondary.IsEmpty() {
		id = m.Secondary.ID() + "@" + id
	}
	return id
}

func (m *Mentionable) HasSecondary() bool {
	return len(m.Secondary.Cardinal) > 0
}

func (m *Mentionable) PermalinkForPrimary() string {
	return "/" + m.Primary.ID()
}

func (m *Mentionable) Permalink() string {
	return m.PermalinkForPrimary() + "#" + m.Secondary.ID()
}

func (m *Mentionable) PopupPermalink() string {
	return "/api/mentionable/popup/" + m.ID()
}

func (m *Mentionable) Ultimate() *MentionablePart {
	if len(m.Secondary.Cardinal) > 0 || len(m.Secondary.Prefix) > 0 {
		return &m.Secondary
	} else {
		return &m.Primary
	}
}

func (m *Mentionable) IsDuplicate(other Mentionable) bool {
	if m.ID() == other.ID() {
		if m.Primary.Cardinal != other.Primary.Cardinal {
			return true
		}
		return m.Secondary.Cardinal != other.Secondary.Cardinal
	}
	return false
}

func (m *Mentionable) IsMoreComplex(other Mentionable) bool {
	if !m.Primary.HasPrefix() && other.Primary.HasPrefix() {
		return false
	}
	if !m.Secondary.HasPrefix() && other.Secondary.HasPrefix() {
		return false
	}
	return true
}

func (m *Mentionable) AsSignature() string {
	sig := m.Primary.CardinalFirst()
	if len(m.Secondary.CardinalFirst()) > 0 {
		sig += " > " + m.Secondary.CardinalFirst()
	}
	return sig
}

var EmptyMentionablePart = MentionablePart{"", "", ""}
