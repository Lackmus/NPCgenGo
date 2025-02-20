package service

import (
	"maps"
	"slices"
)

// ======================================================
// NPC CREATION OPTIONS
// ======================================================

// CreationDataViewModel holds only the necessary string lists for the UI.
type NPCCreationOptions struct {
	Factions             []string
	Species              []string
	Traits               []string
	NpcTypes             []string
	NpcSubtypeForTypeMap map[string][]string
}

// NewCreationDataViewModel creates a new CreationDataViewModel.
func NewNPCCreationOptions(creationservice *CreationDataService) *NPCCreationOptions {
	factionKeys := slices.Collect(maps.Keys(creationservice.GetFactionMap()))
	speciesKeys := slices.Collect(maps.Keys(creationservice.GetSpeciesMap()))
	traitKeys := slices.Collect(maps.Keys(creationservice.GetTraitMap()))
	npcTypeKeys := slices.Collect(maps.Keys(creationservice.GetNpcTypeMap()))

	npcSubtypeForTypeMap := creationservice.GetNpcSubtypeForTypeMap()

	return &NPCCreationOptions{
		Factions:             factionKeys,
		Species:              speciesKeys,
		Traits:               traitKeys,
		NpcTypes:             npcTypeKeys,
		NpcSubtypeForTypeMap: npcSubtypeForTypeMap,
	}
}
