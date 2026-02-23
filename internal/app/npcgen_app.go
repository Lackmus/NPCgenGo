package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/lackmus/npcgengo/internal/app/handlers"
	"github.com/lackmus/npcgengo/internal/app/view"
	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/product/service"
)

const defaultLocationID = "default"

// NPCGen is the application orchestration layer.
// It wires together creation/loading services and transport-facing controllers.
type NPCGen struct {
	CreationSupplier  *service.NPCCreationSupplier
	NpcService        *service.NPCService
	NPCListController *handlers.NPCListController
}

// NewNPCGen initializes a new NPCGen instance using the default on-disk data directory.
func NewNPCGen() (*NPCGen, error) {
	return NewNPCGenWithDataDir("")
}

// NewNPCGenWithDataDir initializes a new NPCGen instance using the provided data directory.
// If dataDir is empty, it defaults to the repository-relative data directory.
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

// InitNPCListView creates and starts the console list view for a location.
func (n *NPCGen) InitNPCListView(locationID string) {
	npcListController := handlers.NewNPCListController(n.CreationSupplier, n.NpcService, locationID)
	npcListView := view.NewConsoleView(npcListController)
	npcListController.InitView(npcListView)
	npcListView.Run()
}

// GetFactions returns available faction options.
func (n *NPCGen) GetFactions() []string {
	return n.CreationSupplier.CreationOptions.Factions
}
