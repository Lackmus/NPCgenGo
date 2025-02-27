// Description: This file contains the Species struct and its methods.
package npc_components

import (
	"github.com/lackmus/npcgengo/helper"
)

// Species : Represents a species of NPC
// It represents a species of NPC
type Species struct {
	Name        string
	NameSource  string
	Description []string
}

// GetName : Returns the name of the Species
// Returns the name of the Species
func (s Species) GetName() string {
	return s.Name
}

// GetDescription : Return a random description of the Species
// Returns a random description of the Species
func (s Species) GetDescription() string {
	return helper.GetRandomElement(s.Description)
}
