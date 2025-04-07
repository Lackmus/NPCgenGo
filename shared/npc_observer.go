package shared

import "github.com/lackmus/npcgengo/model"

// observer interface
type NPCObserver interface {

	// UpdateNPC updates the NPC with the given model.NPC.
	// It takes a model.NPC as a parameter and returns nothing.
	Update(npcs []model.NPC)
}
