package shared

import (
	"github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components"
)

// observer interface
type NPCEditObserver interface {
	UpdateNPC(npc model.NPC)
	UpdateField(field cp.CompEnum, value any)
	OnNPCEditError(err error) // New method for error reporting
}
