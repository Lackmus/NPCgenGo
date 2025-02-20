package service

import (
	"fmt"

	"github.com/lackmus/npcgengo/shared"
)

type NPCCreationSupplier struct {
	CreationDataService *CreationDataService
	CreationOptions     *NPCCreationOptions
	RandomizerService   *RandomizerService
}

func NewNPCCreationSupplier(loader shared.NPCConfigLoader) *NPCCreationSupplier {
	creationDataService, err := NewCreationDataService(loader)
	if err != nil {
		panic(fmt.Sprintf("Error creating NPCCreationSupplier: %s", err))
	}
	creationOptions := NewNPCCreationOptions(creationDataService)                   // Pass pointer directly
	randomizerService := NewRandomizerService(creationDataService, creationOptions) // Pass pointers

	return &NPCCreationSupplier{
		CreationDataService: creationDataService, // Pass pointers here too
		CreationOptions:     creationOptions,
		RandomizerService:   randomizerService,
	}
}
