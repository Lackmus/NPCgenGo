package npc_components

import (
	"strings"

	"github.com/lackmus/npcgengo/helper"
)

// NpcSubtype is a type of npc.
type NPCSubtype struct {
	Name             string
	NpcTypeName      string
	Description      []string
	Stats            []string
	EquipmentOptions map[string][]string
}

// GetName returns the name of the NPCSubtype.
func (n NPCSubtype) GetName() string {
	return n.Name
}

// NewNPCSubtypeComponentWithStats creates a new component for the NPCSubtype with stats.
func (n NPCSubtype) GetStats() string {
	var sb strings.Builder
	for _, v := range n.Stats {
		sb.WriteString("\n  - " + v + ": " + helper.RandomInt(1, 10))
	}
	return sb.String()
}

// NewNPCSubtypeComponentWithEquipment creates a new component for the NPCSubtype with equipment.
func (n NPCSubtype) GetEquipment() string {
	var sb strings.Builder
	for k, v := range n.EquipmentOptions {
		sb.WriteString("\n  - " + k + ": " + helper.GetRandomElement(v))
	}
	return sb.String()
}

// NewSubtypeDescription returns the description component of the NPCSubtype.
func (n NPCSubtype) GetDescription() string {
	return helper.GetRandomElement(n.Description)
}
