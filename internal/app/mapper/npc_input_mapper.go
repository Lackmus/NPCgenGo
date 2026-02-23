package mapper

import (
	"strings"

	"github.com/lackmus/npcgengo/pkg/product/model"
	cp "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
	"github.com/lackmus/npcgengo/pkg/product/service"
)

type NPCInput struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Species string `json:"species"`
	Faction string `json:"faction"`
	Trait   string `json:"trait"`
	Stats   string `json:"stats"`
	Items   string `json:"items"`
}

// ToModelNPC constructs a new NPC from user input using the builder pattern.
func ToModelNPC(input NPCInput, builder *service.NPCBuilder) (model.NPC, error) {
	return ToModelNPCWithOriginal(input, builder, nil)
}

// ToModelNPCWithOriginal builds an NPC from input, loading the original first if provided.
// If original is nil, creates a new NPC. If original exists, loads it first then only applies changed fields.
func ToModelNPCWithOriginal(input NPCInput, builder *service.NPCBuilder, original *model.NPC) (model.NPC, error) {
	// If editing an existing NPC, load it first to preserve unchanged fields
	if original != nil {
		builder = builder.WithNPC(*original)
	}

	name := preserveOriginalValue(strings.TrimSpace(input.Name), original, cp.CompName)
	npcType := preserveOriginalValue(strings.TrimSpace(input.Type), original, cp.CompType)
	subtype := preserveOriginalValue(strings.TrimSpace(input.Subtype), original, cp.CompSubtype)
	species := preserveOriginalValue(strings.TrimSpace(input.Species), original, cp.CompSpecies)
	faction := preserveOriginalValue(strings.TrimSpace(input.Faction), original, cp.CompFaction)
	trait := preserveOriginalValue(strings.TrimSpace(input.Trait), original, cp.CompTrait)
	stats := preserveOriginalValue(strings.TrimSpace(input.Stats), original, cp.CompStats)
	items := preserveOriginalValue(strings.TrimSpace(input.Items), original, cp.CompItems)

	// Build using single chain - apply all input fields
	return builder.
		WithType(npcType).
		WithSubtype(subtype).
		WithSpecies(species).
		WithFaction(faction).
		WithName(name).
		WithTrait(trait).
		WithSubtypeStats(stats).
		WithID(strings.TrimSpace(input.ID)).
		WithSubtypeEquipment(items).
		Build()
}

func preserveOriginalValue(inputValue string, original *model.NPC, component cp.CompEnum) string {
	if inputValue != "" {
		return inputValue
	}
	if original == nil {
		return ""
	}
	return strings.TrimSpace(original.GetComponent(component))
}
