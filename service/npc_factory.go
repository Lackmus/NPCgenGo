package service

import (
	"fmt"

	"github.com/lackmus/npcgengo/model"
)

func CreateNPCWithOptions(npcType string, faction string, randomizerService *RandomizerService) model.NPC {
	npcSubType := randomizerService.RandomSubtype(npcType)
	fmt.Println("Subtype: " + npcSubType)
	species := randomizerService.RandomSpecies()
	return NewNPCBuilder().
		WithType(npcType).
		WithFaction(faction).
		WithID(randomizerService.GenerateID()).
		WithSpecies(randomizerService.RandomSpecies()).
		WithName(randomizerService.GenerateName(species)).
		WithTrait(randomizerService.RandomTrait()).
		WithComponent("Stats", fmt.Sprintf("%v", randomizerService.ApplySubtypeStats(npcSubType))).
		WithComponent("Items", fmt.Sprintf("%v", randomizerService.GenerateEquipment(npcSubType))).
		Build()
}
