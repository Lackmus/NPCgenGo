// Description: This file contains the implementation of the NPCService struct and its methods.
package service

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

// NPCService manages NPCs and loads them from a storage backend.
// Note: If concurrent access is expected, consider adding a sync.RWMutex.
type NPCService struct {
	loader shared.NPCStorage
	npcs   map[string]model.NPC // Using a map for faster lookup instead of a slice.
	// idCounter is the next new ID to use if no free IDs are available.
	idCounter int
	// freeIDs holds IDs (as integers) that became available when NPCs were deleted.
	freeIDs []int
}

// NewNPCService creates a new NPCService with the provided NPC storage.
// It initializes the service with the NPCs from the storage.
func NewNPCService(loader shared.NPCStorage) *NPCService {
	s := &NPCService{
		loader: loader,
		npcs:   make(map[string]model.NPC),
	}
	s.initNPCService(loader)
	return s
}

// initNPCService initializes the NPCService with NPCs from the storage.
// It also scans the loaded NPCs to set the idCounter to one greater than the highest existing ID.
func (s *NPCService) initNPCService(loader shared.NPCStorage) {
	var err error
	s.npcs, err = loader.LoadAllNPC()
	if err != nil {
		log.Printf("Error loading NPCs: %v", err)
		s.npcs = make(map[string]model.NPC)
	}

	s.idCounter = 0
	s.freeIDs = []int{}

	// Iterate over the loaded NPCs to determine the highest used numeric ID.
	for idStr := range s.npcs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			// If an ID is non-numeric, skip it (or handle it differently).
			log.Printf("Non-numeric NPC ID %q found; skipping for counter purposes", idStr)
			continue
		}
		// Update idCounter if the current ID is greater.
		if id >= s.idCounter {
			s.idCounter = id + 1
		}
	}
}

// generateNewID returns a new unique ID as a string.
// It reuses IDs from freeIDs if available; otherwise, it uses idCounter.
func (s *NPCService) generateNewID() string {
	var id int
	if len(s.freeIDs) > 0 {
		// Use the smallest available ID from freeIDs.
		id = s.freeIDs[0]
		// Remove the used ID from freeIDs.
		s.freeIDs = s.freeIDs[1:]
	} else {
		id = s.idCounter
		s.idCounter++
	}
	return strconv.Itoa(id)
}

// AddNPC adds a new NPC to the service and saves it.
// If the NPC does not already have an ID, it generates one.
func (s *NPCService) AddNPC(npc model.NPC) {
	if npc.ID == "" {
		npc.ID = s.generateNewID() // Generate a new ID if missing
	}

	// Use the updated NPC ID as the key
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

// DeleteNPC removes an NPC from the map and deletes it from the storage.
// It reclaims the NPC's ID for reuse.
func (s *NPCService) DeleteNPC(id string) error {
	if _, found := s.npcs[id]; !found {
		return fmt.Errorf("NPC with ID %s not found", id)
	}
	delete(s.npcs, id)
	// Try to parse the ID as an int and add it to freeIDs for reuse.
	if num, err := strconv.Atoi(id); err == nil {
		s.freeIDs = append(s.freeIDs, num)
		// Optional: sort freeIDs so that the smallest available ID is used first.
		sort.Ints(s.freeIDs)
	}
	return s.loader.DeleteNPC(id)
}

// DeleteAllNPC deletes all NPCs from the storage and clears the map.
// It also resets the ID counter and freeIDs.
func (s *NPCService) DeleteAllNPC() {
	if err := s.loader.DeleteAllNPC(); err != nil {
		log.Printf("Error deleting all NPCs: %v", err)
	}
	s.npcs = make(map[string]model.NPC)
	s.idCounter = 0
	s.freeIDs = []int{}
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
