// Package controller provides controllers for managing NPC editing.
package controllers

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
	saveHandlers      map[cp.CompEnum]func(string)
	randomHandlers    map[cp.CompEnum]func() string
	optionHandlers    map[cp.CompEnum]func() []string
}

func NewNPCEditController(creationSupplier *service.NPCCreationSupplier, controller *NPCListController, locationID string) *NPCEditController {
	c := &NPCEditController{
		creationSupplier:  creationSupplier,
		rand:              creationSupplier.RandomizerService,
		npcBuilder:        service.NewNPCBuilder(creationSupplier, locationID),
		npcListController: controller,
		observers:         []shared.NPCEditObserver{},
		errors:            []error{},
	}
	c.initHandlers()
	return c
}

func (c *NPCEditController) initHandlers() {
	c.saveHandlers = map[cp.CompEnum]func(string){
		cp.CompName:        func(value string) { c.npcBuilder.WithName(value) },
		cp.CompFaction:     func(value string) { c.npcBuilder.WithFaction(value) },
		cp.CompSpecies:     func(value string) { c.npcBuilder.WithSpecies(value) },
		cp.CompType:        func(value string) { c.npcBuilder.WithType(value) },
		cp.CompSubtype:     func(value string) { c.npcBuilder.WithSubtype(value) },
		cp.CompTrait:       func(value string) { c.npcBuilder.WithTrait(value) },
		cp.CompStats:       func(value string) { c.npcBuilder.WithSubtypeStats(value) },
		cp.CompItems:       func(value string) { c.npcBuilder.WithSubtypeEquipment(value) },
		cp.CompDescription: func(value string) { c.npcBuilder.WithDescription(value) },
	}

	c.randomHandlers = map[cp.CompEnum]func() string{
		cp.CompName:        func() string { return c.npcBuilder.WithRandomName().GetNPC().Name() },
		cp.CompFaction:     func() string { return c.npcBuilder.WithRandomFaction().GetNPC().Faction() },
		cp.CompSpecies:     func() string { return c.npcBuilder.WithRandomSpecies().GetNPC().Species() },
		cp.CompType:        func() string { return c.npcBuilder.WithRandomType().GetNPC().Type() },
		cp.CompSubtype:     func() string { return c.npcBuilder.WithRandomSubtype().GetNPC().Subtype() },
		cp.CompTrait:       func() string { return c.npcBuilder.WithRandomTrait().GetNPC().Trait() },
		cp.CompStats:       func() string { return c.npcBuilder.WithRandomSubtypeStats().GetNPC().Stats() },
		cp.CompItems:       func() string { return c.npcBuilder.WithRandomSubtypeEquipment().GetNPC().Items() },
		cp.CompDescription: func() string { return c.npcBuilder.WithRandomDescription().GetNPC().Description() },
	}

	c.optionHandlers = map[cp.CompEnum]func() []string{
		cp.CompType:    func() []string { return c.creationSupplier.CreationOptions.NpcTypes },
		cp.CompFaction: func() []string { return c.creationSupplier.CreationOptions.Factions },
		cp.CompSpecies: func() []string { return c.creationSupplier.CreationOptions.Species },
		cp.CompTrait:   func() []string { return c.creationSupplier.CreationOptions.Traits },
		cp.CompSubtype: func() []string {
			npcType := c.npcBuilder.GetNPC().Type()
			if subtypes, ok := c.creationSupplier.CreationOptions.NpcSubtypeForTypeMap[npcType]; ok {
				return subtypes
			}
			return nil
		},
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
	handler, ok := c.saveHandlers[field]
	if !ok {
		log.Printf("SaveField: unrecognized field %v (%s)", field, field.String())
		return nil
	}
	handler(value)
	if err := c.handleErrors("SaveField", field.String()); err != nil {
		return err
	}

	return nil
}

func (c *NPCEditController) RandomizeField(field cp.CompEnum) string {
	handler, ok := c.randomHandlers[field]
	if !ok {
		log.Printf("RandomizeField: unrecognized field %v (%s)", field, field.String())
		return ""
	}
	value := handler()

	c.handleErrors("RandomizeField", field.String())

	return value
}

func (c *NPCEditController) GetFieldOptions(field cp.CompEnum) []string {
	handler, ok := c.optionHandlers[field]
	if !ok {
		return nil
	}
	return handler()
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
