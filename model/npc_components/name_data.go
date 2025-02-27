// Description: This file contains the NameData struct and the methods associated with it.
package npc_components

import (
	h "github.com/lackmus/npcgengo/helper"
)

// NameData : Represents a name in the game
// It represents a name in the game
type NameData struct {
	Name      string
	Forenames []string
	Surnames  []string
}

// GetName : Returns the name of the name data
// Returns the name of the name data
func (n NameData) GetName() string {
	return n.Name
}

// GenerateName : Generates a random name
// Generates a random name
func (n NameData) GenerateName() string {
	forname := h.GetRandomElement(n.Forenames)
	surname := h.GetRandomElement(n.Surnames)
	return forname + " " + surname
}
