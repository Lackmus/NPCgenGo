package service

import (
	"github.com/lackmus/npcgengo/model"
)

// NPCFactory
type NPCFactory struct {
	randomizerService RandomizerService
}

func NewNPCFactory(randomizerService RandomizerService) *NPCFactory {
	return &NPCFactory{
		randomizerService: randomizerService,
	}
}

func (of *NPCFactory) CreateNPCWithOptions(npcType string, faction string) model.NPC {
	return NewNPC(of.randomizerService,
		WithID(""),
		WithType(npcType),
		WithSubType(""),
		WithStats(nil),
		WithItems(nil),
		WithFaction(faction),
		WithSpecies(""),
		WithName(""),
		WithTrait(""),
	)
}
