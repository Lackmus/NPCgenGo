// model/npc.go
package model

import (
	"fmt"
)

// NPC represents an immutable non-player character.
type NPC struct {
	id, name, faction, species, npcType, npcSubtype, trait, description string
	components                                                          map[string]string
}

// Getters to access the fields
func (n NPC) ID() string         { return n.id }
func (n NPC) Name() string       { return n.name }
func (n NPC) Faction() string    { return n.faction }
func (n NPC) Species() string    { return n.species }
func (n NPC) NPCType() string    { return n.npcType }
func (n NPC) NPCSubtype() string { return n.npcSubtype }
func (n NPC) Trait() string      { return n.trait }
func (n NPC) Description() string {
	return n.description
}
func (n NPC) Components() map[string]string {
	return n.components
}

// NewNPC is the constructor to create an immutable NPC.
func NewNPC(
	id, name, faction, species, npcType, npcSubtype, trait, description string,
	components map[string]string,
) NPC {
	// Optionally perform deep copies of the maps here.
	return NPC{
		id:          id,
		name:        name,
		faction:     faction,
		species:     species,
		npcType:     npcType,
		npcSubtype:  npcSubtype,
		trait:       trait,
		description: description,
		components:  components,
	}
}

// string returns a string representation of the NPC. with all the fields. \n are used to separate the fields.
func (n NPC) String() string {
	return "Name: " + n.name + "\n" +
		"Faction: " + n.faction + "\n" +
		"Species: " + n.species + "\n" +
		"Type: " + n.npcType + "\n" +
		"Subtype: " + n.npcSubtype + "\n" +
		"Trait: " + n.trait + "\n" +
		n.PrintComponents() +
		"Description: " + n.description
}

// print components of the NPC
func (n NPC) PrintComponents() string {
	var result string
	for k, v := range n.components {
		result += fmt.Sprintf("%s: %s\n", k, v)
	}
	return result
}
