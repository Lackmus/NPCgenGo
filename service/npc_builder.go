// Description: This file contains the NPCBuilder struct and methods for constructing an NPC step by step.
package service

import (
	"errors"
	"fmt"
	"strings"

	h "github.com/lackmus/npcgengo/helper"
	m "github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components" // assumed package for component types/keys
	t "github.com/lackmus/npcgengo/model/npc_components/types"
)

// NPCBuilder constructs an NPC step by step.
// It holds an internal error field that accumulates errors encountered during the build process.
type NPCBuilder struct {
	npc         *m.NPC
	supplier    *NPCCreationSupplier
	subtypeData *cp.NPCSubtype
	speciesData *cp.Species
	traitData   *cp.Trait
	npctypeData *t.NPCType
	description string
	err         error
}

// NewNPCBuilder creates a new NPCBuilder using the proper NPC constructor.
func NewNPCBuilder(supplier *NPCCreationSupplier) *NPCBuilder {
	return &NPCBuilder{
		npc:      m.NewNPC(),
		supplier: supplier,
		err:      nil,
	}
}

// WithNPC sets the NPC to the provided value and updates internal fields.
// It requires that the NPC has components for type, subtype, species, and trait.
func (b *NPCBuilder) WithNPC(npc m.NPC) *NPCBuilder {
	if b.checkError() {
		return b
	}
	b.npc = &npc
	b.fetchAndSetComponents(npc)
	return b
}

// fetchAndSetComponents retrieves the NPC's components from the data service and sets them internally.
// It requires that the NPC has components for type, subtype, species, and trait.
func (b *NPCBuilder) fetchAndSetComponents(npc m.NPC) {
	npctypeID := npc.GetComponent(cp.CompType)
	subtypeID := npc.GetComponent(cp.CompSubtype)
	speciesID := npc.GetComponent(cp.CompSpecies)
	traitID := npc.GetComponent(cp.CompTrait)

	b.npctypeData = ptr(b.supplier.CreationDataService.GetNpcTypeData(npctypeID))
	b.subtypeData = ptr(b.supplier.CreationDataService.GetNpcSubtypeData(subtypeID))
	b.speciesData = ptr(b.supplier.CreationDataService.GetSpeciesData(speciesID))
	b.traitData = ptr(b.supplier.CreationDataService.GetTraitData(traitID))
}

// ptr is a helper function that returns a pointer to the provided value.
func ptr[T any](t T) *T { return &t }

// ----- Type Methods -----

// WithType sets the NPC's type to the provided value.
// It requires that the type exists in the data service.
func (b *NPCBuilder) WithType(npctype string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	data := b.supplier.CreationDataService.GetNpcTypeData(npctype)
	b.npctypeData = &data
	b.npc.AddComponent(cp.NewComponent(cp.CompType, npctype))
	return b
}

// WithRandomType sets the NPC's type by selecting a random type.
// It requires that the NPC type is already set.
func (b *NPCBuilder) WithRandomType() *NPCBuilder {
	if b.checkError() {
		return b
	}
	randomType := b.supplier.RandomizerService.RandomType()
	return b.WithType(randomType)
}

// ----- Subtype Methods -----

// WithSubtype sets the NPC's subtype to the provided value.
// It requires that the subtype exists in the data service.
func (b *NPCBuilder) WithSubtype(subtype string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	data := b.supplier.CreationDataService.GetNpcSubtypeData(subtype)
	b.subtypeData = &data
	b.npc.AddComponent(cp.NewComponent(cp.CompSubtype, subtype))
	return b
}

// WithRandomSubtype sets the NPC's subtype by selecting a random subtype.
// It requires that the NPC type is already set.
func (b *NPCBuilder) WithRandomSubtype() *NPCBuilder {
	if b.checkError() {
		return b
	}
	if h.IsNilOrEmpty(b.npctypeData) {
		b.err = errors.New("type must be set before subtype can be added")
		return b
	}
	randomSubtype := b.supplier.RandomizerService.RandomSubtype(b.npctypeData.Name)
	return b.WithSubtype(randomSubtype)
}

// WithSubtypeStats sets the NPC's subtype stats from a provided string.
// It requires that the stats are not empty.
func (b *NPCBuilder) WithSubtypeStats(stats string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	b.npc.AddComponent(cp.NewComponent(cp.CompStats, stats))
	return b
}

// WithRandomSubtypeStats sets the NPC's subtype stats using a random generator.
// It requires that a subtype has been set.
func (b *NPCBuilder) WithRandomSubtypeStats() *NPCBuilder {
	if b.checkError() {
		return b
	}
	if h.IsNilOrEmpty(b.subtypeData) {
		b.err = errors.New("subtype must be set before stats can be added")
		return b
	}
	b.npc.AddComponent(cp.NewComponent(cp.CompStats, b.subtypeData.GetStats()))
	return b
}

// WithSubtypeEquipment sets the NPC's equipment with the provided string.
// It requires that the equipment is not empty.
func (b *NPCBuilder) WithSubtypeEquipment(equipment string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	b.npc.AddComponent(cp.NewComponent(cp.CompItems, equipment))
	return b
}

// WithRandomSubtypeEquipment sets the NPC's equipment using random generation.
// It requires that a subtype has been set.
func (b *NPCBuilder) WithRandomSubtypeEquipment() *NPCBuilder {
	if b.checkError() {
		return b
	}
	if h.IsNilOrEmpty(b.subtypeData) {
		b.err = errors.New("subtype must be set before equipment can be added")
		return b
	}
	b.npc.AddComponent(cp.NewComponent(cp.CompItems, b.subtypeData.GetEquipment()))
	return b
}

// ----- Species and Name Methods -----

// WithSpecies sets the NPC's species to the provided value.
// It requires that the species exists in the data service.
func (b *NPCBuilder) WithSpecies(species string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	data := b.supplier.CreationDataService.GetSpeciesData(species)
	b.speciesData = &data
	b.npc.AddComponent(cp.NewComponent(cp.CompSpecies, species))
	return b
}

// WithRandomSpecies sets the NPC's species by selecting a random species.
// It requires that the NPC type is already set.
func (b *NPCBuilder) WithRandomSpecies() *NPCBuilder {
	if b.checkError() {
		return b
	}
	randomSpecies := b.supplier.RandomizerService.RandomSpecies()
	return b.WithSpecies(randomSpecies)
}

// WithName sets the NPC's name to the provided value.
// It requires that the name is not empty.
func (b *NPCBuilder) WithName(name string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	b.npc.AddComponent(cp.NewComponent(cp.CompName, name))
	return b
}

// WithRandomName sets the NPC's name using random generation based on the current species.
// Requires that a species has been set.
func (b *NPCBuilder) WithRandomName() *NPCBuilder {
	if b.checkError() {
		return b
	}
	if h.IsNilOrEmpty(b.speciesData) {
		b.err = errors.New("species must be set before name can be added")
		return b
	}
	data := b.supplier.CreationDataService.GetNameData(b.speciesData.NameSource)
	b.npc.AddComponent(cp.NewComponent(cp.CompName, data.GenerateName()))
	return b
}

// ----- Faction and Trait Methods -----

// WithFaction sets the NPC's faction to the provided value.
// It requires that the faction exists in the data service.
func (b *NPCBuilder) WithFaction(faction string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	b.npc.AddComponent(cp.NewComponent(cp.CompFaction, faction))
	return b
}

// WithRandomFaction sets the NPC's faction by selecting a random faction.
// It requires that the NPC type is already set.
func (b *NPCBuilder) WithRandomFaction() *NPCBuilder {
	if b.checkError() {
		return b
	}
	randomFaction := b.supplier.RandomizerService.RandomFaction()
	return b.WithFaction(randomFaction)
}

// WithTrait sets the NPC's trait to the provided value.
// It requires that the trait exists in the data service.
func (b *NPCBuilder) WithTrait(trait string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	data := b.supplier.CreationDataService.GetTraitData(trait)
	b.traitData = &data
	b.npc.AddComponent(cp.NewComponent(cp.CompTrait, trait))
	return b
}

// WithRandomTrait sets the NPC's trait by selecting a random trait.
// It requires that the NPC type is already set.
func (b *NPCBuilder) WithRandomTrait() *NPCBuilder {
	if b.checkError() {
		return b
	}
	randomTrait := b.supplier.RandomizerService.RandomTrait()
	return b.WithTrait(randomTrait)
}

// ----- Description Methods -----

// WithDescription sets the NPC's description to the provided value.
// It requires that the description is not empty.
func (b *NPCBuilder) WithDescription(description string) *NPCBuilder {
	if b.checkError() {
		return b
	}
	// Wrap the current description component with additional info.
	b.description = description
	b.npc.AddComponent(cp.NewComponent(cp.CompDescription, description))
	return b
}

// formatDescription combines the subtype, species, and trait into a single string.
// It returns a formatted string with each part separated by a newline.
func formatDescription(subtype, species, trait string) string {
	var sb strings.Builder
	sb.WriteString("\n")
	if trimmed := strings.TrimSpace(subtype); trimmed != "" {
		sb.WriteString("  - " + trimmed + "\n")
	}
	if trimmed := strings.TrimSpace(species); trimmed != "" {
		sb.WriteString("  - " + trimmed + "\n")
	}
	if trimmed := strings.TrimSpace(trait); trimmed != "" {
		sb.WriteString("  - " + trimmed + "\n")
	}
	return sb.String()
}

// WithRandomDescription sets the NPC's description using the subtype, species, and trait.
// It requires that the subtype, species, and trait have been set.
func (b *NPCBuilder) WithRandomDescription() *NPCBuilder {
	if b.checkError() {
		return b
	}
	desc := formatDescription(b.subtypeData.GetDescription(), b.speciesData.GetDescription(), b.traitData.GetDescription())
	b.description = desc
	if !h.IsNilOrEmpty(desc) {
		b.npc.AddComponent(cp.NewComponent(cp.CompDescription, desc))
	}
	return b
}

// ----- Build Method -----

// Validate checks if the NPCBuilder is in a valid state to build an NPC.
// It returns an error if the NPC type, subtype, species, or trait are not set.
func (b *NPCBuilder) Validate() error {
	if h.IsNilOrEmpty(b.npctypeData) {
		return errors.New("NPC type is not set")
	}
	if h.IsNilOrEmpty(b.subtypeData) {
		return errors.New("NPC subtype is not set")
	}
	if h.IsNilOrEmpty(b.speciesData) {
		return errors.New("NPC species is not set")
	}
	if h.IsNilOrEmpty(b.traitData) {
		return errors.New("NPC trait is not set")
	}
	// Add further validations as needed.
	return nil
}

func (b *NPCBuilder) checkError() bool {
	return b.err != nil
}

// Build constructs the NPC from the current state of the builder.
// It returns an error if the builder is in an invalid state.
func (b *NPCBuilder) Build() (m.NPC, error) {
	if b.checkError() {
		return m.NPC{}, fmt.Errorf("failed to build NPC: %w", b.err)
	}
	if err := b.Validate(); err != nil {
		return m.NPC{}, err
	}
	return *b.npc, nil
}
