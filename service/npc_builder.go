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

// build from npc
func NewNPCBuilderFromNPC(c *NPCCreationSupplier, npc m.NPC) *NPCBuilder {
	subtype := c.CreationDataService.GetNpcSubtypeData(npc.GetComponent(cp.CompSubtype))
	species := c.CreationDataService.GetSpeciesData(npc.GetComponent(cp.CompSpecies))
	return &NPCBuilder{
		npc:     &npc,
		c:       c,
		subtype: &subtype,
		species: &species,
		npctype: npc.GetComponent(cp.CompType),
	}
}

// WithRandomType sets a random type for the NPC.
func (b *NPCBuilder) WithType(npctype string) *NPCBuilder {
	data := b.c.CreationDataService.GetNpcTypeData(npctype)
	b.npctype = npctype
	b.npc.AddComponent(*data.NewNPCTypeComponent())
	return b
}

// WithRandomType sets a random type for the NPC.
func (b *NPCBuilder) WithRandomType() *NPCBuilder {
	npctype := b.c.RandomizerService.RandomType()
	return b.WithType(npctype)
}

// WithRandomSubtype sets a random subtype for the NPC.
func (b *NPCBuilder) WithSubtype(subtype string) *NPCBuilder {
	data := b.c.CreationDataService.GetNpcSubtypeData(subtype)
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

func (b *NPCBuilder) WithSubtypeStats(stats string) *NPCBuilder {
	b.npc.AddComponent(cp.Component{Name: cp.CompStats, Value: stats})
	return b
}

// WithRandomSubtypeStats sets random stats for the NPC subtype.
func (b *NPCBuilder) WithRandomSubtypeStats() *NPCBuilder {
	if b.subtype == nil {
		log.Fatal("Subtype must be set before stats can be added.")
	}
	b.npc.AddComponent(*b.subtype.NewNPCSubtypeComponentWithStats())
	return b
}

// withsubtypeequipment
func (b *NPCBuilder) WithSubtypeEquipment(ítems string) *NPCBuilder {
	b.npc.AddComponent(cp.Component{Name: cp.CompItems, Value: ítems})
	return b
}

// WithRandomSubtypeEquipment sets random equipment for the NPC subtype.
func (b *NPCBuilder) WithRandomSubtypeEquipment() *NPCBuilder {
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
	species := b.c.RandomizerService.RandomSpecies()
	return b.WithSpecies(species)
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
func (b *NPCBuilder) WithFaction(faction string) *NPCBuilder {
	data := b.c.CreationDataService.GetFactionData(faction)
	b.npc.AddComponent(*data.NewFactionComponent())
	return b
}

// WithRandomFaction sets a random faction for the NPC.
func (b *NPCBuilder) WithRandomFaction() *NPCBuilder {
	faction := b.c.RandomizerService.RandomFaction()
	return b.WithFaction(faction)
}

// WithRandomTrait sets a random trait for the NPC.
func (b *NPCBuilder) WithTrait(trait string) *NPCBuilder {
	data := b.c.CreationDataService.GetTraitData(trait)
	b.npc.AddComponent(*data.NewTraitComponent())
	return b
}

// WithRandomTrait sets a random trait for the NPC.
func (b *NPCBuilder) WithRandomTrait() *NPCBuilder {
	trait := b.c.RandomizerService.RandomTrait()
	return b.WithTrait(trait)
}

// Build constructs the NPC.
func (b *NPCBuilder) Build() m.NPC {
	return *b.npc
}
