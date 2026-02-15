// Description: This file contains the implementation of the NPCService struct and its methods.
package service

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"slices"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

type NPCService struct {
	loader    shared.NPCStorage
	npcs      map[string]model.NPC
	idCounter int
	freeIDs   []int
}

func NewNPCService(loader shared.NPCStorage) *NPCService {
	s := &NPCService{
		loader: loader,
		npcs:   make(map[string]model.NPC),
	}
	s.initNPCService()
	return s
}

func (s *NPCService) initNPCService() {
	var err error
	s.npcs, err = s.loader.LoadAllNPC()
	if err != nil {
		log.Printf("Error loading NPCs: %v", err)
		s.npcs = make(map[string]model.NPC)
	}

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

	if err := s.loader.SaveNPC(npc); err != nil {
		log.Printf("Error saving NPC (ID %s): %v", npc.ID, err)
	}
}

func (s *NPCService) GetAllNPC() []model.NPC {
	npcList := make([]model.NPC, 0, len(s.npcs))
	for _, npc := range s.npcs {
		npcList = append(npcList, npc)
	}
	return slices.Clone(npcList)
}

func (s *NPCService) GetNPCByLocation(locationID string) []model.NPC {
	var result []model.NPC
	for _, npc := range s.npcs {
		if npc.LocationID == locationID {
			result = append(result, npc)
		}
	}
	return result
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
	return s.loader.DeleteNPC(id)
}

func (s *NPCService) DeleteAllNPC() {
	if err := s.loader.DeleteAllNPC(); err != nil {
		log.Printf("Error deleting all NPCs: %v", err)
	}
	s.npcs = make(map[string]model.NPC)
	s.idCounter = 0
	s.freeIDs = []int{}
}

func (s *NPCService) CountNPC() int {
	return len(s.npcs)
}

func (s *NPCService) PrintAllNPC() {
	for _, npc := range s.npcs {
		fmt.Println(npc)
	}
}
