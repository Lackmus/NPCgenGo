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
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Subtype    string   `json:"subtype"`
	Species    string   `json:"species"`
	Faction    string   `json:"faction"`
	Traits     []string `json:"traits"`
	Stats      string   `json:"stats"`
	Items      string   `json:"items"`
	LocationID string   `json:"locationID"`
}

type SubtypeRoll struct {
	Stats string `json:"stats"`
	Items string `json:"items"`
}

type WailsAPI struct {
	ctx           context.Context
	npcController *controllers.NPCListController
}

func NewWailsAPI(npcController *controllers.NPCListController) *WailsAPI {
	return &WailsAPI{
		npcController: npcController,
	}
}

func (a *WailsAPI) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *WailsAPI) ListNPCs() []model.NPC {
	return a.npcController.GetAllNpcs()
}

func (a *WailsAPI) GetCreationOptions() *service.NPCCreationOptions {
	return a.npcController.GetCreationOptions()
}

func (a *WailsAPI) RollSubtypeFields(subtype string) (SubtypeRoll, error) {
	stats, items, err := a.npcController.GetSubtypeFields(subtype)
	if err != nil {
		return SubtypeRoll{}, err
	}

	return SubtypeRoll{
		Stats: stats,
		Items: items,
	}, nil
}

func (a *WailsAPI) RollSpeciesName(species string) (string, error) {
	return a.npcController.GetSpeciesName(species)
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
	// Get the original NPC for preserving unchanged fields
	var originalNPC *model.NPC
	if input.ID != "" {
		if existing, err := a.npcController.GetNpcByID(input.ID); err == nil {
			originalNPC = &existing
		}
	}

	npc, err := mapper.ToModelNPCWithOriginal(toMapperInput(input), a.npcController.GetNPCBuilder(), originalNPC)
	if err != nil {
		return model.NPC{}, err
	}
	if npc.ID == "" {
		return model.NPC{}, fmt.Errorf("cannot save without an ID (use Generate to create new NPCs)")
	}
	if err := a.npcController.ValidateNPC(npc); err != nil {
		return model.NPC{}, err
	}
	a.npcController.UpdateNpc(npc)
	return npc, nil
}

func (a *WailsAPI) CreateNPC(input NPCInput) (model.NPC, error) {
	npc, err := mapper.ToModelNPC(toMapperInput(input), a.npcController.GetNPCBuilder())
	if err != nil {
		return model.NPC{}, err
	}
	if err := a.npcController.ValidateNPC(npc); err != nil {
		return model.NPC{}, err
	}
	a.npcController.AddNpc(npc)
	return npc, nil
}

func toMapperInput(input NPCInput) mapper.NPCInput {
	return mapper.NPCInput{
		ID:         input.ID,
		Name:       input.Name,
		Type:       input.Type,
		Subtype:    input.Subtype,
		Species:    input.Species,
		Faction:    input.Faction,
		Traits:     input.Traits,
		Stats:      input.Stats,
		Items:      input.Items,
		LocationID: input.LocationID,
	}
}
