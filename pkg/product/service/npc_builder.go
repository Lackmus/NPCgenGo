// Description: This file contains the NPCBuilder struct and methods for constructing an NPC step by step.
package service

import (
	"errors"
	"fmt"
	"strings"

	h "github.com/lackmus/npcgengo/internal/platform/helpers"
	m "github.com/lackmus/npcgengo/pkg/product/model"
	cp "github.com/lackmus/npcgengo/pkg/product/model/npc_components" // assumed package for component types/keys
	t "github.com/lackmus/npcgengo/pkg/product/model/npc_components/types"
)

// NPCBuilder constructs an NPC step by step.
// It holds an internal error field that accumulates errors encountered during the build process.
type NPCBuilder struct {
	npc         *m.NPC
	supplier    *NPCCreationSupplier
	idData      *string
	subtypeData *cp.NPCSubtype
	speciesData *cp.Species
	traitData   *cp.Trait
	npctypeData *t.NPCType
	errors      []error
}

// NewNPCBuilder creates a new NPCBuilder using the proper NPC constructor.
func NewNPCBuilder(supplier *NPCCreationSupplier) *NPCBuilder {
	npc := m.NewNPC()
	return &NPCBuilder{
		npc:      npc,
		supplier: supplier,
		errors:   make([]error, 0),
	}
}

// getnpcTypeData returns the NPC type data for the current NPC.
func (b *NPCBuilder) GetNPCType() string {
	return b.npctypeData.Name
}

func (b *NPCBuilder) updateComponent(compType cp.CompEnum, value string) {
	b.npc.AddComponent(cp.NewComponent(compType, value))
}

// ----- NPC Methods -----

// WithNPC sets the NPC to the provided value and updates internal fields.
// It requires that the NPC has components for type, subtype, species, and trait.
func (b *NPCBuilder) WithNPC(npc m.NPC) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	b.npc = &npc
	b.fetchAndSetComponents(npc)
	return b
}

// GetNPC returns the NPC that is being built.
func (b *NPCBuilder) GetNPC() *m.NPC {
	return b.npc
}

// fetchAndSetComponents retrieves the NPC's components from the data service and sets them internally.
// It requires that the NPC has components for type, subtype, species, and trait.
func (b *NPCBuilder) fetchAndSetComponents(npc m.NPC) {
	npctypeID := npc.GetComponent(cp.CompType)
	subtypeID := npc.GetComponent(cp.CompSubtype)
	speciesID := npc.GetComponent(cp.CompSpecies)
	traitID := npc.GetComponent(cp.CompTrait)

	if data, err := b.supplier.CreationDataService.GetNpcTypeData(npctypeID); err == nil {
		b.npctypeData = ptr(data)
	} else {
		b.addErrorWithContext("fetchAndSetComponents", err)
	}

	if data, err := b.supplier.CreationDataService.GetNpcSubtypeData(subtypeID); err == nil {
		b.subtypeData = ptr(data)
	} else {
		b.addErrorWithContext("fetchAndSetComponents", err)
	}

	if data, err := b.supplier.CreationDataService.GetSpeciesData(speciesID); err == nil {
		b.speciesData = ptr(data)
	} else {
		b.addErrorWithContext("fetchAndSetComponents", err)
	}

	if data, err := b.supplier.CreationDataService.GetTraitData(traitID); err == nil {
		b.traitData = ptr(data)
	} else {
		b.addErrorWithContext("fetchAndSetComponents", err)
	}
}

// ----- ID Methods -----

func (b *NPCBuilder) WithID(d string) *NPCBuilder {
	if value := strings.TrimSpace(d); value != "" {
		b.npc.ID = value
		b.idData = &value
	}
	return b
}

// ptr is a helper function that returns a pointer to the provided value.
func ptr[T any](t T) *T { return &t }

// ----- Type Methods -----

// WithType sets the NPC's type to the provided value.
// It requires that the type exists in the data service.
func (b *NPCBuilder) WithType(npctype string) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	data, err := b.supplier.CreationDataService.GetNpcTypeData(npctype)
	if err != nil {
		b.addErrorWithContext("WithType", err)
		return b
	}
	b.npctypeData = &data
	b.updateComponent(cp.CompType, npctype)
	return b
}

// WithRandomType sets the NPC's type by selecting a random type.
// It requires that the NPC type is already set.
func (b *NPCBuilder) WithRandomType() *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	randomType := b.supplier.RandomizerService.RandomType()
	return b.WithType(randomType)
}

// ----- Subtype Methods -----

// WithSubtype sets the NPC's subtype to the provided value.
// It requires that the subtype exists in the data service.
func (b *NPCBuilder) WithSubtype(subtype string) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	data, err := b.supplier.CreationDataService.GetNpcSubtypeData(subtype)
	if err != nil {
		b.addErrorWithContext("WithSubtype", err)
		return b
	}
	b.subtypeData = &data
	b.updateComponent(cp.CompSubtype, subtype)
	return b
}

// WithRandomSubtype sets the NPC's subtype by selecting a random subtype.
// It requires that the NPC type is already set.
func (b *NPCBuilder) WithRandomSubtype() *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	if h.IsNilOrEmpty(b.npctypeData) {
		b.addErrorWithContext("WithRandomSubtype", errors.New("type must be set before subtype can be added"))
		return b
	}
	randomSubtype := b.supplier.RandomizerService.RandomSubtype(b.npctypeData.Name)
	return b.WithSubtype(randomSubtype)
}

// WithSubtypeStats sets the NPC's subtype stats from a provided string.
// It requires that the stats are not empty.
func (b *NPCBuilder) WithSubtypeStats(stats string) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	b.updateComponent(cp.CompStats, stats)
	return b
}

// WithRandomSubtypeStats sets the NPC's subtype stats using a random generator.
// It requires that a subtype has been set.
func (b *NPCBuilder) WithRandomSubtypeStats() *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	if h.IsNilOrEmpty(b.subtypeData) {
		b.addErrorWithContext("WithRandomSubtypeStats", errors.New("subtype must be set before stats can be added"))
		return b
	}
	b.updateComponent(cp.CompStats, b.subtypeData.GetStats())
	return b
}

// WithSubtypeEquipment sets the NPC's equipment with the provided string.
// It requires that the equipment is not empty.
func (b *NPCBuilder) WithSubtypeEquipment(equipment string) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	b.updateComponent(cp.CompItems, equipment)
	return b
}

// WithRandomSubtypeEquipment sets the NPC's equipment using random generation.
// It requires that a subtype has been set.
func (b *NPCBuilder) WithRandomSubtypeEquipment() *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	if h.IsNilOrEmpty(b.subtypeData) {
		b.addErrorWithContext("WithRandomSubtypeEquipment", errors.New("subtype must be set before equipment can be added"))
		return b
	}
	b.updateComponent(cp.CompItems, b.subtypeData.GetEquipment())
	return b
}

// ----- Species and Name Methods -----

func (b *NPCBuilder) WithSpecies(species string) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	data, err := b.supplier.CreationDataService.GetSpeciesData(species)
	if err != nil {
		b.addErrorWithContext("WithSpecies", err)
		return b
	}
	b.speciesData = &data
	b.updateComponent(cp.CompSpecies, species)
	return b
}

func (b *NPCBuilder) WithRandomSpecies() *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	randomSpecies := b.supplier.RandomizerService.RandomSpecies()
	return b.WithSpecies(randomSpecies)
}

func (b *NPCBuilder) WithName(name string) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	b.updateComponent(cp.CompName, name)
	return b
}

func (b *NPCBuilder) WithRandomName() *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	if h.IsNilOrEmpty(b.speciesData) {
		b.addErrorWithContext("WithRandomName", errors.New("species must be set before name can be added"))
		return b
	}
	data, err := b.supplier.CreationDataService.GetNameData(b.speciesData.NameSource)
	if err != nil {
		b.addErrorWithContext("WithRandomName", err)
		return b
	}
	b.WithName(data.GenerateName())
	return b
}

// ----- Faction and Trait Methods -----

func (b *NPCBuilder) WithFaction(faction string) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	b.updateComponent(cp.CompFaction, faction)
	return b
}

func (b *NPCBuilder) WithRandomFaction() *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	randomFaction := b.supplier.RandomizerService.RandomFaction()
	return b.WithFaction(randomFaction)
}

func (b *NPCBuilder) WithTrait(trait string) *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	data, err := b.supplier.CreationDataService.GetTraitData(trait)
	if err != nil {
		b.addErrorWithContext("WithTrait", err)
		return b
	}
	b.traitData = &data
	b.updateComponent(cp.CompTrait, trait)
	return b
}

func (b *NPCBuilder) WithRandomTrait() *NPCBuilder {
	if b.HasErrors() {
		return b
	}
	randomTrait := b.supplier.RandomizerService.RandomTrait()
	return b.WithTrait(randomTrait)
}

// ----- Build Method -----

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
	return nil
}

func (b *NPCBuilder) HasErrors() bool {
	return len(b.errors) > 0
}

func (b *NPCBuilder) Errors() []error {
	return b.errors
}

func (b *NPCBuilder) addErrorWithContext(context string, err error) {
	b.errors = append(b.errors, fmt.Errorf("%s: %w", context, err))
}

func (b *NPCBuilder) Build() (m.NPC, error) {
	if b.HasErrors() {
		return m.NPC{}, fmt.Errorf("cannot build NPC: %w", errors.Join(b.errors...))
	}
	if err := b.Validate(); err != nil {
		return m.NPC{}, err
	}
	return *b.npc, nil
}
