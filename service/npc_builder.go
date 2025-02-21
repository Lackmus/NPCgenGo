package service

import (
	"maps"

	"github.com/lackmus/npcgengo/model"
)

type NPCBuilder struct {
	ID, Name, Faction, Species, NPCType, NPCSubType, Trait, Drive, Description string
	Stats                                                                      map[string]int
	Items, Abilities                                                           map[string]string
}

// NewNPCBuilder creates a new NPCBuilder with default values.
func NewNPCBuilder() *NPCBuilder {
	return &NPCBuilder{
		Stats:     make(map[string]int),
		Items:     make(map[string]string),
		Abilities: make(map[string]string),
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
		Drive:       npc.Drive(),
		Description: npc.Description(),
		Stats:       maps.Clone(npc.Stats()),
		Items:       maps.Clone(npc.Items()),
		Abilities:   maps.Clone(npc.Abilities()),
	}
}

// Build constructs the NPC from the builder.
func (b *NPCBuilder) Build() model.NPC {
	return model.NewNPC(
		b.ID, b.Name, b.Faction, b.Species, b.NPCType, b.NPCSubType,
		b.Trait, b.Drive, b.Description, b.Stats, b.Items, b.Abilities,
	)
}

// BuildWithRandom fills missing fields using randomization.
func (b *NPCBuilder) BuildWithRandom(rand *RandomizerService) model.NPC {
	if b.ID == "" {
		b.WithID(rand.GenerateID())
	}
	if b.NPCType == "" {
		b.WithType(rand.RandomType())
	}
	if b.NPCSubType == "" {
		b.WithSubType(rand.RandomSubtype(b.NPCType))
	}
	if b.Faction == "" {
		b.WithFaction(rand.RandomFaction())
	}
	if b.Species == "" {
		b.WithSpecies(rand.RandomSpecies())
	}
	if b.Name == "" {
		b.WithName(rand.GenerateName(b.Species))
	}
	if b.Trait == "" {
		b.WithTrait(rand.GenerateTraitDescription())
	}
	if b.Description == "" {
		//nb.WithDescription(rand.GenerateDescription(nb.Name, nb.Species, nb.Type))
	}
	if len(b.Items) == 0 {
		b.WithItems(rand.GenerateEquipment(b.NPCSubType))
	}
	if len(b.Stats) == 0 {
		b.WithStats(rand.ApplySubtypeStats(b.NPCSubType))
	}
	if len(b.Abilities) == 0 {
		//nb.WithAbilities(rand.GenerateAbilities(nb.SubType))
	}

	return b.Build()
}

// Setter methods for fluent API
func (b *NPCBuilder) WithID(id string) *NPCBuilder           { b.ID = id; return b }
func (b *NPCBuilder) WithName(name string) *NPCBuilder       { b.Name = name; return b }
func (b *NPCBuilder) WithFaction(faction string) *NPCBuilder { b.Faction = faction; return b }
func (b *NPCBuilder) WithSpecies(species string) *NPCBuilder { b.Species = species; return b }
func (b *NPCBuilder) WithType(t string) *NPCBuilder          { b.NPCType = t; return b }
func (b *NPCBuilder) WithSubType(st string) *NPCBuilder      { b.NPCSubType = st; return b }
func (b *NPCBuilder) WithTrait(trait string) *NPCBuilder     { b.Trait = trait; return b }
func (b *NPCBuilder) WithDrive(drive string) *NPCBuilder     { b.Drive = drive; return b }
func (b *NPCBuilder) WithStats(stats map[string]int) *NPCBuilder {
	b.Stats = maps.Clone(stats)
	return b
}
func (b *NPCBuilder) WithItems(items map[string]string) *NPCBuilder {
	b.Items = maps.Clone(items)
	return b
}
func (b *NPCBuilder) WithAbilities(abilities map[string]string) *NPCBuilder {
	b.Abilities = maps.Clone(abilities)
	return b
}
func (b *NPCBuilder) WithDescription(desc string) *NPCBuilder {
	b.Description = desc
	return b
}
