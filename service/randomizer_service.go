package service

import (
	"math/rand"

	"github.com/lackmus/npcgengo/helper"
	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/model/types"
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
	return helper.GetRandomElement(r.npcCreationOptions.Traits)
}

func (r *RandomizerService) RandomFaction() string {
	return helper.GetRandomElement(r.npcCreationOptions.Factions)
}

func (r *RandomizerService) RandomSpecies() string {
	return helper.GetRandomElement(r.npcCreationOptions.Species)
}

func (r *RandomizerService) RandomType() string {
	return helper.GetRandomElement(r.npcCreationOptions.NpcTypes)
}

func (r *RandomizerService) RandomSubtype(npcType string) string {
	subtypes := r.npcCreationOptions.NpcSubtypeForTypeMap[npcType]
	return helper.GetRandomElement(subtypes)
}

// abilities
func (r *RandomizerService) RandomAbilities() map[string]int {
	return nil
}

// generating random ID
func (r *RandomizerService) GenerateID() string {
	return helper.GenerateID()
}

// generate trait
func (r *RandomizerService) GenerateTraitDescription() string {
	return r.creationData.GetTraitData(r.RandomTrait()).String()

}

// generating random name
func (r *RandomizerService) GenerateName(species string) string {
	nameDataKey := r.creationData.GetSpeciesNameMap()[species]
	nameData := r.creationData.GetNameData(nameDataKey)
	forename := helper.GetRandomElement(nameData.Forenames)
	surname := helper.GetRandomElement(nameData.Surnames)
	return forename + " " + surname
}

func (r *RandomizerService) GenerateEquipment(npcSubtype string) map[string]string {
	items := make(map[string]string)
	subtype := r.creationData.GetNpcSubtypeData(npcSubtype)
	for category, equipmentOptions := range subtype.EquipmentOptions {
		items[category] = helper.GetRandomElement(equipmentOptions)
	}
	return items
}

func (r *RandomizerService) ApplyTraitStats(trait string) map[string]int {
	return r.applyStats(r.creationData.GetTraitData(trait))
}

func (r *RandomizerService) ApplyTypeStats(npcType string) map[string]int {
	return r.applyStats(r.creationData.GetNpcTypeData(npcType))
}

func (r *RandomizerService) ApplySubtypeStats(npcSubtype string) map[string]int {
	return r.applyStats(r.creationData.GetNpcSubtypeData(npcSubtype))
}

func (r *RandomizerService) applyStats(data interface{}) map[string]int {
	stats := make(map[string]int)
	switch d := data.(type) {
	case model.Trait:
		for key, value := range d.Stats {
			stats[key] += value
		}
	case types.NPCType:
		for _, stat := range d.Stats {
			stats[stat] += rand.Intn(10) + 1
		}
	case types.NPCSubtype:
		for _, stat := range d.Stats {
			stats[stat] += rand.Intn(10) + 1
		}
	}
	return stats
}
