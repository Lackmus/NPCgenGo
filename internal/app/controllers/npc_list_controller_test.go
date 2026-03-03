package controllers

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/service"
)

func newControllerForTests(t *testing.T) *NPCListController {
	t.Helper()

	baseDir := t.TempDir()
	creationDir := filepath.Join(baseDir, "creation_data")
	dbDir := filepath.Join(baseDir, "npc_database")

	if err := loader.CreateSampleCreationData(creationDir); err != nil {
		t.Fatalf("failed to create sample creation data: %v", err)
	}

	creationSupplier, err := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader(creationDir))
	if err != nil {
		t.Fatalf("failed to create creation supplier: %v", err)
	}

	npcService, err := service.NewNPCService(context.Background(), loader.NewJSONNPCStorage(dbDir))
	if err != nil {
		t.Fatalf("failed to create NPC service: %v", err)
	}

	return NewNPCListController(creationSupplier, npcService)
}

func TestNPCListController_GetSubtypeFields_ValidSubtype(t *testing.T) {
	controller := newControllerForTests(t)

	stats, items, err := controller.GetSubtypeFields("someCivilianSubtypeID")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(stats, "Str:") || !strings.Contains(stats, "Dex:") {
		t.Fatalf("expected rolled stats to include Str and Dex, got: %q", stats)
	}
	if !strings.Contains(items, "Weapon:") {
		t.Fatalf("expected rolled items to include Weapon, got: %q", items)
	}
}

func TestNPCListController_GetSubtypeFields_InvalidSubtype(t *testing.T) {
	controller := newControllerForTests(t)

	_, _, err := controller.GetSubtypeFields("unknown-subtype")
	if err == nil {
		t.Fatalf("expected error for unknown subtype")
	}
}

func TestNPCListController_GetSpeciesName_ValidSpecies(t *testing.T) {
	controller := newControllerForTests(t)

	name, err := controller.GetSpeciesName("someSpeciesID")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if name != "Alice Smith" {
		t.Fatalf("expected deterministic generated name 'Alice Smith', got %q", name)
	}
}

func TestNPCListController_GetSpeciesName_InvalidSpecies(t *testing.T) {
	controller := newControllerForTests(t)

	_, err := controller.GetSpeciesName("unknown-species")
	if err == nil {
		t.Fatalf("expected error for unknown species")
	}
}
