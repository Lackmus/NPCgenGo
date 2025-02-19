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
	Factions             []string            `json:"factions"`
	Species              []string            `json:"species"`
	Traits               []string            `json:"traits"`
	NpcTypes             []string            `json:"npcTypes"`
	NpcSubtypeForTypeMap map[string][]string `json:"npcSubtypesMap"`
}

// NewCreationDataViewModel creates a new CreationDataViewModel.
func NewNPCCreationOptions(creationservice CreationDataService) *NPCCreationOptions {
	factionKeys := slices.Collect(maps.Keys(creationservice.FactionMap))
	speciesKeys := slices.Collect(maps.Keys(creationservice.SpeciesMap))
	traitKeys := slices.Collect(maps.Keys(creationservice.TraitMap))
	npcTypeKeys := slices.Collect(maps.Keys(creationservice.NpcTypeMap))

	npcSubtypeForTypeMap := creationservice.NpcSubtypeForTypeMap

	return &NPCCreationOptions{
		Factions:             factionKeys,
		Species:              speciesKeys,
		Traits:               traitKeys,
		NpcTypes:             npcTypeKeys,
		NpcSubtypeForTypeMap: npcSubtypeForTypeMap,
	}
}
