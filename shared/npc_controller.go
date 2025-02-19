package shared

import "github.com/lackmus/npcgengo/model"

type NPCController interface {
	GetAllNpcs() []model.NPC
	RegisterObserver(o NPCObserver)
	NotifyObservers()
}
