package loader

import (
	"context"
	"strings"
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

	// Check for a specific trait ID (modify as needed). Trait GetName includes Opposes text,
	// so match by prefix instead of map key.
	found := false
	for _, tr := range traits {
		if strings.HasPrefix(tr.GetName(), "someTraitID") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected trait with name prefix 'someTraitID' to exist")
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

func TestJSONNPCConfigLoader_LoadNpcCivilianSubtypeMap(t *testing.T) {
	dir := t.TempDir()
	if err := CreateSampleCreationData(dir); err != nil {
		t.Fatal(err)
	}
	loader := NewJSONNPCConfigLoader(dir)
	ctx := context.Background()

	// Test loading civilian subtype data
	subtypes, err := loader.LoadNpcCivilianSubtypeMap(ctx)

	// Check if error occurred
	if err != nil {
		t.Errorf("unexpected error loading civilian subtypes: %v", err)
	}

	// Check if the civilian subtype map is not empty
	if len(subtypes) == 0 {
		t.Errorf("expected civilian subtypes map to be non-empty, got empty map")
	}

	// Check for a specific civilian subtype ID (modify as needed)
	if _, exists := subtypes["someCivilianSubtypeID"]; !exists {
		t.Errorf("expected civilian subtype ID 'someCivilianSubtypeID' to exist")
	}
}

func TestJSONNPCConfigLoader_LoadNpcMilitarySubtypeMap(t *testing.T) {
	dir := t.TempDir()
	if err := CreateSampleCreationData(dir); err != nil {
		t.Fatal(err)
	}
	loader := NewJSONNPCConfigLoader(dir)
	ctx := context.Background()

	// Test loading military subtype data
	subtypes, err := loader.LoadNpcMilitarySubtypeMap(ctx)

	// Check if error occurred
	if err != nil {
		t.Errorf("unexpected error loading military subtypes: %v", err)
	}

	// Check if the military subtype map is not empty
	if len(subtypes) == 0 {
		t.Errorf("expected military subtypes map to be non-empty, got empty map")
	}

	// Check for a specific military subtype ID (modify as needed)
	if _, exists := subtypes["someMilitarySubtypeID"]; !exists {
		t.Errorf("expected military subtype ID 'someMilitarySubtypeID' to exist")
	}
}
