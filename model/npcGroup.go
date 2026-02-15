package model

import "slices"

// NPCGroup represents a collection of related NPCs
type NPCGroup struct {
	Name        string
	LocationID  string
	NPCIDs      []string
	Description string
}

// NewNPCGroup creates a new empty NPC group
func NewNPCGroup(name string) *NPCGroup {
	return &NPCGroup{
		Name:        name,
		NPCIDs:      []string{},
		Description: "",
	}
}

// AddNPC adds an NPC to the group if not already present
func (g *NPCGroup) AddNPC(npcID string) {
	// Check if NPC is already in the group
	for _, id := range g.NPCIDs {
		if id == npcID {
			return // NPC already in group
		}
	}
	g.NPCIDs = append(g.NPCIDs, npcID)
}

// RemoveNPC removes an NPC from the group
func (g *NPCGroup) RemoveNPC(npcID string) {
	for i, id := range g.NPCIDs {
		if id == npcID {
			// Remove the NPC ID from the slice
			g.NPCIDs = slices.Delete(g.NPCIDs, i, i+1)
			return
		}
	}
}

func (g *NPCGroup) ContainsNPC(npcID string) bool {
	return slices.Contains(g.NPCIDs, npcID)
}
