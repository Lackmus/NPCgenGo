package service

import (
	"path/filepath"
	"testing"

	"github.com/lackmus/npcgengo/internal/platform/loader"
)

func newCreationSupplierForSeedTests(t *testing.T) *NPCCreationSupplier {
	t.Helper()

	baseDir := t.TempDir()
	creationDir := filepath.Join(baseDir, "creation_data")
	if err := loader.CreateSampleCreationData(creationDir); err != nil {
		t.Fatalf("failed to create sample creation data: %v", err)
	}

	supplier, err := NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader(creationDir))
	if err != nil {
		t.Fatalf("failed to create supplier: %v", err)
	}
	return supplier
}

func TestCreateNPCWithOptionsAndSeed_SameSeedIsReproducible(t *testing.T) {
	supplier := newCreationSupplierForSeedTests(t)

	first, err := CreateNPCWithOptionsAndSeed("random", "random", 1337, supplier)
	if err != nil {
		t.Fatalf("first generation failed: %v", err)
	}
	second, err := CreateNPCWithOptionsAndSeed("random", "random", 1337, supplier)
	if err != nil {
		t.Fatalf("second generation failed: %v", err)
	}

	if first.String() != second.String() {
		t.Fatalf("expected same output for same seed\nfirst: %s\nsecond: %s", first.String(), second.String())
	}
}

func TestCreateNPCWithOptionsAndSeed_DifferentSeedCanChangeOutput(t *testing.T) {
	supplier := newCreationSupplierForSeedTests(t)

	first, err := CreateNPCWithOptionsAndSeed("random", "random", 1337, supplier)
	if err != nil {
		t.Fatalf("first generation failed: %v", err)
	}
	second, err := CreateNPCWithOptionsAndSeed("random", "random", 1338, supplier)
	if err != nil {
		t.Fatalf("second generation failed: %v", err)
	}

	if first.String() == second.String() {
		t.Fatalf("expected different output for different seeds, but got same output: %s", first.String())
	}
}
