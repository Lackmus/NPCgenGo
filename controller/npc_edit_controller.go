// Package controller provides controllers for managing NPC editing.
package controller

import (
	"errors"
	"log"

	"github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/shared"
)

// NPCEditController manages the editing of an NPC.
type NPCEditController struct {
	creationSupplier  *service.NPCCreationSupplier
	rand              *service.RandomizerService
	npcBuilder        *service.NPCBuilder
	npcListController *NPCListController
	observers         []shared.NPCEditObserver
	errors            []error
}

// NewNPCEditController creates a new NPCEditController with the provided view, creation supplier, and list controller.
func NewNPCEditController(view shared.NPCEditViewer, creationSupplier *service.NPCCreationSupplier, controller *NPCListController) *NPCEditController {
	return &NPCEditController{
		creationSupplier:  creationSupplier,
		rand:              creationSupplier.RandomizerService,
		npcBuilder:        service.NewNPCBuilder(creationSupplier),
		npcListController: controller,
		observers:         []shared.NPCEditObserver{view},
		errors:            []error{},
	}
}

// EditNPC returns the current NPC from the builder for editing.
func (c *NPCEditController) EditNPC() model.NPC {
	return *c.npcBuilder.GetNPC()
}

// RegisterObserver registers a new observer to receive NPC edit updates.
func (c *NPCEditController) RegisterObserver(o shared.NPCEditObserver) {
	c.observers = append(c.observers, o)
}

// NotifyObservers notifies all registered observers that the NPC has been updated.
func (c *NPCEditController) NotifyObservers() {
	for _, o := range c.observers {
		o.UpdateNPC(*c.npcBuilder.GetNPC())
	}
}

// NotifyObserversField notifies observers that a specific field has been updated.
func (c *NPCEditController) NotifyObserversField(field cp.CompEnum, value any) {
	for _, o := range c.observers {
		o.UpdateField(field, value)
	}
}

// notify observers of error
func (c *NPCEditController) NotifyObserversError(err error) {
	for _, o := range c.observers {
		o.OnNPCEditError(err)
	}
	log.Printf("Error: %v", err)
}

// LoadNPC loads an NPC into the builder and notifies observers.
func (c *NPCEditController) LoadNPC(npc model.NPC) {
	c.npcBuilder.WithNPC(npc)
	c.NotifyObservers()
}

// SaveNPC builds the current NPC from the builder and updates the NPC list.
// It logs an error if the build fails.
func (c *NPCEditController) SaveNPC() {
	npc, err := c.npcBuilder.Build()
	if err != nil {
		c.handleErrors("SaveNPC", "NPC")
		return
	}
	c.NotifyObservers()
	c.npcListController.UpdateNpc(npc)
}

// SaveField updates a single field in the NPC builder using the provided value
// and notifies observers of the change.
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
	default:
		log.Printf("SaveField: unrecognized field %v (%s)", field, field.String())
		return nil
	}
	// Handle any errors that may have been added to the builder.
	if err := c.handleErrors("SaveField", field.String()); err != nil {
		return err
	}

	c.NotifyObserversField(field, value)
	return nil
}

// RandomizeField updates a single field in the NPC builder with a random value,
// notifies observers, and returns the new value.
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
	default:
		log.Printf("RandomizeField: unrecognized field %v (%s)", field, field.String())
		return ""
	}

	// lookup c.npcBuilder.Errors() and do something with it
	c.handleErrors("RandomizeField", field.String())

	c.NotifyObserversField(field, value)
	return value
}

// GetFieldOptions returns available options for a given field using the creation supplier.
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
		// Get the current NPC type from the builder (adjust method name if needed).
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

// handleErrors checks for builder errors, aggregates them, notifies observers,
// and returns the aggregated error.
func (c *NPCEditController) handleErrors(context string, field string) error {
	if c.npcBuilder.HasErrors() {
		// Aggregate errors from the builder.
		err := errors.Join(c.npcBuilder.Errors()...)
		log.Printf("%s: error saving field %v: %v", context, field, err)
		c.NotifyObserversError(err)
		// Optionally, you can clear builder errors here if desired:
		// c.npcBuilder.ClearErrors()
		return err
	}
	return nil
}
