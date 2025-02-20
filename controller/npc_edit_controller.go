package controller

import (
	"fmt"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/shared"
)

// NPCController : The controller for the NPC model.

type NPCEditControllerImp struct {
	creationDataService service.CreationDataService
	creationOptions     service.NPCCreationOptions
	randomizerService   service.RandomizerService
	factory             service.NPCFactory
	npcBuilder          service.NPCBuilder
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

func (c *NPCEditControllerImp) NotifyObserversField(field string, value any) {
	for _, o := range c.observers {
		o.UpdateField(field, value)
	}
}

// create a new NPC
func (c *NPCEditControllerImp) CreateNPC(npcType string, faction string) {
	c.npcBuilder = c.factory.CreateNPCWithOptions(npcType, faction)
}

func (c *NPCEditControllerImp) LoadNPC(npc model.NPC) {
	c.npc = npc
	c.npcBuilder = *service.NewNPCBuilder().SetNPC(npc)
	c.NotifyObservers()
}

// Save NPC (return updated NPC)
func (c *NPCEditControllerImp) SaveNPC() model.NPC {
	c.npc = c.npcBuilder.BuildWithRandom(c.randomizerService)
	c.NotifyObservers()
	return c.npc
}

func (c *NPCEditControllerImp) RandomizeField(field string) {
	var updatedValue any
	switch field {
	case "name":
		updatedValue = c.randomizerService.GenerateName(c.npcBuilder.Species)
	case "faction":
		updatedValue = c.randomizerService.RandomFaction()
	case "species":
		updatedValue = c.randomizerService.RandomSpecies()
	case "npcType":
		updatedValue = c.randomizerService.RandomType()
	case "npcSubtype":
		updatedValue = c.randomizerService.RandomSubtype(c.npcBuilder.NPCType)
	case "trait":
		updatedValue = c.randomizerService.RandomTrait()
	case "drive":
		//option = c.randomizerService.GenerateDrive()
		//c.Builder.WithDrive(option.(string))
	case "stats":
		updatedValue = c.randomizerService.ApplySubtypeStats(c.npcBuilder.NPCSubtype)
	case "items":
		updatedValue = c.randomizerService.GenerateEquipment(c.npcBuilder.NPCSubtype)

	case "abilities":
		//option = c.randomizerService.GenerateAbilities(c.Builder.NPCSubtype)
		//c.Builder.WithAbilities(option.(map[string]string))
	default:
		return
	}
	c.SaveField(field, updatedValue)
	c.NotifyObserversField(field, updatedValue)
}

func (c *NPCEditControllerImp) SaveField(field string, value any) {
	switch field {
	case "name":
		c.npcBuilder.WithName(value.(string))
	case "faction":
		c.npcBuilder.WithFaction(value.(string))
	case "species":
		c.npcBuilder.WithSpecies(value.(string))
	case "npcType":
		c.npcBuilder.WithType(value.(string))
	case "npcSubtype":
		c.npcBuilder.WithSubType(value.(string))
	case "trait":
		c.npcBuilder.WithTrait(value.(string))
	case "drive":
		c.npcBuilder.WithDrive(value.(string))
	case "stats":
		c.npcBuilder.WithStats(value.(map[string]int))
	case "items":
		c.npcBuilder.WithItems(value.(map[string]string))
	case "abilities":
		c.npcBuilder.WithAbilities(value.(map[string]string))
	default:
		return
	}
}

func (c *NPCEditControllerImp) GetFieldOptions(field string) []string {
	switch field {
	case "npcType":
		return c.creationOptions.NpcTypes
	case "faction":
		return c.creationOptions.Factions
	case "species":
		return c.creationOptions.Species
	case "npcSubtype":
		return c.creationOptions.NpcSubtypeForTypeMap[c.npcBuilder.NPCType]
	case "trait":
		return c.creationOptions.Traits
	default:
		return nil
	}
}
