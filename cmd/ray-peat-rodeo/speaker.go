package main

type Speaker struct {
	ID               string
	Name             string
	AvatarPath       string
	IsPrimarySpeaker bool
}

// Implements ast.Speaker interface

func (s *Speaker) GetID() string {
	return s.ID
}

func (s *Speaker) GetName() string {
	return s.Name
}

func (s *Speaker) GetAvatarPath() string {
	return s.AvatarPath
}

func (s *Speaker) GetIsPrimarySpeaker() bool {
	return s.IsPrimarySpeaker
}
