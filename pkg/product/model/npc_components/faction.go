// Description: This file contains the Faction struct which is used to represent a faction in the game.
package npc_components

type Faction struct {
	Name        string
	SpeciesList []string
}

func (f Faction) GetName() string {
	return f.Name
}
