package service

import (
	"github.com/lackmus/npcgengo/model"
)

// NPCBuilder is a mutable structure used to collect values before constructing the immutable model.NPC.
type NPCBuilder struct {
	ID          string
	Name        string
	Faction     string
	Species     string
	NPCType     string
	NPCSubtype  string
	Trait       string
	Drive       string
	Stats       map[string]int
	Items       map[string]string
	Abilities   map[string]string
	Description string
}

// NewBuilder creates a new instance of NPCBuilder with sensible defaults.
func NewNPCBuilder() *NPCBuilder {
	return &NPCBuilder{
		Stats:     make(map[string]int),
		Items:     make(map[string]string),
		Abilities: make(map[string]string),
	}
}

// NewBuilderFromNPC creates a new instance of NPCBuilder with the values from the provided NPC.
func NewNPCBuilderFromNPC(npc model.NPC) *NPCBuilder {
	return &NPCBuilder{
		ID:          npc.ID(),
		Name:        npc.Name(),
		Faction:     npc.Faction(),
		Species:     npc.Species(),
		NPCType:     npc.NPCType(),
		NPCSubtype:  npc.NPCSubtype(),
		Trait:       npc.Trait(),
		Drive:       npc.Drive(),
		Description: npc.Description(),
		Stats:       npc.Stats(),
		Items:       npc.Items(),
		Abilities:   npc.Abilities(),
	}
}

func (b *NPCBuilder) Build() model.NPC {
	return model.NewNPC(
		b.ID, b.Name, b.Faction, b.Species, b.NPCType, b.NPCSubtype,
		b.Trait, b.Drive, b.Description, b.Stats, b.Items, b.Abilities,
	)
}

// Build constructs the final NPC object from the builder.
func (b *NPCBuilder) BuildWithRandom(rand *RandomizerService) model.NPC {
	// Apply random options for any empty fields
	if b.ID == "" {
		b.WithID(rand.GenerateID())
	}
	if b.NPCType == "" {
		b.NPCType = rand.RandomType()
	}
	if b.NPCSubtype == "" {
		b.NPCSubtype = rand.RandomSubtype(b.NPCType)
	}
	if b.Faction == "" {
		b.Faction = rand.RandomFaction()
	}
	if b.Species == "" {
		b.Species = rand.RandomSpecies()
	}
	if b.Name == "" {
		b.Name = rand.GenerateName(b.Species)
	}
	if b.Trait == "" {
		b.Trait = rand.GenerateTraitDescription()
	}
	if b.Description == "" {
		//TODO: Generate a description
	}
	if len(b.Items) == 0 {
		b.Items = rand.GenerateEquipment(b.NPCSubtype)
	}
	if len(b.Stats) == 0 {
		b.Stats = rand.ApplySubtypeStats(b.NPCSubtype)
	}
	if len(b.Abilities) == 0 {
		//TODO: Generate abilities
	}
	return b.Build()
}

func (b *NPCBuilder) WithID(id string) *NPCBuilder                  { b.ID = id; return b }
func (b *NPCBuilder) WithName(name string) *NPCBuilder              { b.Name = name; return b }
func (b *NPCBuilder) WithFaction(faction string) *NPCBuilder        { b.Faction = faction; return b }
func (b *NPCBuilder) WithSpecies(species string) *NPCBuilder        { b.Species = species; return b }
func (b *NPCBuilder) WithType(npcType string) *NPCBuilder           { b.NPCType = npcType; return b }
func (b *NPCBuilder) WithSubType(npcSubtype string) *NPCBuilder     { b.NPCSubtype = npcSubtype; return b }
func (b *NPCBuilder) WithTrait(trait string) *NPCBuilder            { b.Trait = trait; return b }
func (b *NPCBuilder) WithDrive(drive string) *NPCBuilder            { b.Drive = drive; return b }
func (b *NPCBuilder) WithStats(stats map[string]int) *NPCBuilder    { b.Stats = stats; return b }
func (b *NPCBuilder) WithItems(items map[string]string) *NPCBuilder { b.Items = items; return b }
func (b *NPCBuilder) WithAbilities(abilities map[string]string) *NPCBuilder {
	b.Abilities = abilities
	return b
}
func (b *NPCBuilder) WithDescription(description string) *NPCBuilder {
	b.Description = description
	return b
}
