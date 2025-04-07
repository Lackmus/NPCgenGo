// Description: This file contains the RandomizerService struct and its methods. The RandomizerService struct provides randomization services for NPC creation.
package service

import (
	"github.com/lackmus/npcgengo/helper"
)

// RandomizerService provides randomization services for NPC creation.
// It provides randomization services for NPC creation.
type RandomizerService struct {
	creationData       *CreationDataService
	npcCreationOptions *NPCCreationOptions
}

// NewRandomizerService creates a new RandomizerService.
// It returns a new RandomizerService.
func NewRandomizerService(creationData *CreationDataService, npcCreationOptions *NPCCreationOptions) *RandomizerService {
	return &RandomizerService{
		creationData:       creationData,
		npcCreationOptions: npcCreationOptions,
	}
}

// RandomTrait returns a random trait from the available traits.
func (r *RandomizerService) RandomTrait() string {
	trait := helper.GetRandomElement(r.npcCreationOptions.Traits)
	return trait
}

// RandomFaction returns a random faction from the available factions.
func (r *RandomizerService) RandomFaction() string {
	faction := helper.GetRandomElement(r.npcCreationOptions.Factions)
	return faction
}

// RandomSpecies returns a random species from the available species.
func (r *RandomizerService) RandomSpecies() string {
	species := helper.GetRandomElement(r.npcCreationOptions.Species)
	return species
}

// RandomType returns a random NPC type from the available types.
func (r *RandomizerService) RandomType() string {
	npcType := helper.GetRandomElement(r.npcCreationOptions.NpcTypes)
	return npcType
}

// RandomSubtype returns a random subtype for the given NPC type.
func (r *RandomizerService) RandomSubtype(npcType string) string {
	subtype := helper.GetRandomElement(r.npcCreationOptions.NpcSubtypeForTypeMap[npcType])
	return subtype
}
