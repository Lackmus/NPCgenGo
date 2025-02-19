package controller

import (
	"fmt"

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
	factory             service.NPCFactory
	builder             service.NPCBuilder
	npc                 model.NPC
	observers           []shared.NPCEditObserver
}

// NewNPCEditController : Returns a new NPC controller.
func NewNPCEditController(view shared.NPCEditViewer, loader shared.NPCConfigLoader) (*NPCEditControllerImp, error) {
	creationDataService, err := service.NewCreationDataService(loader)
	if err != nil {
		return nil, fmt.Errorf("failed to create NPCEditController: %w", err)
	}
	creationOptions := service.NewNPCCreationOptions(*creationDataService)
	randomizerService := service.NewRandomizerService(*creationDataService, *creationOptions)

	return &NPCEditControllerImp{
		view:                view,
		creationDataService: *creationDataService,
		creationOptions:     *creationOptions,
		randomizerService:   *randomizerService,
		factory:             *service.NewNPCFactory(*randomizerService),
		observers:           []shared.NPCEditObserver{view},
	}, nil
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

// create a new NPC
func (c *NPCEditControllerImp) CreateNPC(npcType string, faction string) {
	c.builder = c.factory.CreateNPCWithOptions(npcType, faction)
}

func (c *NPCEditControllerImp) LoadNPC(npc model.NPC) {
	c.npc = npc
	c.builder = *service.NewBuilder().
		WithID(npc.ID()).
		WithName(npc.Name()).
		WithFaction(npc.Faction()).
		WithSpecies(npc.Species()).
		WithType(npc.NPCType()).
		WithSubType(npc.NPCSubtype()).
		WithTrait(npc.Trait()).
		WithStats(npc.Stats()).
		WithItems(npc.Items()).
		WithAbilities(npc.Abilities())

	c.NotifyObservers()
}

// Save NPC (return updated NPC)
func (c *NPCEditControllerImp) SaveNPC() model.NPC {
	c.npc = c.builder.Build(c.randomizerService)
	c.NotifyObservers()
	return c.npc
}
