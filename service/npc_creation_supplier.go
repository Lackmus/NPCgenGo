// Description: This file contains the NPCCreationSupplier struct, which encapsulates the services required for NPC creation.
package service

import (
	"log"

	"github.com/lackmus/npcgengo/shared"
)

type NPCCreationSupplier struct {
	CreationDataService *CreationDataService
	CreationOptions     *NPCCreationOptions
	RandomizerService   *RandomizerService
}

func NewNPCCreationSupplier(loader shared.NPCConfigLoader) *NPCCreationSupplier {
	c := &NPCCreationSupplier{}
	c.initCreationSupplier(loader)
	return c
}

func (c *NPCCreationSupplier) initCreationSupplier(loader shared.NPCConfigLoader) {
	var err error
	c.CreationDataService, err = NewCreationDataService(loader)
	if err != nil {
		log.Fatalf("Failed to initialize CreationDataService: %v", err)
	}
	c.CreationOptions = NewNPCCreationOptions(c.CreationDataService)
	c.RandomizerService = NewRandomizerService(c.CreationDataService, c.CreationOptions)
}
