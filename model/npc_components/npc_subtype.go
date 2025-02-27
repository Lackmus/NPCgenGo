// Description: This file contains the NPCSubtype struct and its methods.
package npc_components

import (
	"strings"

	"github.com/lackmus/npcgengo/helper"
)

// NPCSubtype : Represents a subtype of NPC
// It represents a subtype of NPC
type NPCSubtype struct {
	Name             string
	NpcTypeName      string
	Description      []string
	Stats            []string
	EquipmentOptions map[string][]string
}

// GetName : Returns the name of the NPCSubtype
// Returns the name of the NPCSubtype
func (n NPCSubtype) GetName() string {
	return n.Name
}

// GetStats : Returns the stats of the NPCSubtype
// Returns the stats of the NPCSubtype
func (n NPCSubtype) GetStats() string {
	var sb strings.Builder
	for _, v := range n.Stats {
		sb.WriteString("\n  - " + v + ": " + helper.RandomInt(1, 10))
	}
	return sb.String()
}

// GetEquipment : Returns the equipment of the NPCSubtype
// Returns the equipment of the NPCSubtype
func (n NPCSubtype) GetEquipment() string {
	var sb strings.Builder
	for k, v := range n.EquipmentOptions {
		sb.WriteString("\n  - " + k + ": " + helper.GetRandomElement(v))
	}
	return sb.String()
}

// GetDescription : Returns the description of the NPCSubtype
// Returns the description of the NPCSubtype
func (n NPCSubtype) GetDescription() string {
	return helper.GetRandomElement(n.Description)
}
