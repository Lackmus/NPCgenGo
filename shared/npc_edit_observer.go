package shared

import "github.com/lackmus/npcgengo/model"

// observer interface
type NPCEditObserver interface {
	UpdateNPC(npc model.NPC)
	UpdateField(field string, value any)
}
