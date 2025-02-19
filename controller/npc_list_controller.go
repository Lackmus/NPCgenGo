package controller

import (
	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/shared"
)

// =============================================================================
// DefaultNpcController
// =============================================================================

// DefaultNpcController is a controller that mediates between the NPC service and the UI (observers).
type NPCListController struct {
	npcService service.NPCService
	loader     shared.NPCConfigLoader
	view       shared.NPCViewer
	observers  []shared.NPCObserver
}

// NewDefaultNpcController creates a new instance of DefaultNpcController.
func NewNPCListController(storage shared.NPCStorage, loader shared.NPCConfigLoader, view shared.NPCViewer) *NPCListController {
	return &NPCListController{
		npcService: service.NewNPCService(storage),
		loader:     loader,
		view:       view,
		observers:  []shared.NPCObserver{view},
	}
}

func (c *NPCListController) InitEditController(editView shared.NPCEditViewer) shared.NPCEditController {
	return NewNPCEditController(editView, c.loader)
}

// intview
func (c *NPCListController) InitView() {
	c.NotifyObservers()
}

// UpdateNpc updates an NPC and notifies observers
func (c *NPCListController) UpdateNpc(npc model.NPC) {

	c.npcService.UpdateNPC(npc)
	c.NotifyObservers()
}

// RegisterObserver registers a view as an observer
func (c *NPCListController) RegisterObserver(o shared.NPCObserver) {
	c.observers = append(c.observers, o)
}

// NotifyObservers notifies all registered observers (views)
func (c *NPCListController) NotifyObservers() {
	npcs := c.npcService.GetAllNpcs()
	for _, o := range c.observers {
		o.Update(npcs)
	}
}

// GetAllNPCs fetches NPCs from the service
func (c *NPCListController) GetAllNpcs() []model.NPC {
	return c.npcService.GetAllNpcs()
}

// AddNpc adds an NPC and notifies observers
func (c *NPCListController) AddNpc(npc model.NPC) {
	c.npcService.AddNpc(npc)
	c.NotifyObservers() // Now updates all views
}
