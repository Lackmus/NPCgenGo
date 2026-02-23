package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/product/service"
	consoleui "github.com/lackmus/npcgengo/ui/console"
)

// NPCGen is the application orchestration layer.
// It wires together creation/loading services and transport-facing controllers.
type NPCGen struct {
	CreationSupplier  *service.NPCCreationSupplier
	NpcService        *service.NPCService
	NPCListController *controllers.NPCListController
}

// NewNPCGen initializes a new NPCGen instance using the default on-disk data directory.
func NewNPCGen() (*NPCGen, error) {
	return NewNPCGenWithDataDir("")
}

// NewNPCGenWithDataDir initializes a new NPCGen instance using the provided data directory.
// If dataDir is empty, it defaults to the repository-relative data directory.
func NewNPCGenWithDataDir(dataDir string) (*NPCGen, error) {
	base := resolveDataDir(dataDir)
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

	npcListController := controllers.NewNPCListController(creationSupplier, npcService)

	return &NPCGen{
		CreationSupplier:  creationSupplier,
		NpcService:        npcService,
		NPCListController: npcListController,
	}, nil
}

func resolveDataDir(dataDir string) string {
	if base := normalizeDataDir(dataDir); base != "" {
		return base
	}

	if env := os.Getenv("NPCGEN_DATA"); env != "" {
		if base := normalizeDataDir(env); base != "" {
			return base
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		if base := findDataDirUp(cwd); base != "" {
			return base
		}
	}

	if executablePath, err := os.Executable(); err == nil {
		if base := findDataDirUp(filepath.Dir(executablePath)); base != "" {
			return base
		}
	}

	return "data"
}

func normalizeDataDir(base string) string {
	if base == "" {
		return ""
	}

	if hasCreationData(base) {
		return base
	}

	candidate := filepath.Join(base, "data")
	if hasCreationData(candidate) {
		return candidate
	}

	return base
}

func findDataDirUp(start string) string {
	current := start

	for {
		if hasCreationData(filepath.Join(current, "data")) {
			return filepath.Join(current, "data")
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return ""
}

func hasCreationData(base string) bool {
	info, err := os.Stat(filepath.Join(base, "creation_data", "factiondata"))
	if err != nil {
		return false
	}
	return info.IsDir()
}

// InitNPCListView creates and starts the console list view.
func (n *NPCGen) InitNPCListView() {
	npcListView := consoleui.NewConsoleView(n.NPCListController)
	n.NPCListController.InitView(npcListView)
	npcListView.Run()
}

// GetFactions returns available faction options.
func (n *NPCGen) GetFactions() []string {
	return n.CreationSupplier.CreationOptions.Factions
}
