// Description: This file contains the NameData struct and the methods associated with it.
package npc_components

import (
	h "github.com/lackmus/npcgengo/helper"
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
