package loader

import (
	"context"
	"testing"

	npcgendata "github.com/lackmus/npcgengo/data"
)

func TestFSNPCConfigLoader_LoadsEmbeddedCreationData(t *testing.T) {
	configLoader := NewFSNPCConfigLoader(npcgendata.CreationDataFS(), "creation_data")
	ctx := context.Background()

	factions, err := configLoader.LoadFactionMap(ctx)
	if err != nil {
		t.Fatalf("failed loading embedded factions: %v", err)
	}
	if len(factions) == 0 {
		t.Fatalf("expected embedded factions to be non-empty")
	}

	species, err := configLoader.LoadSpeciesMap(ctx)
	if err != nil {
		t.Fatalf("failed loading embedded species: %v", err)
	}
	if len(species) == 0 {
		t.Fatalf("expected embedded species to be non-empty")
	}

	traits, err := configLoader.LoadTraitMap(ctx)
	if err != nil {
		t.Fatalf("failed loading embedded traits: %v", err)
	}
	if len(traits) == 0 {
		t.Fatalf("expected embedded traits to be non-empty")
	}

	names, err := configLoader.LoadNameMap(ctx)
	if err != nil {
		t.Fatalf("failed loading embedded names: %v", err)
	}
	if len(names) == 0 {
		t.Fatalf("expected embedded names to be non-empty")
	}

	subtypes, err := configLoader.LoadNPCSubtypeMaps(ctx)
	if err != nil {
		t.Fatalf("failed loading embedded subtype maps: %v", err)
	}
	if len(subtypes) == 0 {
		t.Fatalf("expected embedded subtype maps to be non-empty")
	}
}
