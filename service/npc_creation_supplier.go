package service

import (
	"fmt"

	"github.com/lackmus/npcgengo/shared"
)

type NPCCreationSupplier struct {
	CreationDataService CreationDataService
	CreationOptions     NPCCreationOptions
	RandomizerService   RandomizerService
}

func NewNPCCreationSupplier(loader shared.NPCConfigLoader) *NPCCreationSupplier {
	creationDataService, err := NewCreationDataService(loader)
	if err != nil {
		panic(fmt.Sprintf("Error creating NPCCreationSupplier: %s", err))
	}
	creationOptions := NewNPCCreationOptions(*creationDataService)
	randomizerService := NewRandomizerService(*creationDataService, *creationOptions)

	return &NPCCreationSupplier{
		CreationDataService: *creationDataService,
		CreationOptions:     *creationOptions,
		RandomizerService:   *randomizerService,
	}
}
