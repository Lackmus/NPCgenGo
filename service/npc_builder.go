package service

import (
	"log"

	m "github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components" // assumed package for component types/keys
)

// NPCBuilder constructs an NPC step by step.
type NPCBuilder struct {
	npc     *m.NPC
	c       *NPCCreationSupplier
	subtype *cp.NPCSubtype
	species *cp.Species
	npctype string
}

// NewNPCBuilder creates a new NPCBuilder.
func NewNPCBuilder(c *NPCCreationSupplier) *NPCBuilder {
	return &NPCBuilder{
		npc: m.NewNPC(c.RandomizerService.GenerateID()),
		c:   c,
	}
}

// WithRandomType sets a random type for the NPC.
func (b *NPCBuilder) WithType(t string) *NPCBuilder {
	data := b.c.CreationDataService.GetNpcTypeData(t)
	b.npctype = t
	b.npc.AddComponent(*data.NewNPCTypeComponent())
	return b
}

// WithRandomType sets a random type for the NPC.
func (b *NPCBuilder) WithRandomType() *NPCBuilder {
	t := b.c.RandomizerService.RandomType()
	return b.WithType(t)
}

// WithRandomSubtype sets a random subtype for the NPC.
func (b *NPCBuilder) WithSubtype(t string) *NPCBuilder {
	data := b.c.CreationDataService.GetNpcSubtypeData(t)
	b.subtype = &data
	b.npc.AddComponent(*b.subtype.NewNPCSubtypeComponent())
	return b
}

// WithRandomSubtype sets a random subtype for the NPC.
func (b *NPCBuilder) WithRandomSubtype() *NPCBuilder {
	if b.npctype == "" {
		b.WithRandomType()
	}
	s := b.c.RandomizerService.RandomSubtype(b.npctype)
	return b.WithSubtype(s)
}

// WithRandomSubtypeStats sets random stats for the NPC subtype.
func (b *NPCBuilder) WithSubtypeStats() *NPCBuilder {
	if b.subtype == nil {
		log.Fatal("Subtype must be set before stats can be added.")
	}
	b.npc.AddComponent(*b.subtype.NewNPCSubtypeComponentWithStats())
	return b
}

// WithRandomSubtypeEquipment sets random equipment for the NPC subtype.
func (b *NPCBuilder) WithSubtypeEquipment() *NPCBuilder {
	if b.subtype == nil {
		log.Fatal("Subtype must be set before equipment can be added.")
	}
	b.npc.AddComponent(*b.subtype.NewNPCSubtypeComponentWithEquipment())
	return b
}

// WithRandomSpeciesAndName sets a random species and name for the NPC.
func (b *NPCBuilder) WithSpecies(species string) *NPCBuilder {
	data := b.c.CreationDataService.GetSpeciesData(species)
	b.species = &data
	b.npc.AddComponent(*b.species.NewSpeciesComponent())
	return b
}

// WithRandomSpeciesAndName sets a random species and name for the NPC.
func (b *NPCBuilder) WithRandomSpecies() *NPCBuilder {
	s := b.c.RandomizerService.RandomSpecies()
	return b.WithSpecies(s)
}

// WithRandomName sets a random name for the NPC.
func (b *NPCBuilder) WithName() *NPCBuilder {
	if b.species == nil {
		log.Fatal("Species must be set before name can be added.")
	}
	data := b.c.CreationDataService.GetNameData(b.species.NameSource)
	b.npc.AddComponent(*data.NewNameComponent())
	return b
}

// WithRandomFaction sets a random faction for the NPC.
func (b *NPCBuilder) WithFaction(f string) *NPCBuilder {
	data := b.c.CreationDataService.GetFactionData(f)
	b.npc.AddComponent(*data.NewFactionComponent())
	return b
}

// WithRandomFaction sets a random faction for the NPC.
func (b *NPCBuilder) WithRandomFaction() *NPCBuilder {
	f := b.c.RandomizerService.RandomFaction()
	return b.WithFaction(f)
}

// WithRandomTrait sets a random trait for the NPC.
func (b *NPCBuilder) WithTrait(t string) *NPCBuilder {
	data := b.c.CreationDataService.GetTraitData(t)
	b.npc.AddComponent(*data.NewTraitComponent())
	return b
}

// WithRandomTrait sets a random trait for the NPC.
func (b *NPCBuilder) WithRandomTrait() *NPCBuilder {
	t := b.c.RandomizerService.RandomTrait()
	return b.WithTrait(t)
}

// Build constructs the NPC.
func (b *NPCBuilder) Build() m.NPC {
	return *b.npc
}
