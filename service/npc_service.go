package service

import (
	"fmt"
	"log"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

// NPCService manages NPCs and loads them from a storage backend.
// Note: If concurrent access is expected, consider adding a sync.RWMutex.
type NPCService struct {
	loader shared.NPCStorage
	npcs   map[string]model.NPC // Using a map for faster lookup instead of a slice.
}

// NewNPCService creates a new NPCService and loads existing NPCs.
// In case of an error, it logs the error and initializes an empty map.
func NewNPCService(loader shared.NPCStorage) *NPCService {
	npcMap, err := loader.LoadAllNPC()
	if err != nil {
		log.Printf("Error loading NPCs: %v", err)
		npcMap = make(map[string]model.NPC)
	}

	return &NPCService{
		loader: loader,
		npcs:   npcMap,
	}
}

// AddNPC adds a new NPC and immediately saves it.
func (s *NPCService) AddNPC(npc model.NPC) {
	s.npcs[npc.ID] = npc
	if err := s.loader.SaveNPC(npc); err != nil {
		log.Printf("Error saving NPC (ID %s): %v", npc.ID, err)
	}
}

// GetAllNPCs returns all NPCs as a slice.
func (s *NPCService) GetAllNPC() []model.NPC {
	npcList := make([]model.NPC, 0, len(s.npcs))
	for _, npc := range s.npcs {
		npcList = append(npcList, npc)
	}
	return npcList
}

// GetNPCByID returns the NPC with the specified ID.
func (s *NPCService) GetNPCByID(id string) (model.NPC, error) {
	npc, found := s.npcs[id]
	if !found {
		return model.NPC{}, fmt.Errorf("NPC with ID %s not found", id)
	}
	return npc, nil
}

// UpdateNPC updates an existing NPC and saves it.
func (s *NPCService) UpdateNPC(updatedNPC model.NPC) error {
	id := updatedNPC.ID
	if _, found := s.npcs[id]; !found {
		return fmt.Errorf("NPC with ID %s not found", id)
	}
	s.npcs[id] = updatedNPC
	if err := s.loader.SaveNPC(updatedNPC); err != nil {
		return fmt.Errorf("failed to save updated NPC: %w", err)
	}
	return nil
}

// DeleteNPC removes an NPC from the map and deletes it from the storage.
func (s *NPCService) DeleteNPC(id string) error {
	if _, found := s.npcs[id]; !found {
		return fmt.Errorf("NPC with ID %s not found", id)
	}
	delete(s.npcs, id)
	return s.loader.DeleteNPC(id)
}

// DeleteAllNPC deletes all NPCs from the storage and clears the map.
func (s *NPCService) DeleteAllNPC() {
	if err := s.loader.DeleteAllNPC(); err != nil {
		log.Printf("Error deleting all NPCs: %v", err)
	}
	s.npcs = make(map[string]model.NPC)
}

// CountNPC returns the number of NPCs.
func (s *NPCService) CountNPC() int {
	return len(s.npcs)
}

// PrintAllNPCs prints all NPCs to the console.
// (For debugging purposes; consider using structured logging in production.)
func (s *NPCService) PrintAllNPC() {
	for _, npc := range s.npcs {
		fmt.Println(npc)
	}
}
