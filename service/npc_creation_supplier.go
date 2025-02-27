// Description: This file contains the NPCCreationSupplier struct, which encapsulates the services required for NPC creation.
package service

import (
	"log"

	"github.com/lackmus/npcgengo/shared"
)

// NPCCreationSupplier encapsulates the services required for NPC creation.
// It provides the services required for NPC creation.
type NPCCreationSupplier struct {
	CreationDataService *CreationDataService
	CreationOptions     *NPCCreationOptions
	RandomizerService   *RandomizerService
}

// NewNPCCreationSupplier creates a new NPCCreationSupplier.
// It returns an error if the data cannot be loaded.
func NewNPCCreationSupplier(loader shared.NPCConfigLoader) *NPCCreationSupplier {
	c := &NPCCreationSupplier{}
	c.initCreationSupplier(loader)
	return c
}

// initCreationSupplier initializes the services required for NPC creation.
// It panics if the data cannot be loaded.
func (c *NPCCreationSupplier) initCreationSupplier(loader shared.NPCConfigLoader) {
	var err error
	c.CreationDataService, err = NewCreationDataService(loader)
	if err != nil {
		log.Fatalf("Failed to initialize CreationDataService: %v", err)
	}
	c.CreationOptions = NewNPCCreationOptions(c.CreationDataService)
	c.RandomizerService = NewRandomizerService(c.CreationDataService, c.CreationOptions)
}
