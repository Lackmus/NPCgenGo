package controller

import (
	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/shared"
)

// NPCEditController : The controller for the NPC model.
type NPCEditController struct {
	creationSupplier *service.NPCCreationSupplier
	rand             *service.RandomizerService
	npcBuilder       *service.NPCBuilder
	npc              model.NPC
	observers        []shared.NPCEditObserver
}

// NewNPCEditController : Returns a new NPC controller.
func NewNPCEditController(view shared.NPCEditViewer, creationSupplier *service.NPCCreationSupplier) *NPCEditController {
	return &NPCEditController{
		creationSupplier: creationSupplier,
		rand:             creationSupplier.RandomizerService,
		observers:        []shared.NPCEditObserver{view},
	}
}

// EditNPC : Edit NPC (return updated NPC)
func (c *NPCEditController) EditNPC(npc model.NPC) {
}

func (c *NPCEditController) RegisterObserver(o shared.NPCEditObserver) {
	c.observers = append(c.observers, o)
}

func (c *NPCEditController) NotifyObservers() {
	for _, o := range c.observers {
		o.UpdateNPC(c.npc)
	}
}

func (c *NPCEditController) NotifyObserversField(field string, value any) {
	for _, o := range c.observers {
		o.UpdateField(field, value)
	}
}

func (c *NPCEditController) LoadNPC(npc model.NPC) {
	c.npc = npc
	c.npcBuilder = service.NewNPCBuilderFromNPC(npc)
	c.NotifyObservers()
}

// Save NPC (return updated NPC)
func (c *NPCEditController) SaveNPC() model.NPC {
	c.npc = c.npcBuilder.BuildWithRandom(c.rand)
	c.NotifyObservers()
	return c.npc
}

func (c *NPCEditController) RandomizeField(field string) {
	var updatedValue any
	switch field {
	case "name":
		updatedValue = c.rand.GenerateName(c.npcBuilder.Species)
	case "faction":
		updatedValue = c.rand.RandomFaction()
	case "species":
		updatedValue = c.rand.RandomSpecies()
	case "npcType":
		updatedValue = c.rand.RandomType()
	case "npcSubtype":
		updatedValue = c.rand.RandomSubtype(c.npcBuilder.NPCType)
	case "trait":
		updatedValue = c.rand.RandomTrait()
	case "drive":
		//option = c.randomizerService.GenerateDrive()
		//c.Builder.WithDrive(option.(string))
	case "stats":
		updatedValue = c.rand.ApplySubtypeStats(c.npcBuilder.NPCSubtype)
	case "items":
		updatedValue = c.rand.GenerateEquipment(c.npcBuilder.NPCSubtype)

	case "abilities":
		//option = c.randomizerService.GenerateAbilities(c.Builder.NPCSubtype)
		//c.Builder.WithAbilities(option.(map[string]string))
	default:
		return
	}
	c.SaveField(field, updatedValue)
	c.NotifyObserversField(field, updatedValue)
}

func (c *NPCEditController) SaveField(field string, value any) {
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

func (c *NPCEditController) GetFieldOptions(field string) []string {
	var options = c.creationSupplier.CreationOptions
	switch field {
	case "npcType":
		return options.NpcTypes
	case "faction":
		return options.Factions
	case "species":
		return options.Species
	case "npcSubtype":
		return options.NpcSubtypeForTypeMap[c.npcBuilder.NPCType]
	case "trait":
		return options.Traits
	default:
		return nil
	}
}
