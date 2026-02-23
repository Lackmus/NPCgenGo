// Description: This file contains the controller for the list of NPCs.
package controllers

import (
	"fmt"
	"log"

	h "github.com/lackmus/npcgengo/internal/platform/helpers"
	"github.com/lackmus/npcgengo/pkg/product/model"
	"github.com/lackmus/npcgengo/pkg/product/service"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

type NPCListController struct {
	NpcService       *service.NPCService
	CreationSupplier *service.NPCCreationSupplier
	validator        *service.NPCValidationService
	observers        []shared.NPCObserver
	CreationOptions  *service.NPCCreationOptions
}

func NewNPCListController(creationSupplier *service.NPCCreationSupplier, npcService *service.NPCService) *NPCListController {
	log.Println("Creating NPCListController...")
	validator := service.NewNPCValidationService(creationSupplier.CreationDataService)
	return &NPCListController{
		NpcService:       npcService,
		CreationSupplier: creationSupplier,
		validator:        validator,
		CreationOptions:  creationSupplier.CreationOptions,
		observers:        []shared.NPCObserver{},
	}
}

func (c *NPCListController) GetNPCBuilder() *service.NPCBuilder {
	return service.NewNPCBuilder(c.CreationSupplier)
}

func (c *NPCListController) CreateRandomNPC() (model.NPC, error) {
	return c.createAndAddNPC(h.Random, h.Random)
}

func (c *NPCListController) CreateNPC(npctype string, faction string) (model.NPC, error) {
	return c.createAndAddNPC(npctype, faction)
}

func (c *NPCListController) createAndAddNPC(npctype string, faction string) (model.NPC, error) {
	npc, err := service.CreateNPCWithOptions(npctype, faction, c.CreationSupplier)
	if err != nil {
		log.Printf("Error creating NPC: %v", err)
		return model.NPC{}, err
	}
	c.AddNpc(npc)
	return npc, nil
}

func (c *NPCListController) InitView(view shared.NPCListViewer) {
	log.Println("Initializing view...")
	c.RegisterObserver(view)
	c.NotifyObservers()
}

func (c *NPCListController) UpdateNpc(npc model.NPC) {
	c.AddNpc(npc)
}

func (c *NPCListController) RegisterObserver(o shared.NPCObserver) {
	c.observers = append(c.observers, o)
}

func (c *NPCListController) NotifyObservers() {
	npcs := c.NpcService.GetAllNPC()
	for _, o := range c.observers {
		o.Update(npcs)
	}
}

func (c *NPCListController) GetAllNpcs() []model.NPC {
	npcs := c.NpcService.GetAllNPC()
	if len(npcs) == 0 {
		log.Println("No NPCs found.")
	}
	return npcs
}

func (c *NPCListController) GetNpcByID(id string) (model.NPC, error) {
	npc, err := c.NpcService.GetNPCByID(id)
	if err != nil {
		return model.NPC{}, err
	}
	return npc, nil
}

func (c *NPCListController) AddNpc(npc model.NPC) {
	c.NpcService.AddNPC(npc)
	c.NotifyObservers()
}

func (c *NPCListController) DeleteNPC(id string) {
	if err := c.NpcService.DeleteNPC(id); err != nil {
		log.Printf("Error deleting NPC: %v", err)
	}
	c.NotifyObservers()
}

func (c *NPCListController) DeleteAllNPC() {
	c.NpcService.DeleteAllNPC()
	c.NotifyObservers()
}

func (c *NPCListController) CreateNPCGroup() {
	// Implementation will be in the view that uses this controller
}

// GetCreationOptions returns the available creation options for NPC generation.
func (c *NPCListController) GetCreationOptions() *service.NPCCreationOptions {
	if c.CreationSupplier == nil {
		return nil
	}
	return c.CreationSupplier.CreationOptions
}

// ValidateNPC validates an NPC using the controller's validation service.
func (c *NPCListController) ValidateNPC(npc model.NPC) error {
	if c.validator == nil {
		return nil
	}
	return c.validator.ValidateNPC(npc)
}

// GetSubtypeFields returns the rolled stats and items for a given subtype.
func (c *NPCListController) GetSubtypeFields(subtype string) (stats, items string, err error) {
	if c.CreationSupplier == nil {
		return "", "", fmt.Errorf("creation supplier not initialized")
	}
	subtypeData, err := c.CreationSupplier.CreationDataService.GetNpcSubtypeData(subtype)
	if err != nil {
		return "", "", err
	}
	return subtypeData.GetStats(), subtypeData.GetEquipment(), nil
}

// GetSpeciesName returns a generated name for a given species.
func (c *NPCListController) GetSpeciesName(species string) (string, error) {
	if c.CreationSupplier == nil {
		return "", fmt.Errorf("creation supplier not initialized")
	}
	speciesData, err := c.CreationSupplier.CreationDataService.GetSpeciesData(species)
	if err != nil {
		return "", err
	}
	nameData, err := c.CreationSupplier.CreationDataService.GetNameData(speciesData.NameSource)
	if err != nil {
		return "", err
	}
	return nameData.GenerateName(), nil
}
