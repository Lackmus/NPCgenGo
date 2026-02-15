package model

import "slices"

type NPCGroup struct {
	Name        string
	LocationID  string
	NPCIDs      []string
	Description string
}

func NewNPCGroup(name string) *NPCGroup {
	return &NPCGroup{
		Name:        name,
		NPCIDs:      []string{},
		Description: "",
	}
}

func (g *NPCGroup) AddNPC(npcID string) {
	for _, id := range g.NPCIDs {
		if id == npcID {
			return
		}
	}
	g.NPCIDs = append(g.NPCIDs, npcID)
}

func (g *NPCGroup) RemoveNPC(npcID string) {
	for i, id := range g.NPCIDs {
		if id == npcID {
			g.NPCIDs = slices.Delete(g.NPCIDs, i, i+1)
			return
		}
	}
}

func (g *NPCGroup) ContainsNPC(npcID string) bool {
	return slices.Contains(g.NPCIDs, npcID)
}
