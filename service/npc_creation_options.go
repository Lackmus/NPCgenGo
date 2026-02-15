// Description: This file contains the NPCCreationOptions struct and the NewNPCCreationOptions function.
package service

import (
	"maps"
	"slices"
)

type NPCCreationOptions struct {
	Factions             []string
	Species              []string
	Traits               []string
	NpcTypes             []string
	NpcSubtypeForTypeMap map[string][]string
}

func NewNPCCreationOptions(creationService *CreationDataService) *NPCCreationOptions {
	return &NPCCreationOptions{
		Factions:             slices.Collect(maps.Keys(creationService.GetFactionMap())),
		Species:              slices.Collect(maps.Keys(creationService.GetSpeciesMap())),
		Traits:               slices.Collect(maps.Keys(creationService.GetTraitMap())),
		NpcTypes:             slices.Collect(maps.Keys(creationService.GetNpcTypeMap())),
		NpcSubtypeForTypeMap: creationService.GetNpcSubtypeForTypeMap(),
	}
}
