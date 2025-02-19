package shared

import "github.com/lackmus/npcgengo/model"

// observer interface
type NPCObserver interface {
	Update(npcs []model.NPC)
}
