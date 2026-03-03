package shared

import (
	"github.com/lackmus/npcgengo/pkg/model"
	cp "github.com/lackmus/npcgengo/pkg/model/npc_components"
)

// observer interface
type NPCEditObserver interface {

	// UpdateNPC updates the NPC with the given model.NPC.
	// It takes a model.NPC as a parameter and returns nothing.
	UpdateNPC(npc model.NPC)

	// UpdateField updates the field of the NPC with the given value.
	// It takes a field of type cp.CompEnum and a value of any type as parameters and returns nothing.
	UpdateField(field cp.CompEnum, value any)

	// UpdateFieldWithName updates the field of the NPC with the given name and value.
	// It takes a field of type cp.CompEnum, a name of type string, and a value of any type as parameters and returns nothing.
	OnNPCEditError(err error) // New method for error reporting
}

