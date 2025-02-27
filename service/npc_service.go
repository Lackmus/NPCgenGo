// Description: This file contains the implementation of the NPCService struct and its methods.
// The NPCService struct manages NPCs and loads them from a storage backend.
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

// NewNPCService creates a new NPCService with the provided NPC storage.
// It initializes the service with the NPCs from the storage.
func NewNPCService(loader shared.NPCStorage) *NPCService {
	n := &NPCService{loader: loader}
	n.initNPCService(loader)
	return n
}

// initNPCService initializes the NPCService with NPCs from the storage.
// It is called during the creation of the NPCService.
func (s *NPCService) initNPCService(loader shared.NPCStorage) {
	var err error
	s.npcs, err = loader.LoadAllNPC()
	if err != nil {
		log.Printf("Error loading NPCs: %v", err)
		s.npcs = make(map[string]model.NPC)
	}
}

// AddNPC adds a new NPC to the service and saves it.
// It is safe for concurrent use.
func (s *NPCService) AddNPC(npc model.NPC) {
	s.npcs[npc.ID] = npc
	if err := s.loader.SaveNPC(npc); err != nil {
		log.Printf("Error saving NPC (ID %s): %v", npc.ID, err)
	}
}

// GetAllNPCs returns all NPCs as a slice.
// It creates a copy of the NPCs to avoid concurrent map access.
// Consider using a sync.RWMutex if concurrent access is expected.
func (s *NPCService) GetAllNPC() []model.NPC {
	npcList := make([]model.NPC, 0, len(s.npcs))
	for _, npc := range s.npcs {
		npcList = append(npcList, npc)
	}
	return append([]model.NPC(nil), npcList...)
}

// GetNPCByID returns the NPC with the specified ID.
// It returns an error if the NPC is not found.
func (s *NPCService) GetNPCByID(id string) (model.NPC, error) {
	npc, found := s.npcs[id]
	// use log.Printf for logging errors
	if !found {
		return model.NPC{}, fmt.Errorf("NPC with ID %s not found", id)
	}
	return npc, nil
}

// UpdateNPC updates an existing NPC and saves it.
// It returns an error if the NPC with the specified ID is not found.
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
// It returns an error if the NPC with the specified ID is not found.
func (s *NPCService) DeleteNPC(id string) error {
	if _, found := s.npcs[id]; !found {
		return fmt.Errorf("NPC with ID %s not found", id)
	}
	delete(s.npcs, id)
	return s.loader.DeleteNPC(id)
}

// DeleteAllNPC deletes all NPCs from the storage and clears the map.
// It returns an error if the deletion fails.
func (s *NPCService) DeleteAllNPC() {
	if err := s.loader.DeleteAllNPC(); err != nil {
		log.Printf("Error deleting all NPCs: %v", err)
	}
	s.npcs = make(map[string]model.NPC)
}

// CountNPC returns the number of NPCs.
// It is useful for testing and debugging.
func (s *NPCService) CountNPC() int {
	return len(s.npcs)
}

// PrintAllNPCs prints all NPCs to the console.
// Note: This method is not recommended for large datasets.
// (For debugging purposes; consider using structured logging in production.)
func (s *NPCService) PrintAllNPC() {
	for _, npc := range s.npcs {
		fmt.Println(npc)
	}
}
