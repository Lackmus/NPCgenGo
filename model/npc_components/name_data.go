package npc_components

import (
	h "github.com/lackmus/npcgengo/helper"
)

// NameData : A name has a name, a list of forenames and a list of surnames.
type NameData struct {
	Name      string
	Forenames []string
	Surnames  []string
}

// Name : Return the name of the NameData
func (n NameData) GetName() string {
	return n.Name
}

// NewNameComponent : Create a new component for the name
func (n NameData) GenerateName() string {
	forname := h.GetRandomElement(n.Forenames)
	surname := h.GetRandomElement(n.Surnames)
	return forname + " " + surname
}
