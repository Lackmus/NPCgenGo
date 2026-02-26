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
	Notes   string `json:"notes"`
}

func (input NPCInput) normalized() NPCInput {
	input.ID = strings.TrimSpace(input.ID)
	input.Name = strings.TrimSpace(input.Name)
	input.Type = strings.TrimSpace(input.Type)
	input.Subtype = strings.TrimSpace(input.Subtype)
	input.Species = strings.TrimSpace(input.Species)
	input.Faction = strings.TrimSpace(input.Faction)
	input.Trait = strings.TrimSpace(input.Trait)
	input.Stats = strings.TrimSpace(input.Stats)
	input.Items = strings.TrimSpace(input.Items)
	input.Notes = strings.TrimSpace(input.Notes)
	return input
}

func ToNPCInput(npc model.NPC) NPCInput {
	return NPCInput{
		ID:      npc.ID,
		Name:    npc.Name(),
		Type:    npc.Type(),
		Subtype: npc.Subtype(),
		Species: npc.Species(),
		Faction: npc.Faction(),
		Trait:   npc.Trait(),
		Stats:   npc.Stats(),
		Items:   npc.Items(),
		Notes:   npc.Notes(),
	}.normalized()
}

func ToNPCInputs(npcs []model.NPC) []NPCInput {
	out := make([]NPCInput, 0, len(npcs))
	for _, npc := range npcs {
		out = append(out, ToNPCInput(npc))
	}
	return out
}

// ToModelNPC constructs a new NPC from user input using the builder pattern.
func ToModelNPC(input NPCInput, builder *service.NPCBuilder) (model.NPC, error) {
	return ToModelNPCWithOriginal(input, builder, nil)
}

// ToModelNPCWithOriginal builds an NPC from input, loading the original first if provided.
// If original is nil, creates a new NPC. If original exists, loads it first then only applies changed fields.
func ToModelNPCWithOriginal(input NPCInput, builder *service.NPCBuilder, original *model.NPC) (model.NPC, error) {
	if original != nil {
		builder = builder.WithNPC(*original)
	}

	input = input.normalized()

	name := preserveOriginalValue(input.Name, original, cp.CompName)
	npcType := preserveOriginalValue(input.Type, original, cp.CompType)
	subtype := preserveOriginalValue(input.Subtype, original, cp.CompSubtype)
	species := preserveOriginalValue(input.Species, original, cp.CompSpecies)
	faction := preserveOriginalValue(input.Faction, original, cp.CompFaction)
	trait := preserveOriginalValue(input.Trait, original, cp.CompTrait)
	stats := preserveOriginalValue(input.Stats, original, cp.CompStats)
	items := preserveOriginalValue(input.Items, original, cp.CompItems)
	notes := preserveOriginalValue(input.Notes, original, cp.CompNotes)

	// Build using single chain - apply all input fields
	return builder.
		WithType(npcType).
		WithSubtype(subtype).
		WithSpecies(species).
		WithFaction(faction).
		WithName(name).
		WithTrait(trait).
		WithSubtypeStats(stats).
		WithID(input.ID).
		WithSubtypeEquipment(items).
		WithNotes(notes).
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
