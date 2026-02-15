package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

type NPCGroupService struct {
	loader     shared.NPCGroupStorage
	npcService *NPCService
	groups     map[string]model.NPCGroup
	idCounter  int
	freeIDs    []int
}

func NewNPCGroupService(loader shared.NPCGroupStorage, npcService *NPCService) *NPCGroupService {
	s := &NPCGroupService{
		loader:     loader,
		npcService: npcService,
		groups:     make(map[string]model.NPCGroup),
	}
	s.initNPCGroupService()
	return s
}

func (s *NPCGroupService) initNPCGroupService() {
	var err error
	s.groups, err = s.loader.LoadAllNPCGroups()
	if err != nil {
		log.Printf("Error loading NPC groups: %v", err)
		s.groups = make(map[string]model.NPCGroup)
	}

	s.idCounter = 0
	s.freeIDs = []int{}

	for idStr := range s.groups {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Non-numeric NPC group ID %q found; skipping for counter purposes", idStr)
			continue
		}
		if id >= s.idCounter {
			s.idCounter = id + 1
		}
	}
}

func (s *NPCGroupService) generateNewID() string {
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

func (s *NPCGroupService) AddNPCGroup(group model.NPCGroup) {
	if group.Name == "" {
		group.Name = s.generateNewID()
	}

	s.groups[group.Name] = group

	if err := s.loader.SaveNPCGroup(group); err != nil {
		log.Printf("Error saving NPC group (ID %s): %v", group.Name, err)
	}
}

func (s *NPCGroupService) GetNPCGroupsByLocation(locationID string) []model.NPCGroup {
	var result []model.NPCGroup
	for _, group := range s.groups {
		if group.LocationID == locationID {
			result = append(result, group)
		}
	}
	return result
}

func (s *NPCGroupService) GetNPCGroupByID(id string) (model.NPCGroup, error) {
	group, found := s.groups[id]
	if !found {
		return model.NPCGroup{}, fmt.Errorf("NPC group with ID %s not found", id)
	}
	return group, nil
}

func (s *NPCGroupService) GetNPCsInGroup(groupID string) ([]model.NPC, error) {
	group, err := s.GetNPCGroupByID(groupID)
	if err != nil {
		return nil, err
	}

	var npcs []model.NPC
	for _, npcID := range group.NPCIDs {
		npc, err := s.npcService.GetNPCByID(npcID)
		if err != nil {
			log.Printf("Warning: NPC %s in group %s not found", npcID, groupID)
			continue
		}
		npcs = append(npcs, npc)
	}

	return npcs, nil
}

func (s *NPCGroupService) CreateGroup(name, locationID string, npcIDs []string) model.NPCGroup {
	group := *model.NewNPCGroup(name)
	group.LocationID = locationID

	for _, npcID := range npcIDs {
		if _, err := s.npcService.GetNPCByID(npcID); err == nil {
			group.AddNPC(npcID)
		} else {
			log.Printf("Warning: NPC %s not found, not adding to group", npcID)
		}
	}

	s.AddNPCGroup(group)
	return group
}
