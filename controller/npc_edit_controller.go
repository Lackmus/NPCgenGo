package controller

import (
	//"log"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/shared"
)

// Field name constants for consistency.
const (
	FieldName       = "name"
	FieldFaction    = "faction"
	FieldSpecies    = "species"
	FieldNPCType    = "npcType"
	FieldNPCSubtype = "npcSubtype"
	FieldTrait      = "trait"
	FieldDrive      = "drive"
	FieldStats      = "stats"
	FieldItems      = "items"
	FieldAbilities  = "abilities"
)

// NPCEditController manages the editing of an NPC.
type NPCEditController struct {
	creationSupplier *service.NPCCreationSupplier
	rand             *service.RandomizerService
	npcBuilder       *service.NPCBuilder
	npc              model.NPC
	observers        []shared.NPCEditObserver
}

func NewNPCEditController(view shared.NPCEditViewer, creationSupplier *service.NPCCreationSupplier) *NPCEditController {
	return &NPCEditController{
		creationSupplier: creationSupplier,
		rand:             creationSupplier.RandomizerService,
		observers:        []shared.NPCEditObserver{view},
	}
}

func (c *NPCEditController) EditNPC() model.NPC {
	return c.npc
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
	c.NotifyObservers()
}

func (c *NPCEditController) SaveNPC() model.NPC {
	c.NotifyObservers()
	return c.npc
}

/*
// RandomizeField randomizes a specific field based on its type and saves the change.
func (c *NPCEditController) RandomizeField(field string) {

	c.SaveField(field, updatedValue)
	c.NotifyObserversField(field, updatedValue)
}

// SaveField updates a single field in the NPC builder using a type-safe assertion.
func (c *NPCEditController) SaveField(field string, value any) {
	switch field {
	case FieldName:
		if v, ok := value.(string); ok {
			c.npcBuilder.WithName(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected string", FieldName)
		}
	case FieldFaction:
		if v, ok := value.(string); ok {
			c.npcBuilder.WithFaction(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected string", FieldFaction)
		}
	case FieldSpecies:
		if v, ok := value.(string); ok {
			c.npcBuilder.WithSpecies(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected string", FieldSpecies)
		}
	case FieldNPCType:
		if v, ok := value.(string); ok {
			c.npcBuilder.WithType(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected string", FieldNPCType)
		}
	case FieldNPCSubtype:
		if v, ok := value.(string); ok {
			c.npcBuilder.WithSubType(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected string", FieldNPCSubtype)
		}
	case FieldTrait:
		if v, ok := value.(string); ok {
			c.npcBuilder.WithTrait(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected string", FieldTrait)
		}
	case FieldDrive:
		if v, ok := value.(string); ok {
			c.npcBuilder.WithDrive(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected string", FieldDrive)
		}
	case FieldStats:
		if v, ok := value.(map[string]int); ok {
			c.npcBuilder.WithStats(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected map[string]int", FieldStats)
		}
	case FieldItems:
		if v, ok := value.(map[string]string); ok {
			c.npcBuilder.WithItems(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected map[string]string", FieldItems)
		}
	case FieldAbilities:
		if v, ok := value.(map[string]string); ok {
			c.npcBuilder.WithAbilities(v)
		} else {
			log.Printf("Type assertion failed for field %s; expected map[string]string", FieldAbilities)
		}
	default:
		// Unrecognized field; do nothing.
		return
	}
}

// GetFieldOptions returns available options for a given field using the creation supplier.
func (c *NPCEditController) GetFieldOptions(field string) []string {
	options := c.creationSupplier.CreationOptions
	switch field {
	case FieldNPCType:
		return options.NpcTypes
	case FieldFaction:
		return options.Factions
	case FieldSpecies:
		return options.Species
	case FieldNPCSubtype:
		return options.NpcSubtypeForTypeMap[c.npcBuilder.NPCType]
	case FieldTrait:
		return options.Traits
	default:
		return nil
	}
}
*/
