package shared

import "github.com/lackmus/npcgengo/model"

type NPCEditController interface {
	CreateNPC(npcType string, faction string)
	LoadNPC(npc model.NPC)
	EditNPC(npc model.NPC)
	SaveNPC() model.NPC
	RandomizeField(field string)
	SaveField(field string, value any)
	GetFieldOptions(field string) []string
	RegisterObserver(o NPCEditObserver)
	NotifyObserversField(field string, value any)
}
