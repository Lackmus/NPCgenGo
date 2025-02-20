package service

import "github.com/lackmus/npcgengo/model"

func CreateNPCWithOptions(npcType string, faction string, randomizerService RandomizerService) model.NPC {
	return NewNPCBuilder().
		WithType(npcType).
		WithFaction(faction).
		BuildWithRandom(randomizerService)
}
