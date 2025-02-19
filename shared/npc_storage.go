package shared

import "github.com/lackmus/npcgengo/model"

type NPCStorage interface {
	LoadNPC(id string) (model.NPC, error)
	LoadAllNPC() (map[string]model.NPC, error)
	SaveAllNPC(npcs map[string]model.NPC) error
	SaveNPC(npc model.NPC) error
	DeleteNPC(id string) error
	DeleteAllNPC() error
}
