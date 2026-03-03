package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/model"
)

// TestNewNPCService_AllValid tests that NewNPCService successfully loads all valid NPCs from storage.
func TestNewNPCService_AllValid(t *testing.T) {
	dir, err := os.MkdirTemp("", "npcstorage_test_valid_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	storage := loader.NewJSONNPCStorage(dir)

	// create and save two NPCs
	npc1 := model.NewNPC()
	npc1.ID = "1"
	if err := storage.SaveNPC(context.Background(), *npc1); err != nil {
		t.Fatalf("failed to save npc1: %v", err)
	}

	npc2 := model.NewNPC()
	npc2.ID = "2"
	if err := storage.SaveNPC(context.Background(), *npc2); err != nil {
		t.Fatalf("failed to save npc2: %v", err)
	}

	s, err := NewNPCService(context.Background(), storage)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(s.npcs) != 2 {
		t.Fatalf("expected 2 npcs loaded, got %d", len(s.npcs))
	}
	if _, ok := s.npcs["1"]; !ok {
		t.Fatalf("expected npc '1' to be present")
	}
	if _, ok := s.npcs["2"]; !ok {
		t.Fatalf("expected npc '2' to be present")
	}
}

func TestNewNPCService_PartialLoadOnCorruptFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "npcstorage_test_partial_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	storage := loader.NewJSONNPCStorage(dir)

	// create and save one valid NPC
	npc1 := model.NewNPC()
	npc1.ID = "1"
	if err := storage.SaveNPC(context.Background(), *npc1); err != nil {
		t.Fatalf("failed to save npc1: %v", err)
	}

	// create a corrupt JSON file for id '2'
	corruptPath := filepath.Join(dir, "2.json")
	if err := os.WriteFile(corruptPath, []byte("not-json"), 0644); err != nil {
		t.Fatalf("failed to write corrupt file: %v", err)
	}

	s, err := NewNPCService(context.Background(), storage)
	if err == nil {
		t.Fatalf("expected error due to corrupt file, got nil")
	}
	// partial data should be present
	if _, ok := s.npcs["1"]; !ok {
		t.Fatalf("expected npc '1' to be present despite error")
	}
	if _, ok := s.npcs["2"]; ok {
		t.Fatalf("did not expect npc '2' to be present")
	}
}
