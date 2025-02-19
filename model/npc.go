package model

import (
	"fmt"
)

// Npc is a struct that represents a row in the npc table.
type NPC struct {
	ID          string
	Name        string
	Faction     string
	Species     string
	NpcType     string
	NpcSubtype  string
	Trait       string
	Drive       string
	Stats       map[string]int
	Items       map[string]string
	Abilities   map[string]string
	Description string
}

// print the struct in a human readable format
func (n NPC) String() string {
	return fmt.Sprintf("Name: %s\nType: %s\nSpecies: %s\nSubtype: %s\nDescription: %s\nFaction: %s\nTrait: %s\nDrive: %s\nStats: %v\nItems: %v\nAbilities: %v",
		n.Name,
		n.NpcType,
		n.Species,
		n.NpcSubtype,
		n.Description,
		n.Faction,
		n.Trait,
		n.Drive,
		n.Stats,
		n.Items,
		n.Abilities)
}
