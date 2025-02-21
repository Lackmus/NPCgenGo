package service

import (
	"log"

	"github.com/lackmus/npcgengo/shared"
)

// NPCCreationSupplier encapsulates the services required for NPC creation.
type NPCCreationSupplier struct {
	CreationDataService *CreationDataService
	CreationOptions     *NPCCreationOptions
	RandomizerService   *RandomizerService
}

// NewNPCCreationSupplier initializes an NPCCreationSupplier, logging an error if initialization fails.
func NewNPCCreationSupplier(loader shared.NPCConfigLoader) *NPCCreationSupplier {
	creationDataService, err := NewCreationDataService(loader)
	if err != nil {
		log.Fatalf("Failed to create NPCCreationSupplier: %v", err) // Logs and exits the program
	}

	creationOptions := NewNPCCreationOptions(creationDataService)
	randomizerService := NewRandomizerService(creationDataService, creationOptions)

	return &NPCCreationSupplier{
		CreationDataService: creationDataService,
		CreationOptions:     creationOptions,
		RandomizerService:   randomizerService,
	}
}
