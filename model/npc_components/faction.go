// Description: This file contains the Faction struct which is used to represent a faction in the game.
package npc_components

// Faction represents a faction in the game.
// It represents a faction in the game.
type Faction struct {
	Name        string
	SpeciesList []string
	Description []string
}

// GetName : Returns the name of the faction.
// Returns the name of the faction.
func (f Faction) GetName() string {
	return f.Name
}
