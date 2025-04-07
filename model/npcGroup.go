package model

import (
	"slices"
	"strings"

	cp "github.com/lackmus/npcgengo/model/npc_components"
)

type NPCGroup struct {
	NPC
	NPCs []string
}

// NewNPCGroup creates a new NPCGroup with the given ID and location ID.
// It returns a pointer to the new NPCGroup.
func NewNPCGroup() *NPCGroup {
	return &NPCGroup{
		NPC: NPC{
			Components: make(map[cp.CompEnum]string),
		},
		NPCs: []string{},
	}
}

// AddNPC adds an NPC to the group.
func (g *NPCGroup) AddNPC(npcID string) {
	g.NPCs = append(g.NPCs, npcID)
}

// RemoveNPC removes an NPC from the group.
func (g *NPCGroup) RemoveNPC(npcID string) {
	for i, npc := range g.NPCs {
		if npc == npcID {
			g.NPCs = slices.Delete(g.NPCs, i, i+1)
			break
		}
	}
}

// String returns a string representation of the NPCGroup and its components.
// It returns a string representation of the NPCGroup and its components.
func (g *NPCGroup) String() string {
	var sb strings.Builder
	sb.WriteString(g.NPC.String())
	sb.WriteString("\n  NPCs:")
	for _, npc := range g.NPCs {
		sb.WriteString("\n    " + npc)
	}
	return sb.String()
}

// ShortString returns a short string representation of the NPCGroup and its components.
// It returns a short string representation of the NPCGroup and its components.
func (g *NPCGroup) ShortString() string {
	var sb strings.Builder
	sb.WriteString(g.NPC.ShortString())
	sb.WriteString("\n  NPCs:")
	for _, npc := range g.NPCs {
		sb.WriteString("\n    " + npc)
	}
	return sb.String()
}
