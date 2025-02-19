package controller

import (
	"log"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/shared"
)

// NPCController : The controller for the NPC model.

type NPCEditControllerImp struct {
	view                shared.NPCEditViewer
	creationDataService service.CreationDataService
	creationOptions     service.NPCCreationOptions
	randomizerService   service.RandomizerService
	npc                 model.NPC
	observers           []shared.NPCEditObserver
}

// NewNPCEditController : Returns a new NPC controller.
func NewNPCEditController(view shared.NPCEditViewer, loader shared.NPCConfigLoader) *NPCEditControllerImp {
	creationDataService, err := service.NewCreationDataService(loader)
	if err != nil {
		log.Fatalf("Error creating NPCEditController %v", err)
	}
	creationOptions := service.NewNPCCreationOptions(*creationDataService)
	randomizerService := service.NewRandomizerService(*creationDataService, *creationOptions)
	return &NPCEditControllerImp{
		view:                view,
		creationDataService: *creationDataService,
		creationOptions:     *creationOptions,
		randomizerService:   *randomizerService,
		observers:           []shared.NPCEditObserver{view},
	}
}

func (c *NPCEditControllerImp) CreateNPC(npcType string, faction string) {
	factory := service.NewNPCFactory(c.randomizerService)
	c.npc = factory.CreateNPCWithOptions(npcType, faction)

}

// EditNPC : Edit NPC (return updated NPC)
func (c *NPCEditControllerImp) EditNPC(npc model.NPC) {
}

func (c *NPCEditControllerImp) RegisterObserver(o shared.NPCEditObserver) {
	c.observers = append(c.observers, o)
}

func (c *NPCEditControllerImp) NotifyObservers() {
	for _, o := range c.observers {
		o.UpdateNPC(c.npc)
	}
}

// Save NPC (return updated NPC)
func (c *NPCEditControllerImp) SaveNPC() model.NPC {
	c.NotifyObservers()
	return c.npc
}
