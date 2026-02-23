// Description: This file contains the NPCCreationSupplier struct, which encapsulates the services required for NPC creation.
package service

import (
	"context"
	"fmt"

	"github.com/lackmus/npcgengo/shared"
)

type NPCCreationSupplier struct {
	CreationDataService *CreationDataService
	CreationOptions     *NPCCreationOptions
	RandomizerService   *RandomizerService
}

func NewNPCCreationSupplier(loader shared.NPCConfigLoader) (*NPCCreationSupplier, error) {
	c := &NPCCreationSupplier{}
	if err := c.initCreationSupplier(loader); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *NPCCreationSupplier) initCreationSupplier(loader shared.NPCConfigLoader) error {
	var err error
	// use background context for initialization; callers may opt to initialize differently
	c.CreationDataService, err = NewCreationDataService(context.Background(), loader)
	if err != nil {
		return fmt.Errorf("failed to initialize CreationDataService: %w", err)
	}
	c.CreationOptions = NewNPCCreationOptions(c.CreationDataService)
	c.RandomizerService = NewRandomizerService(c.CreationDataService, c.CreationOptions)
	return nil
}
