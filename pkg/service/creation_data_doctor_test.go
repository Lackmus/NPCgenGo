package service

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lackmus/npcgengo/internal/platform/loader"
	c "github.com/lackmus/npcgengo/pkg/model/npc_components"
)

type mockConfigLoader struct {
	factions map[string]c.Faction
	species  map[string]c.Species
	traits   map[string]c.Trait
	names    map[string]c.NameData
	subtypes map[string]map[string]c.NPCSubtype
}

func (m mockConfigLoader) LoadFactionMap(context.Context) (map[string]c.Faction, error) {
	return m.factions, nil
}

func (m mockConfigLoader) LoadSpeciesMap(context.Context) (map[string]c.Species, error) {
	return m.species, nil
}

func (m mockConfigLoader) LoadTraitMap(context.Context) (map[string]c.Trait, error) {
	return m.traits, nil
}

func (m mockConfigLoader) LoadNameMap(context.Context) (map[string]c.NameData, error) {
	return m.names, nil
}

func (m mockConfigLoader) LoadNPCSubtypeMaps(context.Context) (map[string]map[string]c.NPCSubtype, error) {
	return m.subtypes, nil
}

func TestValidateCreationData_ValidSampleData(t *testing.T) {
	baseDir := t.TempDir()
	creationDir := filepath.Join(baseDir, "creation_data")
	if err := loader.CreateSampleCreationData(creationDir); err != nil {
		t.Fatalf("failed to create sample creation data: %v", err)
	}

	err := ValidateCreationData(context.Background(), loader.NewJSONNPCConfigLoader(creationDir))
	if err != nil {
		t.Fatalf("expected valid sample data, got error: %v", err)
	}
}

func TestValidateCreationData_InvalidReferences(t *testing.T) {
	loader := mockConfigLoader{
		factions: map[string]c.Faction{
			"someFactionID": {Name: "someFactionID", SpeciesList: []string{"missingSpecies"}},
		},
		species: map[string]c.Species{
			"someSpeciesID": {Name: "someSpeciesID", NameSource: "missingName"},
		},
		traits: map[string]c.Trait{
			"someTraitID": {Name: "someTraitID", Opposes: "missingTrait"},
		},
		names: map[string]c.NameData{
			"someNameID": {Name: "someNameID", Forenames: []string{"Alice"}, Surnames: []string{"Smith"}},
		},
		subtypes: map[string]map[string]c.NPCSubtype{
			"Civilian": {
				"someCivilianSubtypeID": {
					Name:             "someCivilianSubtypeID",
					NpcTypeName:      "Military",
					Stats:            []string{},
					EquipmentOptions: map[string][]string{"Weapon": {}},
				},
			},
		},
	}

	err := ValidateCreationData(context.Background(), loader)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	message := err.Error()
	expectedParts := []string{
		"references unknown species",
		"references unknown name source",
		"references unknown opposing trait",
		"has type \"Military\" but is listed under \"Civilian\"",
		"has no stats",
		"empty equipment option list",
	}
	for _, part := range expectedParts {
		if !strings.Contains(message, part) {
			t.Fatalf("expected error to contain %q, got: %s", part, message)
		}
	}
}
