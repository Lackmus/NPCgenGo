package npc_components

import (
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

// NewNPCSubtypeComponent creates a new component for the NPCSubtype.
func (n NPCSubtype) NewNPCSubtypeComponent() *NPCComponent {
	return &NPCComponent{
		Name:  CompSubtype,
		Value: n.Name,
	}
}

// NewNPCSubtypeComponentWithStats creates a new component for the NPCSubtype with stats.
func (n NPCSubtype) NewNPCSubtypeStatsComponent() *NPCComponent {
	statsString := ""
	for _, v := range n.Stats {
		statsString += v + ": " + helper.RandomInt(1, 10) + ", "
	}
	return &NPCComponent{
		Name:  CompStats,
		Value: statsString[:len(statsString)-2],
	}
}

// NewNPCSubtypeComponentWithEquipment creates a new component for the NPCSubtype with equipment.
func (n NPCSubtype) NewNPCSubtypeEquipmentComponent() *NPCComponent {
	itemString := ""
	for k, v := range n.EquipmentOptions {
		itemString += k + ": " + helper.GetRandomElement(v) + ", "
	}
	return &NPCComponent{
		Name:  CompItems,
		Value: itemString[:len(itemString)-2],
	}
}

// NewSubtypeDescription returns the description component of the NPCSubtype.
func (n NPCSubtype) GetDescription() string {
	return helper.GetRandomElement(n.Description)
}
