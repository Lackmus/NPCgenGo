package npc_components

import (
	h "github.com/lackmus/npcgengo/internal/platform/helpers"
)

type NameData struct {
	Name      string
	Forenames []string
	Surnames  []string
}

func (n NameData) GetName() string {
	return n.Name
}

func (n NameData) GenerateName() string {
	forname := h.GetRandomElement(n.Forenames)
	surname := h.GetRandomElement(n.Surnames)
	return forname + " " + surname
}
