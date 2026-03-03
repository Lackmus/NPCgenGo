package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/pkg/model"
	cp "github.com/lackmus/npcgengo/pkg/model/npc_components"
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

type App struct {
	ctx           context.Context
	npcController *controllers.NPCListController
}

func NewApp(npcController *controllers.NPCListController) *App {
	return &App{npcController: npcController}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListNPCs() []model.NPC {
	return a.npcController.GetAllNpcs()
}

func (a *App) GetNPC(id string) (model.NPC, error) {
	return a.npcController.GetNpcByID(id)
}

func (a *App) GenerateNPC() (model.NPC, error) {
	return a.npcController.CreateRandomNPC()
}

func (a *App) DeleteNPC(id string) error {
	a.npcController.DeleteNPC(id)
	return nil
}

func (a *App) DeleteAllNPCs() error {
	a.npcController.DeleteAllNPC()
	return nil
}

func (a *App) SaveNPC(input NPCInput) (model.NPC, error) {
	npc := toModelNPC(input)
	if npc.ID == "" {
		npc.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	a.npcController.UpdateNpc(npc)
	return npc, nil
}

func (a *App) CreateNPC(input NPCInput) (model.NPC, error) {
	npc := toModelNPC(input)
	if npc.ID == "" {
		npc.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	a.npcController.AddNpc(npc)
	return npc, nil
}

func toModelNPC(input NPCInput) model.NPC {
	locationID := strings.TrimSpace(input.LocationID)
	if locationID == "" {
		locationID = "default"
	}

	npc := model.NPC{
		ID:         strings.TrimSpace(input.ID),
		LocationID: locationID,
		Components: make(map[cp.CompEnum]string),
	}

	if value := strings.TrimSpace(input.Name); value != "" {
		npc.Components[cp.CompName] = value
	}
	if value := strings.TrimSpace(input.Type); value != "" {
		npc.Components[cp.CompType] = value
	}
	if value := strings.TrimSpace(input.Subtype); value != "" {
		npc.Components[cp.CompSubtype] = value
	}
	if value := strings.TrimSpace(input.Species); value != "" {
		npc.Components[cp.CompSpecies] = value
	}
	if value := strings.TrimSpace(input.Faction); value != "" {
		npc.Components[cp.CompFaction] = value
	}
	if len(input.Traits) > 0 {
		traits := make([]string, 0, len(input.Traits))
		for _, trait := range input.Traits {
			if trimmed := strings.TrimSpace(trait); trimmed != "" {
				traits = append(traits, trimmed)
			}
		}
		if len(traits) > 0 {
			npc.Components[cp.CompTrait] = strings.Join(traits, ", ")
		}
	}
	if value := strings.TrimSpace(input.Stats); value != "" {
		npc.Components[cp.CompStats] = value
	}
	if value := strings.TrimSpace(input.Items); value != "" {
		npc.Components[cp.CompItems] = value
	}
	if value := strings.TrimSpace(input.Description); value != "" {
		npc.Components[cp.CompDescription] = value
	}

	return npc
}
