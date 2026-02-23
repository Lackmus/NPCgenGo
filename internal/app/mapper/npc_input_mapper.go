package mapper

import (
	"strings"

	"github.com/lackmus/npcgengo/pkg/product/model"
	"github.com/lackmus/npcgengo/pkg/product/service"
)

const defaultLocationID = "default"

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

// toTraitString normalizes trait array to comma-separated string
func toTraitString(traits []string) string {
	if len(traits) == 0 {
		return ""
	}
	trimmed := make([]string, 0, len(traits))
	for _, t := range traits {
		if s := strings.TrimSpace(t); s != "" {
			trimmed = append(trimmed, s)
		}
	}
	if len(trimmed) == 0 {
		return ""
	}
	return strings.Join(trimmed, ", ")
}

// ToModelNPC constructs a new NPC from user input using the builder pattern.
func ToModelNPC(input NPCInput, builder *service.NPCBuilder) (model.NPC, error) {
	return ToModelNPCWithOriginal(input, builder, nil)
}

// ToModelNPCWithOriginal builds an NPC from input, loading the original first if provided.
// If original is nil, creates a new NPC. If original exists, loads it first then only applies changed fields.
func ToModelNPCWithOriginal(input NPCInput, builder *service.NPCBuilder, original *model.NPC) (model.NPC, error) {
	locationID := strings.TrimSpace(input.LocationID)
	if locationID == "" {
		locationID = defaultLocationID
	}

	// If editing an existing NPC, load it first to preserve unchanged fields
	if original != nil {
		builder = builder.WithNPC(*original)
	}

	// Set location
	builder.GetNPC().LocationID = locationID

	// Build using single chain - apply all input fields
	npc, err := builder.
		WithType(strings.TrimSpace(input.Type)).
		WithSubtype(strings.TrimSpace(input.Subtype)).
		WithSpecies(strings.TrimSpace(input.Species)).
		WithFaction(strings.TrimSpace(input.Faction)).
		WithName(strings.TrimSpace(input.Name)).
		WithTrait(toTraitString(input.Traits)).
		WithSubtypeStats(strings.TrimSpace(input.Stats)).
		WithSubtypeEquipment(strings.TrimSpace(input.Items)).
		Build()

	if err != nil {
		return model.NPC{}, err
	}

	// Set ID if provided
	if id := strings.TrimSpace(input.ID); id != "" {
		npc.ID = id
	}

	return npc, nil
}
