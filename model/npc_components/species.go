package npc_components

import (
	"github.com/lackmus/npcgengo/helper"
)

// Species : A species in the game has a name and a name source.
type Species struct {
	Name        string
	NameSource  string
	Description []string
}

// Name : Return the name of the Species
func (s Species) GetName() string {
	return s.Name
}

// GetDescription : Return the description of the Species
func (s Species) GetDescription() string {
	return helper.GetRandomElement(s.Description)
}
