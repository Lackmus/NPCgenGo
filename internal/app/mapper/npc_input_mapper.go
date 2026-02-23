package mapper

import (
	"strings"

	"github.com/lackmus/npcgengo/pkg/product/model"
)

const defaultLocationID = "default"

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

func ToModelNPC(input NPCInput) model.NPC {
	locationID := strings.TrimSpace(input.LocationID)
	if locationID == "" {
		locationID = defaultLocationID
	}

	npc := model.NPC{
		ID:         strings.TrimSpace(input.ID),
		LocationID: locationID,
	}

	if value := strings.TrimSpace(input.Name); value != "" {
		npc.SetName(value)
	}
	if value := strings.TrimSpace(input.Type); value != "" {
		npc.SetType(value)
	}
	if value := strings.TrimSpace(input.Subtype); value != "" {
		npc.SetSubtype(value)
	}
	if value := strings.TrimSpace(input.Species); value != "" {
		npc.SetSpecies(value)
	}
	if value := strings.TrimSpace(input.Faction); value != "" {
		npc.SetFaction(value)
	}
	if len(input.Traits) > 0 {
		traits := make([]string, 0, len(input.Traits))
		for _, trait := range input.Traits {
			if trimmed := strings.TrimSpace(trait); trimmed != "" {
				traits = append(traits, trimmed)
			}
		}
		if len(traits) > 0 {
			npc.SetTrait(strings.Join(traits, ", "))
		}
	}
	if value := strings.TrimSpace(input.Stats); value != "" {
		npc.SetStats(value)
	}
	if value := strings.TrimSpace(input.Items); value != "" {
		npc.SetItems(value)
	}
	if value := strings.TrimSpace(input.Description); value != "" {
		npc.SetDescription(value)
	}

	return npc
}
