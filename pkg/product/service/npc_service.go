package service

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"

	"slices"

	"github.com/lackmus/npcgengo/pkg/product/model"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

type NPCService struct {
	loader    shared.NPCStorage
	npcs      map[string]model.NPC
	idCounter int
	freeIDs   []int
}

// NewNPCService creates an NPCService and attempts to load existing NPCs from storage.
// It returns the service instance and any error encountered while loading existing NPCs.
// Partial loads are preserved: the returned error may be non-nil while the service
// contains any successfully loaded NPCs.
func NewNPCService(ctx context.Context, loader shared.NPCStorage) (*NPCService, error) {
	s := &NPCService{
		loader: loader,
		npcs:   make(map[string]model.NPC),
	}

	data, err := s.loader.LoadAllNPCs(ctx)
	if data != nil {
		s.npcs = data
	}
	if err != nil {
		// return service with partial data and the initialization error
		// caller may choose to log/fail depending on context
		s.initializeCounters()
		return s, err
	}

	s.initializeCounters()
	return s, nil
}

func (s *NPCService) initializeCounters() {
	s.idCounter = 0
	s.freeIDs = []int{}

	for idStr := range s.npcs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Non-numeric NPC ID %q found; skipping for counter purposes", idStr)
			continue
		}
		if id >= s.idCounter {
			s.idCounter = id + 1
		}
	}
}

func (s *NPCService) generateNewID() string {
	var id int
	if len(s.freeIDs) > 0 {
		id = s.freeIDs[0]
		s.freeIDs = s.freeIDs[1:]
	} else {
		id = s.idCounter
		s.idCounter++
	}
	return strconv.Itoa(id)
}

func (s *NPCService) AddNPC(npc model.NPC) {
	if npc.ID == "" {
		npc.ID = s.generateNewID()
	}

	s.npcs[npc.ID] = npc

	if err := s.loader.SaveNPC(context.Background(), npc); err != nil {
		log.Printf("Error saving NPC (ID %s): %v", npc.ID, err)
	}
}

func (s *NPCService) GetAllNPCs() []model.NPC {
	npcList := make([]model.NPC, 0, len(s.npcs))
	for _, npc := range s.npcs {
		npcList = append(npcList, npc)
	}
	return slices.Clone(npcList)
}

func (s *NPCService) GetNPCByID(id string) (model.NPC, error) {
	npc, found := s.npcs[id]
	if !found {
		return model.NPC{}, fmt.Errorf("NPC with ID %s not found", id)
	}
	return npc, nil
}

func (s *NPCService) DeleteNPC(id string) error {
	if _, found := s.npcs[id]; !found {
		return fmt.Errorf("NPC with ID %s not found", id)
	}
	delete(s.npcs, id)
	if num, err := strconv.Atoi(id); err == nil {
		s.freeIDs = append(s.freeIDs, num)
		sort.Ints(s.freeIDs)
	}
	return s.loader.DeleteNPC(context.Background(), id)
}

func (s *NPCService) DeleteAllNPCs() {
	if err := s.loader.DeleteAllNPCs(context.Background()); err != nil {
		log.Printf("Error deleting all NPCs: %v", err)
	}
	s.npcs = make(map[string]model.NPC)
	s.idCounter = 0
	s.freeIDs = []int{}
}

func (s *NPCService) CountNPCs() int {
	return len(s.npcs)
}

func (s *NPCService) PrintAllNPCs() {
	for _, npc := range s.npcs {
		log.Printf("%+v", npc)
	}
}
