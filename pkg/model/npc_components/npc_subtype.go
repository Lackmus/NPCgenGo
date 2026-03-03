package npc_components

import (
	"strings"

	helper "github.com/lackmus/npcgengo/internal/platform/helpers"
)

type NPCSubtype struct {
	Name             string
	NpcTypeName      string
	Stats            []string
	EquipmentOptions map[string][]string
}

func (n NPCSubtype) GetName() string {
	return n.Name
}

func (n NPCSubtype) GetStats() string {
	var sb strings.Builder
	for _, v := range n.Stats {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(v + ": " + helper.RandomInt(1, 10))
	}
	return sb.String()
}

func (n NPCSubtype) GetEquipment() string {
	var sb strings.Builder
	for k, v := range n.EquipmentOptions {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(k + ": " + helper.GetRandomElement(v))
	}
	return sb.String()
}
