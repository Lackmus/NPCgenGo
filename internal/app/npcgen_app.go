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

const defaultDataDir = "data"

// NPCGen is the application orchestration layer.
// It wires together creation/loading services and transport-facing controllers.
type NPCGen struct {
	CreationSupplier  *service.NPCCreationSupplier
	NPCService        *service.NPCService
	NPCListController *controllers.NPCListController
}

// NewNPCGen initializes a new NPCGen instance using the default on-disk data directory.
func NewNPCGen() (*NPCGen, error) {
	return NewNPCGenWithDataDir("")
}

// NewNPCGenWithDataDir initializes a new NPCGen instance using the provided data directory.
// If dataDir is empty or invalid, it resolves using environment/discovery fallbacks.
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
		NPCService:        npcService,
		NPCListController: npcListController,
	}, nil
}

// resolveDataDir determines the appropriate data directory to use based on the provided input, environment variables, and directory structure.
func resolveDataDir(dataDir string) string {
	if base := normalizeDataDir(dataDir); base != "" {
		return base
	}

	// Check NPCGEN_DATA environment variable as a fallback if dataDir is not provided or valid.
	if env := os.Getenv("NPCGEN_DATA"); env != "" {
		if base := normalizeDataDir(env); base != "" {
			return base
		}
	}

	// Attempt to find the data directory by traversing up from the current working directory.
	if cwd, err := os.Getwd(); err == nil {
		if base := findDataDirUp(cwd); base != "" {
			return base
		}
	}

	// As a last resort, check the directory of the executable itself, which is useful for bundled applications.
	if executablePath, err := os.Executable(); err == nil {
		if base := findDataDirUp(filepath.Dir(executablePath)); base != "" {
			return base
		}
	}

	return defaultDataDir
}

// normalizeDataDir checks if the provided base path directly contains the expected creation data structure.
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

	return ""
}

// findDataDirUp traverses up the directory tree from the starting point to find a directory containing the expected creation data structure.
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

// hasCreationData checks if the given base directory contains the expected creation data structure.
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
