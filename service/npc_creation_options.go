// Description: This file contains the NPCCreationOptions struct and the NewNPCCreationOptions function.
package service

import (
	"maps"
	"slices"
)

// NPCCreationOptions provides the options for creating an NPC.
// It is used to provide the options for creating an NPC.
type NPCCreationOptions struct {
	Factions             []string
	Species              []string
	Traits               []string
	NpcTypes             []string
	NpcSubtypeForTypeMap map[string][]string
}

// NewNPCCreationOptions creates a new NPCCreationOptions.
// It returns an error if the data cannot be loaded.
func NewNPCCreationOptions(creationService *CreationDataService) *NPCCreationOptions {
	return &NPCCreationOptions{
		Factions:             slices.Collect(maps.Keys(creationService.GetFactionMap())),
		Species:              slices.Collect(maps.Keys(creationService.GetSpeciesMap())),
		Traits:               slices.Collect(maps.Keys(creationService.GetTraitMap())),
		NpcTypes:             slices.Collect(maps.Keys(creationService.GetNpcTypeMap())),
		NpcSubtypeForTypeMap: creationService.GetNpcSubtypeForTypeMap(),
	}
}
