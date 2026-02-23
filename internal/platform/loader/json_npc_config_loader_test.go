package loader

import (
	"context"
	"testing"
)

func TestJSONNPCConfigLoader_LoadFactionMap(t *testing.T) {
	dir := t.TempDir()
	if err := CreateSampleCreationData(dir); err != nil {
		t.Fatal(err)
	}
	loader := NewJSONNPCConfigLoader(dir)
	ctx := context.Background()

	// Test loading faction data
	factions, err := loader.LoadFactionMap(ctx)

	// Check if error occurred
	if err != nil {
		t.Errorf("unexpected error loading factions: %v", err)
	}

	// Check if the factions map is not empty
	if len(factions) == 0 {
		t.Errorf("expected factions map to be non-empty, got empty map")
	}

	// Check for a specific faction ID (modify as needed)
	if _, exists := factions["someFactionID"]; !exists {
		t.Errorf("expected faction ID 'someFactionID' to exist")
	}
}

func TestJSONNPCConfigLoader_LoadSpeciesMap(t *testing.T) {
	dir := t.TempDir()
	if err := CreateSampleCreationData(dir); err != nil {
		t.Fatal(err)
	}
	loader := NewJSONNPCConfigLoader(dir)
	ctx := context.Background()

	// Test loading species data
	species, err := loader.LoadSpeciesMap(ctx)

	// Check if error occurred
	if err != nil {
		t.Errorf("unexpected error loading species: %v", err)
	}

	// Check if the species map is not empty
	if len(species) == 0 {
		t.Errorf("expected species map to be non-empty, got empty map")
	}

	// Check for a specific species ID (modify as needed)
	if _, exists := species["someSpeciesID"]; !exists {
		t.Errorf("expected species ID 'someSpeciesID' to exist")
	}
}

func TestJSONNPCConfigLoader_LoadTraitMap(t *testing.T) {
	dir := t.TempDir()
	if err := CreateSampleCreationData(dir); err != nil {
		t.Fatal(err)
	}
	loader := NewJSONNPCConfigLoader(dir)
	ctx := context.Background()

	// Test loading trait data
	traits, err := loader.LoadTraitMap(ctx)

	// Check if error occurred
	if err != nil {
		t.Errorf("unexpected error loading traits: %v", err)
	}

	// Check if the traits map is not empty
	if len(traits) == 0 {
		t.Errorf("expected traits map to be non-empty, got empty map")
	}

	// Check for a specific trait ID
	if _, exists := traits["someTraitID"]; !exists {
		t.Errorf("expected trait ID 'someTraitID' to exist")
	}
}

func TestJSONNPCConfigLoader_LoadNameMap(t *testing.T) {
	dir := t.TempDir()
	if err := CreateSampleCreationData(dir); err != nil {
		t.Fatal(err)
	}
	loader := NewJSONNPCConfigLoader(dir)
	ctx := context.Background()

	// Test loading name data
	names, err := loader.LoadNameMap(ctx)

	// Check if error occurred
	if err != nil {
		t.Errorf("unexpected error loading names: %v", err)
	}

	// Check if the names map is not empty
	if len(names) == 0 {
		t.Errorf("expected names map to be non-empty, got empty map")
	}

	// Check for a specific name ID (modify as needed)
	if _, exists := names["someNameID"]; !exists {
		t.Errorf("expected name ID 'someNameID' to exist")
	}
}

func TestJSONNPCConfigLoader_LoadNpcSubtypeMaps(t *testing.T) {
	dir := t.TempDir()
	if err := CreateSampleCreationData(dir); err != nil {
		t.Fatal(err)
	}
	loader := NewJSONNPCConfigLoader(dir)
	ctx := context.Background()

	// Test loading dynamic subtype data
	subtypeMaps, err := loader.LoadNpcSubtypeMaps(ctx)

	// Check if error occurred
	if err != nil {
		t.Errorf("unexpected error loading subtype maps: %v", err)
	}

	if len(subtypeMaps) == 0 {
		t.Errorf("expected subtype maps to be non-empty, got empty map")
	}

	civilianSubtypes, ok := subtypeMaps["Civilian"]
	if !ok {
		t.Fatalf("expected subtype map for 'Civilian' to exist")
	}

	if _, exists := civilianSubtypes["someCivilianSubtypeID"]; !exists {
		t.Errorf("expected civilian subtype ID 'someCivilianSubtypeID' to exist")
	}

	militarySubtypes, ok := subtypeMaps["Military"]
	if !ok {
		t.Fatalf("expected subtype map for 'Military' to exist")
	}

	if _, exists := militarySubtypes["someMilitarySubtypeID"]; !exists {
		t.Errorf("expected military subtype ID 'someMilitarySubtypeID' to exist")
	}
}
