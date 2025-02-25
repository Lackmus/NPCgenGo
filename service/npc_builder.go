package service

import (
	"errors"
	"fmt"

	m "github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components" // assumed package for component types/keys
)

const (
	DEFAULT_NPC_TYPE = "default"
)

// NPCBuilder constructs an NPC step by step.
// It holds an internal error field that accumulates errors encountered during the build process.
type NPCBuilder struct {
	npc     *m.NPC
	c       *NPCCreationSupplier
	subtype *cp.NPCSubtype
	species *cp.Species
	npctype string
	err     error
}

// NewNPCBuilder creates a new NPCBuilder using the proper NPC constructor.
func NewNPCBuilder(c *NPCCreationSupplier) *NPCBuilder {
	return &NPCBuilder{
		npc:     m.NewNPC(c.RandomizerService.GenerateID()),
		c:       c,
		npctype: DEFAULT_NPC_TYPE,
	}
}

// WithNPC sets the NPC to the provided value and updates internal fields.
func (b *NPCBuilder) WithNPC(npc m.NPC) *NPCBuilder {
	if b.err != nil {
		return b
	}
	b.npc = &npc
	subtypeID := npc.GetComponent(cp.CompSubtype)
	speciesID := npc.GetComponent(cp.CompSpecies)
	if subtypeID == "" || speciesID == "" {
		b.err = errors.New("provided NPC is missing subtype or species components")
		return b
	}
	subtype := b.c.CreationDataService.GetNpcSubtypeData(subtypeID)
	species := b.c.CreationDataService.GetSpeciesData(speciesID)
	b.subtype = &subtype
	b.species = &species
	b.npctype = npc.GetComponent(cp.CompType)
	return b
}

// ----- Type Methods -----

// WithType sets the NPC's type to the provided value.
func (b *NPCBuilder) WithType(npctype string) *NPCBuilder {
	if b.err != nil {
		return b
	}
	data := b.c.CreationDataService.GetNpcTypeData(npctype)
	b.npctype = npctype
	b.npc.AddComponent(*data.NewNPCTypeComponent())
	return b
}

// WithRandomType sets the NPC's type by selecting a random type.
func (b *NPCBuilder) WithRandomType() *NPCBuilder {
	if b.err != nil {
		return b
	}
	randomType := b.c.RandomizerService.RandomType()
	return b.WithType(randomType)
}

// ----- Subtype Methods -----

// WithSubtype sets the NPC's subtype to the provided value.
func (b *NPCBuilder) WithSubtype(subtype string) *NPCBuilder {
	if b.err != nil {
		return b
	}
	data := b.c.CreationDataService.GetNpcSubtypeData(subtype)
	b.subtype = &data
	b.npc.AddComponent(*b.subtype.NewNPCSubtypeComponent())
	return b
}

// WithRandomSubtype sets the NPC's subtype by selecting a random subtype.
// It requires that the NPC type is already set.
func (b *NPCBuilder) WithRandomSubtype() *NPCBuilder {
	if b.err != nil {
		return b
	}
	if b.npctype == DEFAULT_NPC_TYPE {
		b.err = errors.New("type must be set before subtype can be added")
		return b
	}
	randomSubtype := b.c.RandomizerService.RandomSubtype(b.npctype)
	return b.WithSubtype(randomSubtype)
}

// WithSubtypeStats sets the NPC's subtype stats from a provided string.
func (b *NPCBuilder) WithSubtypeStats(stats string) *NPCBuilder {
	if b.err != nil {
		return b
	}
	b.npc.AddComponent(cp.Component{Name: cp.CompStats, Value: stats})
	return b
}

// WithRandomSubtypeStats sets the NPC's subtype stats using a random generator.
// Requires that a subtype has been set.
func (b *NPCBuilder) WithRandomSubtypeStats() *NPCBuilder {
	if b.err != nil {
		return b
	}
	if b.subtype == nil {
		b.err = errors.New("subtype must be set before stats can be added")
		return b
	}
	b.npc.AddComponent(*b.subtype.NewNPCSubtypeStatsComponent())
	return b
}

// WithSubtypeEquipment sets the NPC's equipment with the provided string.
func (b *NPCBuilder) WithSubtypeEquipment(items string) *NPCBuilder {
	if b.err != nil {
		return b
	}
	b.npc.AddComponent(cp.Component{Name: cp.CompItems, Value: items})
	return b
}

// WithRandomSubtypeEquipment sets the NPC's equipment using random generation.
// Requires that a subtype has been set.
func (b *NPCBuilder) WithRandomSubtypeEquipment() *NPCBuilder {
	if b.err != nil {
		return b
	}
	if b.subtype == nil {
		b.err = errors.New("subtype must be set before equipment can be added")
		return b
	}
	b.npc.AddComponent(*b.subtype.NewNPCSubtypeEquipmentComponent())
	return b
}

// ----- Species and Name Methods -----

// WithSpecies sets the NPC's species to the provided value.
func (b *NPCBuilder) WithSpecies(species string) *NPCBuilder {
	if b.err != nil {
		return b
	}
	data := b.c.CreationDataService.GetSpeciesData(species)
	b.species = &data
	b.npc.AddComponent(*b.species.NewSpeciesComponent())
	return b
}

// WithRandomSpecies sets the NPC's species by selecting a random species.
func (b *NPCBuilder) WithRandomSpecies() *NPCBuilder {
	if b.err != nil {
		return b
	}
	randomSpecies := b.c.RandomizerService.RandomSpecies()
	return b.WithSpecies(randomSpecies)
}

// WithName sets the NPC's name to the provided value.
func (b *NPCBuilder) WithName(name string) *NPCBuilder {
	if b.err != nil {
		return b
	}
	b.npc.AddComponent(cp.Component{Name: cp.CompName, Value: name})
	return b
}

// WithRandomName sets the NPC's name using random generation based on the current species.
// Requires that a species has been set.
func (b *NPCBuilder) WithRandomName() *NPCBuilder {
	if b.err != nil {
		return b
	}
	if b.species == nil {
		b.err = errors.New("species must be set before name can be added")
		return b
	}
	data := b.c.CreationDataService.GetNameData(b.species.NameSource)
	b.npc.AddComponent(*data.NewNameComponent())
	return b
}

// ----- Faction and Trait Methods -----

// WithFaction sets the NPC's faction to the provided value.
func (b *NPCBuilder) WithFaction(faction string) *NPCBuilder {
	if b.err != nil {
		return b
	}
	data := b.c.CreationDataService.GetFactionData(faction)
	b.npc.AddComponent(*data.NewFactionComponent())
	return b
}

// WithRandomFaction sets the NPC's faction by selecting a random faction.
func (b *NPCBuilder) WithRandomFaction() *NPCBuilder {
	if b.err != nil {
		return b
	}
	randomFaction := b.c.RandomizerService.RandomFaction()
	return b.WithFaction(randomFaction)
}

// WithTrait sets the NPC's trait to the provided value.
func (b *NPCBuilder) WithTrait(trait string) *NPCBuilder {
	if b.err != nil {
		return b
	}
	data := b.c.CreationDataService.GetTraitData(trait)
	b.npc.AddComponent(*data.NewTraitComponent())
	return b
}

// WithRandomTrait sets the NPC's trait by selecting a random trait.
func (b *NPCBuilder) WithRandomTrait() *NPCBuilder {
	if b.err != nil {
		return b
	}
	randomTrait := b.c.RandomizerService.RandomTrait()
	return b.WithTrait(randomTrait)
}

// ----- Build Method -----

// Build finalizes and returns the constructed NPC. If any error was encountered during the build process,
// it returns the error.
func (b *NPCBuilder) Build() (m.NPC, error) {
	if b.err != nil {
		return m.NPC{}, fmt.Errorf("failed to build NPC: %w", b.err)
	}
	return *b.npc, nil
}
