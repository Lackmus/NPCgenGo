package shared

import "github.com/lackmus/npcgengo/model"

type NPCEditController interface {
	CreateNPC(npcType string, faction string)
	EditNPC(npc model.NPC)
	SaveNPC() model.NPC
	RegisterObserver(o NPCEditObserver)
	NotifyObservers()
}
