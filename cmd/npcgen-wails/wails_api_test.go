package main

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/mapper"
	"github.com/lackmus/npcgengo/pkg/model"
	cp "github.com/lackmus/npcgengo/pkg/model/npc_components"
	"github.com/lackmus/npcgengo/pkg/service"
)

func newWailsAPIForTests(t *testing.T) *WailsAPI {
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

	controller := controllers.NewNPCListController(creationSupplier, npcService)
	return NewWailsAPI(controller)
}

func TestWailsAPI_RerollEndpoints(t *testing.T) {
	api := newWailsAPIForTests(t)

	subtypeRoll, err := api.RollSubtypeFields("someCivilianSubtypeID")
	if err != nil {
		t.Fatalf("expected no subtype roll error, got: %v", err)
	}
	if !strings.Contains(subtypeRoll.Stats, "Str:") || !strings.Contains(subtypeRoll.Items, "Weapon:") {
		t.Fatalf("unexpected subtype roll payload: %+v", subtypeRoll)
	}

	name, err := api.RollSpeciesName("someSpeciesID")
	if err != nil {
		t.Fatalf("expected no species roll error, got: %v", err)
	}
	if name != "Alice Smith" {
		t.Fatalf("expected deterministic generated name 'Alice Smith', got %q", name)
	}
}

func TestWailsAPI_SaveNPC_ValidationError(t *testing.T) {
	api := newWailsAPIForTests(t)

	existing := *model.NewNPC()
	existing.ID = "existing-1"
	existing.SetComponent(cp.CompName, "Existing Name")
	existing.SetComponent(cp.CompType, "Civilian")
	existing.SetComponent(cp.CompSubtype, "someCivilianSubtypeID")
	existing.SetComponent(cp.CompSpecies, "someSpeciesID")
	existing.SetComponent(cp.CompFaction, "someFactionID")
	existing.SetComponent(cp.CompTrait, "someTraitID")
	existing.SetComponent(cp.CompStats, "STR: 2")
	existing.SetComponent(cp.CompItems, "Weapon: Fists")
	api.npcController.AddNPC(existing)

	_, err := api.SaveNPC(mapper.NPCInput{
		ID:      "existing-1",
		Name:    "Updated Name",
		Type:    "Civilian",
		Subtype: "someCivilianSubtypeID",
		Species: "someSpeciesID",
		Faction: "invalid-faction",
		Trait:   "someTraitID",
		Stats:   "STR: 3",
		Items:   "Weapon: Fists",
	})

	if err == nil {
		t.Fatalf("expected validation error for invalid faction")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "invalid faction") {
		t.Fatalf("expected invalid faction error, got: %v", err)
	}
}

func TestWailsAPI_SaveNPC_ValidUpdate(t *testing.T) {
	api := newWailsAPIForTests(t)

	existing := *model.NewNPC()
	existing.ID = "existing-2"
	existing.SetComponent(cp.CompName, "Existing Name")
	existing.SetComponent(cp.CompType, "Civilian")
	existing.SetComponent(cp.CompSubtype, "someCivilianSubtypeID")
	existing.SetComponent(cp.CompSpecies, "someSpeciesID")
	existing.SetComponent(cp.CompFaction, "someFactionID")
	existing.SetComponent(cp.CompTrait, "someTraitID")
	existing.SetComponent(cp.CompStats, "STR: 2")
	existing.SetComponent(cp.CompItems, "Weapon: Fists")
	api.npcController.AddNPC(existing)

	saved, err := api.SaveNPC(mapper.NPCInput{
		ID:      "existing-2",
		Name:    "Updated Name",
		Type:    "Civilian",
		Subtype: "someCivilianSubtypeID",
		Species: "someSpeciesID",
		Faction: "someFactionID",
		Trait:   "someTraitID",
		Stats:   "STR: 3",
		Items:   "Weapon: Fists",
		Notes:   "Tracks patrol timings",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if saved.ID != "existing-2" || saved.Name != "Updated Name" || saved.Notes != "Tracks patrol timings" {
		t.Fatalf("unexpected save response: %+v", saved)
	}
}
