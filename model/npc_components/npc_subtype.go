package npc_components

import (
	"github.com/lackmus/npcgengo/helper"
)

// NpcSubtype is a type of npc.
type NPCSubtype struct {
	Name             string
	NpcTypeName      string
	Description      string
	Stats            []string
	EquipmentOptions map[string][]string
}

// GetName returns the name of the NPCSubtype.
func (n NPCSubtype) GetName() string {
	return n.Name
}

// NewNPCSubtypeComponent creates a new component for the NPCSubtype.
func (n NPCSubtype) NewNPCSubtypeComponent() *Component {
	return &Component{
		Name:  CompSubtype,
		Value: n.Name,
	}
}

// NewNPCSubtypeComponentWithStats creates a new component for the NPCSubtype with stats.
func (n NPCSubtype) NewNPCSubtypeComponentWithStats() *Component {
	// make stats string with random int value 1-10
	statsString := ""
	for _, v := range n.Stats {
		statsString += v + ": " + helper.RandomInt(1, 10) + ", "
	}
	return &Component{
		Name:  CompStats,
		Value: statsString[:len(statsString)-2],
	}
}

// NewNPCSubtypeComponentWithEquipment creates a new component for the NPCSubtype with equipment.
func (n NPCSubtype) NewNPCSubtypeComponentWithEquipment() *Component {
	itemString := ""
	for k, v := range n.EquipmentOptions {
		itemString += k + ": " + helper.GetRandomElement(v) + ", "
	}
	return &Component{
		Name: CompItems,
		// helper.randomElement

		Value: itemString[:len(itemString)-2],
	}
}
