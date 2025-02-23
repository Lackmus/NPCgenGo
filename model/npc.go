// model/npc.go
package model

import (
	"fmt"

	"github.com/lackmus/npcgengo/helper"
)

// NPC represents an immutable non-player character.
type NPC struct {
	id, name, faction, species, npcType, npcSubtype, trait, drive, description string
	stats                                                                      map[string]int
	items, abilities                                                           map[string]string
}

// Getters to access the fields
func (n NPC) ID() string                   { return n.id }
func (n NPC) Name() string                 { return n.name }
func (n NPC) Faction() string              { return n.faction }
func (n NPC) Species() string              { return n.species }
func (n NPC) NPCType() string              { return n.npcType }
func (n NPC) NPCSubtype() string           { return n.npcSubtype }
func (n NPC) Trait() string                { return n.trait }
func (n NPC) Drive() string                { return n.drive }
func (n NPC) Stats() map[string]int        { return helper.CopyMap(n.stats) }
func (n NPC) Items() map[string]string     { return helper.CopyMap(n.items) }
func (n NPC) Abilities() map[string]string { return helper.CopyMap(n.abilities) }
func (n NPC) Description() string          { return n.description }

// NewNPC is the constructor to create an immutable NPC.
func NewNPC(
	id, name, faction, species, npcType, npcSubtype, trait, drive, description string,
	stats map[string]int,
	items map[string]string,
	abilities map[string]string,
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
		drive:       drive,
		description: description,
		stats:       helper.CopyMap(stats),
		items:       helper.CopyMap(items),
		abilities:   helper.CopyMap(abilities),
	}
}

// string returns a string representation of the NPC. with all the fields. \n are used to separate the fields.
func (n NPC) String() string {
	return "Name: " + n.name + "\n" +
		"Faction: " + n.faction + "\n" +
		"Species: " + n.species + "\n" +
		"Type: " + n.npcType + "\n" +
		"Subtype: " + n.npcSubtype + "\n" +
		"Trait: [" + n.trait + "]\n" +
		"Stats: " + n.PrintStats() + "\n" +
		"Items: " + n.PrintItems() + "\n"
}

func (n NPC) PrintStats() string {
	result := ""
	if len(n.stats) == 0 {
		return result
	}
	for k, v := range n.stats {
		result += fmt.Sprintf("[%s: %d] ", k, v)
	}
	return result
}

func (n NPC) PrintItems() string {
	result := ""
	if len(n.items) == 0 {
		return result
	}
	for k, v := range n.items {
		result += fmt.Sprintf("[%s: %s] ", k, v)
	}
	return result
}
