package service

import (
	"github.com/lackmus/npcgengo/model"
)

type NPCBuilder struct {
	ID, Name, Faction, Species, NPCType, NPCSubType, Trait, Description string
	Components                                                          map[string]string
}

// NewNPCBuilder creates a new NPCBuilder with default values.
func NewNPCBuilder() *NPCBuilder {
	return &NPCBuilder{
		Components: make(map[string]string),
	}
}

// NewNPCBuilderFromNPC initializes a builder from an existing NPC.
func NewNPCBuilderFromNPC(npc model.NPC) *NPCBuilder {
	return &NPCBuilder{
		ID:          npc.ID(),
		Name:        npc.Name(),
		Faction:     npc.Faction(),
		Species:     npc.Species(),
		NPCType:     npc.NPCType(),
		NPCSubType:  npc.NPCSubtype(),
		Trait:       npc.Trait(),
		Description: npc.Description(),
		Components:  npc.Components(),
	}
}

// Build constructs the NPC from the builder.
func (b *NPCBuilder) Build() model.NPC {
	return model.NewNPC(
		b.ID, b.Name, b.Faction, b.Species, b.NPCType, b.NPCSubType,
		b.Trait, b.Description, b.Components,
	)
}

// Setter methods for fluent API
func (b *NPCBuilder) WithID(id string) *NPCBuilder           { b.ID = id; return b }
func (b *NPCBuilder) WithName(name string) *NPCBuilder       { b.Name = name; return b }
func (b *NPCBuilder) WithFaction(faction string) *NPCBuilder { b.Faction = faction; return b }
func (b *NPCBuilder) WithSpecies(species string) *NPCBuilder { b.Species = species; return b }
func (b *NPCBuilder) WithType(t string) *NPCBuilder          { b.NPCType = t; return b }
func (b *NPCBuilder) WithSubType(st string) *NPCBuilder      { b.NPCSubType = st; return b }
func (b *NPCBuilder) WithTrait(trait string) *NPCBuilder     { b.Trait = trait; return b }
func (b *NPCBuilder) WithDescription(desc string) *NPCBuilder {
	b.Description = desc
	return b
}
func (b *NPCBuilder) WithComponent(key, value string) *NPCBuilder {
	b.Components[key] = value
	return b
}
