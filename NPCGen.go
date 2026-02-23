package main

import (
	"context"
	"fmt"

	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/loader"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/view"
)

// constants for the application
const (
	defaultLocationID = "default"
	npcDatabasePath   = "data/npc_database"
	creationDataPath  = "data/creation_data"
)

type NPCGen struct {
	CreationSupplier  *service.NPCCreationSupplier
	NpcService        *service.NPCService
	NPCListController *controller.NPCListController
}

// NewNPCGen initializes a new NPCGen instance with the provided parameters.
// It creates a new NPCCreationSupplier and NPCListController, and sets up the Fyne view.
func NewNPCGen() (*NPCGen, error) {
	creationSupplier, err := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader(creationDataPath))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NPCCreationSupplier: %w", err)
	}
	npcService, err := service.NewNPCService(context.Background(), loader.NewJSONNPCStorage(npcDatabasePath))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NPCService: %w", err)
	}

	npcListController := controller.NewNPCListController(creationSupplier, npcService, defaultLocationID)

	return &NPCGen{
		CreationSupplier:  creationSupplier,
		NpcService:        npcService,
		NPCListController: npcListController,
	}, nil
}

// New Controller and view instance
func (n *NPCGen) InitNPCListView(locationID string) {
	npcListController := controller.NewNPCListController(n.CreationSupplier, n.NpcService, locationID)
	npcListView := view.NewConsoleView(npcListController)
	npcListController.InitView(npcListView)
	npcListView.Run()
}

// GetCreationOptions returns the creation options from the supplier.
func (n *NPCGen) GetFactions() []string {
	return n.CreationSupplier.CreationOptions.Factions
}
