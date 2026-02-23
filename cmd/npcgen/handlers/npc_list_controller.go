// Description: This file contains the controller for the list of NPCs.
package handlers

import (
	"log"

	h "github.com/lackmus/npcgengo/internal/platform/helpers"
	"github.com/lackmus/npcgengo/pkg/product/model"
	"github.com/lackmus/npcgengo/pkg/product/service"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

type NPCListController struct {
	NpcService       *service.NPCService
	CreationSupplier *service.NPCCreationSupplier
	observers        []shared.NPCObserver
	LocationID       string
}

func NewNPCListController(creationSupplier *service.NPCCreationSupplier, npcService *service.NPCService, locationID string) *NPCListController {
	log.Println("Creating NPCListController...")
	return &NPCListController{
		NpcService:       npcService,
		CreationSupplier: creationSupplier,
		LocationID:       locationID,
		observers:        []shared.NPCObserver{},
	}
}

func (c *NPCListController) CreateRandomNPC() (model.NPC, error) {
	npc, err := service.CreateNPCWithOptions(h.Random, h.Random, c.CreationSupplier, c.LocationID)
	if err != nil {
		log.Printf("Error creating NPC: %v", err)
		return model.NPC{}, err
	}
	c.AddNpc(npc)
	return npc, nil
}

func (c *NPCListController) CreateNPC(npctype string, faction string) (model.NPC, error) {
	npc, err := service.CreateNPCWithOptions(npctype, faction, c.CreationSupplier, c.LocationID)
	if err != nil {
		log.Printf("Error creating NPC: %v", err)
		return model.NPC{}, err
	}
	c.AddNpc(npc)
	return npc, nil
}

func (c *NPCListController) InitEditController() *NPCEditController {
	log.Println("Initializing edit controller...")
	return NewNPCEditController(c.CreationSupplier, c, c.LocationID)
}

func (c *NPCListController) InitView(view shared.NPCListViewer) {
	log.Println("Initializing view...")
	c.RegisterObserver(view)
	c.NotifyObservers()
}

func (c *NPCListController) UpdateNpc(npc model.NPC) {
	c.AddNpc(npc)
}

func (c *NPCListController) RegisterObserver(o shared.NPCObserver) {
	c.observers = append(c.observers, o)
}

func (c *NPCListController) NotifyObservers() {
	npcs := c.NpcService.GetNPCByLocation(c.LocationID)
	for _, o := range c.observers {
		o.Update(npcs)
	}
}

func (c *NPCListController) GetAllNpcs() []model.NPC {
	npcs := c.NpcService.GetAllNPC()
	if len(npcs) == 0 {
		log.Println("No NPCs found in the current location.")
	}
	return npcs
}

func (c *NPCListController) GetNpcByID(id string) (model.NPC, error) {
	npc, err := c.NpcService.GetNPCByID(id)
	if err != nil {
		return model.NPC{}, err
	}
	return npc, nil
}

func (c *NPCListController) AddNpc(npc model.NPC) {
	c.NpcService.AddNPC(npc)
	c.NotifyObservers()
}

func (c *NPCListController) DeleteNPC(id string) {
	if err := c.NpcService.DeleteNPC(id); err != nil {
		log.Printf("Error deleting NPC: %v", err)
	}
	c.NotifyObservers()
}

func (c *NPCListController) DeleteAllNPC() {
	c.NpcService.DeleteAllNPC()
	c.NotifyObservers()
}

func (c *NPCListController) GetNPCByLocation() []model.NPC {
	return c.NpcService.GetNPCByLocation(c.LocationID)
}

func (c *NPCListController) CreateNPCGroup() {
	// Implementation will be in the view that uses this controller
}
