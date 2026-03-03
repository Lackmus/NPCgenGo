package main

import (
	"context"
	"fmt"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/internal/app/mapper"
	"github.com/lackmus/npcgengo/pkg/model"
	"github.com/lackmus/npcgengo/pkg/service"
)

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

func (a *WailsAPI) ListNPCs() []mapper.NPCInput {
	return mapper.ToNPCInputs(a.npcController.GetAllNPCs())
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

func (a *WailsAPI) GetNPC(id string) (mapper.NPCInput, error) {
	npc, err := a.npcController.GetNPCByID(id)
	if err != nil {
		return mapper.NPCInput{}, err
	}
	return mapper.ToNPCInput(npc), nil
}

func (a *WailsAPI) GenerateNPC() (mapper.NPCInput, error) {
	npc, err := a.npcController.CreateRandomNPC()
	if err != nil {
		return mapper.NPCInput{}, err
	}
	return mapper.ToNPCInput(npc), nil
}

func (a *WailsAPI) DeleteNPC(id string) error {
	a.npcController.DeleteNPC(id)
	return nil
}

func (a *WailsAPI) DeleteAllNPCs() error {
	a.npcController.DeleteAllNPCs()
	return nil
}

func (a *WailsAPI) SaveNPC(input mapper.NPCInput) (mapper.NPCInput, error) {
	// Get the original NPC for preserving unchanged fields
	var originalNPC *model.NPC
	if input.ID != "" {
		existing, err := a.npcController.GetNPCByID(input.ID)
		if err != nil {
			return mapper.NPCInput{}, fmt.Errorf("npc with ID %s not found", input.ID)
		}
		originalNPC = &existing
	}

	npc, err := mapper.ToModelNPCWithOriginal(input, a.npcController.GetNPCBuilder(), originalNPC)
	if err != nil {
		return mapper.NPCInput{}, err
	}
	if npc.ID == "" {
		return mapper.NPCInput{}, fmt.Errorf("cannot save without an ID (use Generate to create new NPCs)")
	}
	return a.validateAndPersistNPC(npc, true)
}

func (a *WailsAPI) CreateNPC(input mapper.NPCInput) (mapper.NPCInput, error) {
	npc, err := mapper.ToModelNPC(input, a.npcController.GetNPCBuilder())
	if err != nil {
		return mapper.NPCInput{}, err
	}
	return a.validateAndPersistNPC(npc, false)
}

func (a *WailsAPI) validateAndPersistNPC(npc model.NPC, isUpdate bool) (mapper.NPCInput, error) {
	if err := a.npcController.ValidateNPC(npc); err != nil {
		return mapper.NPCInput{}, err
	}
	if isUpdate {
		a.npcController.UpdateNPC(npc)
	} else {
		a.npcController.AddNPC(npc)
	}
	return mapper.ToNPCInput(npc), nil
}
