package shared

import "github.com/lackmus/npcgengo/pkg/product/model"

// observer interface
type NPCObserver interface {

	// Update updates observers with the current NPC list.
	Update(npcs []model.NPC)
}
