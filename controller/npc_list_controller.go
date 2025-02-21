package controller

import (
	"log"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/shared"
)

// NPCListController manages the list of NPCs.
type NPCListController struct {
	npcService       *service.NPCService
	creationSupplier *service.NPCCreationSupplier
	observers        []shared.NPCObserver
}

// NewNPCListController creates a new NPCListController.
func NewNPCListController(storage shared.NPCStorage, creationSupplier *service.NPCCreationSupplier, view shared.NPCViewer) *NPCListController {
	npcService := service.NewNPCService(storage)
	return &NPCListController{
		npcService:       npcService,
		creationSupplier: creationSupplier,
		observers:        []shared.NPCObserver{view},
	}
}

// InitEditController initializes the NPC edit controller.
func (c *NPCListController) InitEditController(editView shared.NPCEditViewer) *NPCEditController {
	return NewNPCEditController(editView, c.creationSupplier)
}

// InitView notifies observers to initialize the view.
func (c *NPCListController) InitView() {
	c.NotifyObservers()
}

// UpdateNpc updates an NPC in the service and notifies observers.
func (c *NPCListController) UpdateNpc(npc model.NPC) {
	if err := c.npcService.UpdateNPC(npc); err != nil {
		log.Printf("Error updating NPC: %v", err)
	}
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
