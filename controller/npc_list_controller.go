// Description: This file contains the controller for the list of NPCs.
package controller

import (
	"log"

	h "github.com/lackmus/npcgengo/helper"
	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/shared"
)

// NPCListController manages the list of NPCs.
type NPCListController struct {
	npcService       *service.NPCService
	creationSupplier *service.NPCCreationSupplier
	observers        []shared.NPCObserver
	LocationID       string
}

func (c *NPCListController) GetNPCs() []model.NPC {
	return c.npcService.GetAllNPC()
}

// NewNPCListController creates a new NPCListController.
func NewNPCListController(creationSupplier *service.NPCCreationSupplier, npcService *service.NPCService, locationID string) *NPCListController {
	log.Println("Creating NPCListController...")
	return &NPCListController{
		npcService:       npcService,
		creationSupplier: creationSupplier,
		LocationID:       locationID,
		observers:        []shared.NPCObserver{},
	}
}

// creaTE random NPC
func (c *NPCListController) CreateRandomNPC() {
	npc, err := service.CreateNPCWithOptions(h.Random, h.Random, c.creationSupplier, c.LocationID)
	if err != nil {
		log.Printf("Error creating NPC: %v", err)
		return
	}
	c.AddNpc(npc)
}

// InitEditController initializes the NPC edit controller.
// It returns a new NPCEditController.
func (c *NPCListController) InitEditController() *NPCEditController {
	log.Println("Initializing edit controller...")
	return NewNPCEditController(c.creationSupplier, c, c.LocationID)
}

// InitView notifies observers to initialize the view.
func (c *NPCListController) InitView(view shared.NPCListViewer) {
	log.Println("Initializing view...")
	c.RegisterObserver(view)
	c.NotifyObservers()
}

// UpdateNpc updates an NPC in the service and notifies observers.
func (c *NPCListController) UpdateNpc(npc model.NPC) {
	c.AddNpc(npc)
	c.NotifyObservers()
}

// RegisterObserver adds an observer to the list.
func (c *NPCListController) RegisterObserver(o shared.NPCObserver) {
	c.observers = append(c.observers, o)
}

// NotifyObservers notifies all observers with the current list of NPCs.
func (c *NPCListController) NotifyObservers() {
	npcs := c.npcService.GetAllNPC()
	for _, o := range c.observers {
		o.Update(npcs)
	}
}

// GetAllNpcs returns all NPCs from the service.
func (c *NPCListController) GetAllNpcs() []model.NPC {
	return c.npcService.GetAllNPC()
}

// getNpcByID returns an NPC by id from the service.
func (c *NPCListController) GetNpcByID(id string) (model.NPC, error) {
	return c.npcService.GetNPCByID(id)
}

// AddNpc adds a new NPC to the service and notifies observers.
func (c *NPCListController) AddNpc(npc model.NPC) {
	c.npcService.AddNPC(npc)
	c.NotifyObservers()
}

// DeleteNPC deletes an NPC by id and notifies observers.
func (c *NPCListController) DeleteNPC(id string) {
	if err := c.npcService.DeleteNPC(id); err != nil {
		log.Printf("Error deleting NPC: %v", err)
	}
	c.NotifyObservers()
}

// DeleteAllNPC deletes all NPCs and notifies observers.
func (c *NPCListController) DeleteAllNPC() {
	c.npcService.DeleteAllNPC()
	c.NotifyObservers()
}
