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
func NewBuilder() *NPCBuilder {
	return &NPCBuilder{
		Stats:     make(map[string]int),
		Items:     make(map[string]string),
		Abilities: make(map[string]string),
	}
}

// Build constructs the final NPC object from the builder.
func (b *NPCBuilder) Build(rand RandomizerService) model.NPC {
	// Apply random options for any empty fields
	if b.ID == "" {
		b.ID = rand.GenerateID()
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

	// Return the final NPC object
	return model.NewNPC(
		b.ID,
		b.Name,
		b.Faction,
		b.Species,
		b.NPCType,
		b.NPCSubtype,
		b.Trait,
		b.Drive,
		b.Description,
		b.Stats,
		b.Items,
		b.Abilities,
	)
}

// WithID sets the NPC's ID.
func (b *NPCBuilder) WithID(id string) *NPCBuilder {
	b.ID = id
	return b
}

// WithName sets the NPC's name.
func (b *NPCBuilder) WithName(name string) *NPCBuilder {
	b.Name = name
	return b
}

// WithFaction sets the NPC's faction.
func (b *NPCBuilder) WithFaction(faction string) *NPCBuilder {
	b.Faction = faction
	return b
}

// WithSpecies sets the NPC's species.
func (b *NPCBuilder) WithSpecies(species string) *NPCBuilder {
	b.Species = species
	return b
}

// WithType sets the NPC's type.
func (b *NPCBuilder) WithType(npcType string) *NPCBuilder {
	b.NPCType = npcType
	return b
}

// WithSubType sets the NPC's subtype.
func (b *NPCBuilder) WithSubType(npcSubtype string) *NPCBuilder {
	b.NPCSubtype = npcSubtype
	return b
}

// WithTrait sets the NPC's trait.
func (b *NPCBuilder) WithTrait(trait string) *NPCBuilder {
	b.Trait = trait
	return b
}

// WithDrive sets the NPC's drive.
func (b *NPCBuilder) WithDrive(drive string) *NPCBuilder {
	b.Drive = drive
	return b
}

// WithStats sets the NPC's stats.
func (b *NPCBuilder) WithStats(stats map[string]int) *NPCBuilder {
	b.Stats = stats
	return b
}

// WithItems sets the NPC's items.
func (b *NPCBuilder) WithItems(items map[string]string) *NPCBuilder {
	b.Items = items
	return b
}

// WithAbilities sets the NPC's abilities.
func (b *NPCBuilder) WithAbilities(abilities map[string]string) *NPCBuilder {
	b.Abilities = abilities
	return b
}

// WithDescription sets the NPC's description.
func (b *NPCBuilder) WithDescription(description string) *NPCBuilder {
	b.Description = description
	return b
}

// NewOptionNPC uses a builder to create a new immutable NPC by applying the provided options.
func NewOptionNPC(rand RandomizerService, opts ...func(*NPCBuilder) *NPCBuilder) model.NPC {
	builder := NewBuilder()
	for _, opt := range opts {
		opt(builder)
	}
	return builder.Build(rand)
}
