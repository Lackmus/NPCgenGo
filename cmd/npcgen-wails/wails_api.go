package main

import (
	"context"
	"fmt"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/internal/app/mapper"
	"github.com/lackmus/npcgengo/pkg/product/model"
	"github.com/lackmus/npcgengo/pkg/product/service"
)

type NPCInput struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Subtype     string   `json:"subtype"`
	Species     string   `json:"species"`
	Faction     string   `json:"faction"`
	Traits      []string `json:"traits"`
	Stats       string   `json:"stats"`
	Items       string   `json:"items"`
	Description string   `json:"description"`
	LocationID  string   `json:"locationID"`
}

type SubtypeRoll struct {
	Stats       string `json:"stats"`
	Items       string `json:"items"`
	Description string `json:"description"`
}

type WailsAPI struct {
	ctx           context.Context
	npcController *controllers.NPCListController
	validator     *service.NPCValidationService
}

func NewWailsAPI(npcController *controllers.NPCListController) *WailsAPI {
	return &WailsAPI{
		npcController: npcController,
		validator:     service.NewNPCValidationService(npcController.CreationSupplier.CreationDataService),
	}
}

func (a *WailsAPI) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *WailsAPI) ListNPCs() []model.NPC {
	return a.npcController.GetAllNpcs()
}

func (a *WailsAPI) GetCreationOptions() *service.NPCCreationOptions {
	return a.npcController.CreationSupplier.CreationOptions
}

func (a *WailsAPI) RollSubtypeFields(subtype string) (SubtypeRoll, error) {
	subtypeData, err := a.npcController.CreationSupplier.CreationDataService.GetNpcSubtypeData(subtype)
	if err != nil {
		return SubtypeRoll{}, err
	}

	return SubtypeRoll{
		Stats:       subtypeData.GetStats(),
		Items:       subtypeData.GetEquipment(),
		Description: subtypeData.GetDescription(),
	}, nil
}

func (a *WailsAPI) RollSpeciesName(species string) (string, error) {
	speciesData, err := a.npcController.CreationSupplier.CreationDataService.GetSpeciesData(species)
	if err != nil {
		return "", err
	}

	nameData, err := a.npcController.CreationSupplier.CreationDataService.GetNameData(speciesData.NameSource)
	if err != nil {
		return "", err
	}

	return nameData.GenerateName(), nil
}

func (a *WailsAPI) GetNPC(id string) (model.NPC, error) {
	return a.npcController.GetNpcByID(id)
}

func (a *WailsAPI) GenerateNPC() (model.NPC, error) {
	return a.npcController.CreateRandomNPC()
}

func (a *WailsAPI) DeleteNPC(id string) error {
	a.npcController.DeleteNPC(id)
	return nil
}

func (a *WailsAPI) DeleteAllNPCs() error {
	a.npcController.DeleteAllNPC()
	return nil
}

func (a *WailsAPI) SaveNPC(input NPCInput) (model.NPC, error) {
	npc := mapper.ToModelNPC(toMapperInput(input))
	if npc.ID == "" {
		return model.NPC{}, fmt.Errorf("cannot save without an ID (use Generate to create new NPCs)")
	}
	if err := a.validator.ValidateNPC(npc); err != nil {
		return model.NPC{}, err
	}
	a.npcController.UpdateNpc(npc)
	return npc, nil
}

func (a *WailsAPI) CreateNPC(input NPCInput) (model.NPC, error) {
	npc := mapper.ToModelNPC(toMapperInput(input))
	if err := a.validator.ValidateNPC(npc); err != nil {
		return model.NPC{}, err
	}
	a.npcController.AddNpc(npc)
	return npc, nil
}

func toMapperInput(input NPCInput) mapper.NPCInput {
	return mapper.NPCInput{
		ID:          input.ID,
		Name:        input.Name,
		Type:        input.Type,
		Subtype:     input.Subtype,
		Species:     input.Species,
		Faction:     input.Faction,
		Traits:      input.Traits,
		Stats:       input.Stats,
		Items:       input.Items,
		Description: input.Description,
		LocationID:  input.LocationID,
	}
}
