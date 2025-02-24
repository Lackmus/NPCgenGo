package service

import (
	"github.com/lackmus/npcgengo/helper"
)

// RandomTrait returns a random trait from the provided options.
// struct

type RandomizerService struct {
	creationData       *CreationDataService
	npcCreationOptions *NPCCreationOptions
}

func NewRandomizerService(creationData *CreationDataService, npcCreationOptions *NPCCreationOptions) *RandomizerService {
	return &RandomizerService{
		creationData:       creationData,
		npcCreationOptions: npcCreationOptions,
	}
}

func (r *RandomizerService) RandomTrait() string {
	trait := helper.GetRandomElement(r.npcCreationOptions.Traits)
	return trait
}

func (r *RandomizerService) RandomFaction() string {
	faction := helper.GetRandomElement(r.npcCreationOptions.Factions)
	return faction
}

func (r *RandomizerService) RandomSpecies() string {
	species := helper.GetRandomElement(r.npcCreationOptions.Species)
	return species
}

func (r *RandomizerService) RandomType() string {
	npcType := helper.GetRandomElement(r.npcCreationOptions.NpcTypes)
	return npcType
}

func (r *RandomizerService) RandomSubtype(npcType string) string {
	subtype := helper.GetRandomElement(r.npcCreationOptions.NpcSubtypeForTypeMap[npcType])
	return subtype
}

func (r *RandomizerService) GenerateID() string {
	return helper.GenerateID()
}
