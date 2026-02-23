package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/product/model"
)

func newValidationServiceForTests(t *testing.T) *NPCValidationService {
	t.Helper()

	baseDir := t.TempDir()
	if err := loader.CreateSampleCreationData(filepath.Join(baseDir, "creation_data")); err != nil {
		t.Fatalf("failed to create sample creation data: %v", err)
	}

	configLoader := loader.NewJSONNPCConfigLoader(filepath.Join(baseDir, "creation_data"))
	creationData, err := NewCreationDataService(context.Background(), configLoader)
	if err != nil {
		t.Fatalf("failed to init creation data service: %v", err)
	}

	return NewNPCValidationService(creationData)
}

func TestNPCValidationService_ValidateNPC_Valid(t *testing.T) {
	validator := newValidationServiceForTests(t)

	npc := *model.NewNPC()
	npc.SetType("Civilian")
	npc.SetSubtype("someCivilianSubtypeID")
	npc.SetSpecies("someSpeciesID")
	npc.SetFaction("someFactionID")
	npc.SetTrait("someTraitID")

	if err := validator.ValidateNPC(npc); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
}

func TestNPCValidationService_ValidateNPC_InvalidType(t *testing.T) {
	validator := newValidationServiceForTests(t)

	npc := *model.NewNPC()
	npc.SetType("UnknownType")

	if err := validator.ValidateNPC(npc); err == nil {
		t.Fatalf("expected validation error for invalid type")
	}
}

func TestNPCValidationService_ValidateNPC_SubtypeTypeMismatch(t *testing.T) {
	validator := newValidationServiceForTests(t)

	npc := *model.NewNPC()
	npc.SetType("Military")
	npc.SetSubtype("someCivilianSubtypeID")

	if err := validator.ValidateNPC(npc); err == nil {
		t.Fatalf("expected validation error for subtype/type mismatch")
	}
}

func TestNPCValidationService_ValidateNPC_InvalidTrait(t *testing.T) {
	validator := newValidationServiceForTests(t)

	npc := *model.NewNPC()
	npc.SetTrait("someTraitID, UnknownTrait")

	if err := validator.ValidateNPC(npc); err == nil {
		t.Fatalf("expected validation error for invalid trait")
	}
}
