package goomi

import "strings"

func (s MIQualifiers) FindByName(name string) *MIQualifier {
	name = strings.ToLower(name)
	for _, qualifier := range s {
		if strings.ToLower(qualifier.Name) == name {
			return qualifier
		}
	}
	return nil
}

func (s MIQualifiers) HasIn() bool {
	in := s.FindByName("In")
	if in != nil {
		return in.Value.(bool)
	}
	return false
}

func (s MIQualifiers) HasOut() bool {
	in := s.FindByName("Out")
	if in != nil {
		return in.Value.(bool)
	}
	return false
}
