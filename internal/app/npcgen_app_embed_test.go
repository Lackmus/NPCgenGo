//go:build embeddata

package app

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/lackmus/npcgengo/internal/platform/loader"
)

func TestNewNPCGenWithDataDir_UsesEmbeddedCreationDataWhenFilesystemMissing(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed getting working directory: %v", err)
	}
	tempWD := t.TempDir()
	if err := os.Chdir(tempWD); err != nil {
		t.Fatalf("failed changing working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	app, err := NewNPCGenWithDataDir("")
	if err != nil {
		t.Fatalf("expected embedded creation data fallback to initialize app: %v", err)
	}

	if app == nil || app.CreationSupplier == nil || app.CreationSupplier.CreationOptions == nil {
		t.Fatalf("expected initialized app with creation options")
	}
	if len(app.CreationSupplier.CreationOptions.Factions) == 0 {
		t.Fatalf("expected embedded creation options factions to be non-empty")
	}
}

func TestNewNPCGenWithDataDir_UsesFilesystemCreationDataWhenProvided(t *testing.T) {
	baseDir := t.TempDir()
	creationDir := filepath.Join(baseDir, "creation_data")
	if err := loader.CreateSampleCreationData(creationDir); err != nil {
		t.Fatalf("failed to create sample creation data: %v", err)
	}

	app, err := NewNPCGenWithDataDir(baseDir)
	if err != nil {
		t.Fatalf("expected filesystem creation data to initialize app: %v", err)
	}

	if !slices.Contains(app.CreationSupplier.CreationOptions.Factions, "someFactionID") {
		t.Fatalf("expected sample filesystem faction 'someFactionID' to be present")
	}
}
