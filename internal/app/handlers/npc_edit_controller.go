// Package controller provides controllers for managing NPC editing.
package handlers

import (
	"errors"
	"log"

	"github.com/lackmus/npcgengo/pkg/product/model"
	cp "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
	"github.com/lackmus/npcgengo/pkg/product/service"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

type NPCEditController struct {
	creationSupplier  *service.NPCCreationSupplier
	rand              *service.RandomizerService
	npcBuilder        *service.NPCBuilder
	npcListController *NPCListController
	observers         []shared.NPCEditObserver
	errors            []error
}

func NewNPCEditController(creationSupplier *service.NPCCreationSupplier, controller *NPCListController, locationID string) *NPCEditController {
	return &NPCEditController{
		creationSupplier:  creationSupplier,
		rand:              creationSupplier.RandomizerService,
		npcBuilder:        service.NewNPCBuilder(creationSupplier, locationID),
		npcListController: controller,
		observers:         []shared.NPCEditObserver{},
		errors:            []error{},
	}
}

func (c *NPCEditController) EditNPC() model.NPC {
	return *c.npcBuilder.GetNPC()
}

func (c *NPCEditController) RegisterObserver(o shared.NPCEditObserver) {
	c.observers = append(c.observers, o)
}

func (c *NPCEditController) NotifyObservers() {
	for _, o := range c.observers {
		o.UpdateNPC(*c.npcBuilder.GetNPC())
	}
}

func (c *NPCEditController) NotifyObserversField(field cp.CompEnum, value any) {
	for _, o := range c.observers {
		o.UpdateField(field, value)
	}
}

func (c *NPCEditController) NotifyObserversError(err error) {
	for _, o := range c.observers {
		o.OnNPCEditError(err)
	}
	log.Printf("Error: %v", err)
}

func (c *NPCEditController) LoadNPC(npc model.NPC) {
	c.npcBuilder.WithNPC(npc)
	c.NotifyObservers()
}

func (c *NPCEditController) SaveNPC() {
	npc, err := c.npcBuilder.Build()
	if err != nil {
		c.handleErrors("SaveNPC", "NPC")
		return
	}
	c.NotifyObservers()
	c.npcListController.UpdateNpc(npc)
}

func (c *NPCEditController) SaveField(field cp.CompEnum, value string) error {
	switch field {
	case cp.CompName:
		c.npcBuilder.WithName(value)
	case cp.CompFaction:
		c.npcBuilder.WithFaction(value)
	case cp.CompSpecies:
		c.npcBuilder.WithSpecies(value)
	case cp.CompType:
		c.npcBuilder.WithType(value)
	case cp.CompSubtype:
		c.npcBuilder.WithSubtype(value)
	case cp.CompTrait:
		c.npcBuilder.WithTrait(value)
	case cp.CompStats:
		c.npcBuilder.WithSubtypeStats(value)
	case cp.CompItems:
		c.npcBuilder.WithSubtypeEquipment(value)
	case cp.CompDescription:
		c.npcBuilder.WithDescription(value)
	default:
		log.Printf("SaveField: unrecognized field %v (%s)", field, field.String())
		return nil
	}
	if err := c.handleErrors("SaveField", field.String()); err != nil {
		return err
	}

	return nil
}

func (c *NPCEditController) RandomizeField(field cp.CompEnum) string {
	var value string
	switch field {
	case cp.CompName:
		value = c.npcBuilder.WithRandomName().GetNPC().GetComponent(cp.CompName)
	case cp.CompFaction:
		value = c.npcBuilder.WithRandomFaction().GetNPC().GetComponent(cp.CompFaction)
	case cp.CompSpecies:
		value = c.npcBuilder.WithRandomSpecies().GetNPC().GetComponent(cp.CompSpecies)
	case cp.CompType:
		value = c.npcBuilder.WithRandomType().GetNPC().GetComponent(cp.CompType)
	case cp.CompSubtype:
		value = c.npcBuilder.WithRandomSubtype().GetNPC().GetComponent(cp.CompSubtype)
	case cp.CompTrait:
		value = c.npcBuilder.WithRandomTrait().GetNPC().GetComponent(cp.CompTrait)
	case cp.CompStats:
		value = c.npcBuilder.WithRandomSubtypeStats().GetNPC().GetComponent(cp.CompStats)
	case cp.CompItems:
		value = c.npcBuilder.WithRandomSubtypeEquipment().GetNPC().GetComponent(cp.CompItems)
	case cp.CompDescription:
		value = c.npcBuilder.WithRandomDescription().GetNPC().GetComponent(cp.CompDescription)
	default:
		log.Printf("RandomizeField: unrecognized field %v (%s)", field, field.String())
		return ""
	}

	c.handleErrors("RandomizeField", field.String())

	return value
}

func (c *NPCEditController) GetFieldOptions(field cp.CompEnum) []string {
	options := c.creationSupplier.CreationOptions
	switch field {
	case cp.CompType:
		return options.NpcTypes
	case cp.CompFaction:
		return options.Factions
	case cp.CompSpecies:
		return options.Species
	case cp.CompSubtype:
		npcType := c.npcBuilder.GetNPCType()
		if subtypes, ok := options.NpcSubtypeForTypeMap[npcType]; ok {
			return subtypes
		}
		return nil
	case cp.CompTrait:
		return options.Traits
	default:
		return nil
	}
}

func (c *NPCEditController) handleErrors(context string, field string) error {
	if c.npcBuilder.HasErrors() {
		err := errors.Join(c.npcBuilder.Errors()...)
		log.Printf("%s: error saving field %v: %v", context, field, err)
		c.NotifyObserversError(err)
		// Optionally, you can clear builder errors here if desired:
		// c.npcBuilder.ClearErrors()
		return err
	}
	return nil
}
