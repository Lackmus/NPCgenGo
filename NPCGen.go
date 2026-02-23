package npcgengo

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/lackmus/npcgengo/cmd/npcgen/handlers"
	"github.com/lackmus/npcgengo/cmd/npcgen/view"
	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/product/service"
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
	NPCListController *handlers.NPCListController
}

// NewNPCGen initializes a new NPCGen instance using the default on-disk `data` directory.
func NewNPCGen() (*NPCGen, error) {
	return NewNPCGenWithDataDir("")
}

// NewNPCGenWithDataDir initializes a new NPCGen instance using the provided data directory.
// If dataDir is empty, it defaults to the repository-relative "data" directory.
func NewNPCGenWithDataDir(dataDir string) (*NPCGen, error) {
	base := dataDir
	if base == "" {
		base = "data"
	}
	creationPath := filepath.Join(base, "creation_data")
	dbPath := filepath.Join(base, "npc_database")

	creationSupplier, err := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader(creationPath))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NPCCreationSupplier: %w", err)
	}
	npcService, err := service.NewNPCService(context.Background(), loader.NewJSONNPCStorage(dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NPCService: %w", err)
	}

	npcListController := handlers.NewNPCListController(creationSupplier, npcService, defaultLocationID)

	return &NPCGen{
		CreationSupplier:  creationSupplier,
		NpcService:        npcService,
		NPCListController: npcListController,
	}, nil
}

// New Controller and view instance
func (n *NPCGen) InitNPCListView(locationID string) {
	npcListController := handlers.NewNPCListController(n.CreationSupplier, n.NpcService, locationID)
	npcListView := view.NewConsoleView(npcListController)
	npcListController.InitView(npcListView)
	npcListView.Run()
}

// GetCreationOptions returns the creation options from the supplier.
func (n *NPCGen) GetFactions() []string {
	return n.CreationSupplier.CreationOptions.Factions
}
