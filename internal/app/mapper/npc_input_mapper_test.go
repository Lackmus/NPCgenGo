package mapper

import (
	"path/filepath"
	"testing"

	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/model"
	cp "github.com/lackmus/npcgengo/pkg/model/npc_components"
	"github.com/lackmus/npcgengo/pkg/service"
)

func newBuilderForTests(t *testing.T) *service.NPCBuilder {
	t.Helper()

	baseDir := t.TempDir()
	creationDir := filepath.Join(baseDir, "creation_data")
	if err := loader.CreateSampleCreationData(creationDir); err != nil {
		t.Fatalf("failed to create sample creation data: %v", err)
	}

	supplier, err := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader(creationDir))
	if err != nil {
		t.Fatalf("failed to create NPC creation supplier: %v", err)
	}

	return service.NewNPCBuilder(supplier)
}

func TestToNPCInput_MapsTypedFields(t *testing.T) {
	npc := *model.NewNPC()
	npc.ID = "  npc-123  "
	npc.SetName("  Alice  ")
	npc.SetType(" Civilian ")
	npc.SetSubtype(" someCivilianSubtypeID ")
	npc.SetSpecies(" someSpeciesID ")
	npc.SetFaction(" someFactionID ")
	npc.SetTrait(" someTraitID ")
	npc.SetStats("  STR:1  ")
	npc.SetItems("  Knife  ")
	npc.SetNotes("  Keeps journal entries.  ")

	got := ToNPCInput(npc)

	if got.ID != "npc-123" || got.Name != "Alice" || got.Type != "Civilian" || got.Subtype != "someCivilianSubtypeID" || got.Species != "someSpeciesID" || got.Faction != "someFactionID" || got.Trait != "someTraitID" || got.Stats != "STR:1" || got.Items != "Knife" || got.Notes != "Keeps journal entries." {
		t.Fatalf("unexpected DTO mapping: %+v", got)
	}
}

func TestToNPCInputs_MapsSlice(t *testing.T) {
	first := *model.NewNPC()
	first.ID = "1"
	first.SetName("One")

	second := *model.NewNPC()
	second.ID = "2"
	second.SetName("Two")

	got := ToNPCInputs([]model.NPC{first, second})

	if len(got) != 2 {
		t.Fatalf("expected 2 DTOs, got %d", len(got))
	}
	if got[0].ID != "1" || got[0].Name != "One" || got[1].ID != "2" || got[1].Name != "Two" {
		t.Fatalf("unexpected DTO slice mapping: %+v", got)
	}
}

func TestToModelNPC_BuildsFromTypedInput(t *testing.T) {
	builder := newBuilderForTests(t)
	input := NPCInput{
		ID:      "npc-42",
		Name:    "Alice Smith",
		Type:    "Civilian",
		Subtype: "someCivilianSubtypeID",
		Species: "someSpeciesID",
		Faction: "someFactionID",
		Trait:   "someTraitID",
		Stats:   "STR:2, DEX:1",
		Items:   "Fists",
		Notes:   "Scouting routes near river crossing",
	}

	npc, err := ToModelNPC(input, builder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if npc.ID != "npc-42" {
		t.Fatalf("expected ID npc-42, got %q", npc.ID)
	}
	if npc.GetComponent(cp.CompName) != input.Name ||
		npc.GetComponent(cp.CompType) != input.Type ||
		npc.GetComponent(cp.CompSubtype) != input.Subtype ||
		npc.GetComponent(cp.CompSpecies) != input.Species ||
		npc.GetComponent(cp.CompFaction) != input.Faction ||
		npc.GetComponent(cp.CompTrait) != input.Trait ||
		npc.GetComponent(cp.CompStats) != input.Stats ||
		npc.GetComponent(cp.CompItems) != input.Items ||
		npc.GetComponent(cp.CompNotes) != input.Notes {
		t.Fatalf("unexpected model mapping from DTO: %+v", npc)
	}
}

func TestToModelNPCWithOriginal_PreservesMissingFields(t *testing.T) {
	builder := newBuilderForTests(t)

	original := *model.NewNPC()
	original.ID = "npc-9"
	original.SetName("Old Name")
	original.SetType("Civilian")
	original.SetSubtype("someCivilianSubtypeID")
	original.SetSpecies("someSpeciesID")
	original.SetFaction("someFactionID")
	original.SetTrait("someTraitID")
	original.SetStats("OLD-STATS")
	original.SetItems("OLD-ITEMS")
	original.SetNotes("OLD-NOTES")

	input := NPCInput{
		ID:    "npc-9",
		Name:  "New Name",
		Trait: "someTraitID",
	}

	npc, err := ToModelNPCWithOriginal(input, builder, &original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if npc.ID != "npc-9" {
		t.Fatalf("expected ID npc-9, got %q", npc.ID)
	}
	if npc.GetComponent(cp.CompName) != "New Name" {
		t.Fatalf("expected Name to update, got %q", npc.GetComponent(cp.CompName))
	}
	if npc.GetComponent(cp.CompType) != "Civilian" ||
		npc.GetComponent(cp.CompSubtype) != "someCivilianSubtypeID" ||
		npc.GetComponent(cp.CompSpecies) != "someSpeciesID" ||
		npc.GetComponent(cp.CompFaction) != "someFactionID" ||
		npc.GetComponent(cp.CompTrait) != "someTraitID" ||
		npc.GetComponent(cp.CompStats) != "OLD-STATS" ||
		npc.GetComponent(cp.CompItems) != "OLD-ITEMS" ||
		npc.GetComponent(cp.CompNotes) != "OLD-NOTES" {
		t.Fatalf("expected missing fields to be preserved, got npc: %+v", npc)
	}
}
