package service

import (
	"github.com/lackmus/npcgengo/model"
)

// npcBuilder is an internal, mutable structure used to collect values
// before constructing the immutable model.NPC.
type npcBuilder struct {
	id          string
	name        string
	faction     string
	species     string
	npcType     string
	npcSubtype  string
	trait       string
	drive       string
	stats       map[string]int
	items       map[string]string
	abilities   map[string]string
	description string
}

// NPCOption defines a function that configures the npcBuilder.
type NPCOption func(b *npcBuilder, rand RandomizerService)

// NewNPC creates a new immutable NPC by applying the provided options to a builder
// and then calling model.NewNPC to build the final object.
func NewNPC(rand RandomizerService, opts ...NPCOption) model.NPC {
	// Initialize the builder with sensible defaults (e.g. empty maps).
	b := &npcBuilder{
		stats:     make(map[string]int),
		items:     make(map[string]string),
		abilities: make(map[string]string),
	}
	// Apply all options.
	for _, opt := range opts {
		opt(b, rand)
	}
	// Now call the constructor from your model package.
	// (Assuming model.NewNPC has the signature:
	//  func NewNPC(id, name, faction, species, npcType, npcSubtype, trait, drive, description string,
	//             stats map[string]int, items map[string]string, abilities map[string]string) (NPC, error))
	npc := model.NewNPC(
		b.id,
		b.name,
		b.faction,
		b.species,
		b.npcType,
		b.npcSubtype,
		b.trait,
		b.drive,
		b.description,
		b.stats,
		b.items,
		b.abilities,
	)
	return npc
}

// WithID sets the NPC's ID; if an empty string is provided, a random ID is generated.
func WithID(id string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if id == "" {
			b.id = rand.GenerateID()
		} else {
			b.id = id
		}
	}
}

// WithType sets the NPC's type; if empty, a random type is chosen.
func WithType(npcType string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if npcType == "" {
			b.npcType = rand.RandomType()
		} else {
			b.npcType = npcType
		}
	}
}

// WithSubType sets the NPC's subtype. Note: It uses the already set npcType.
func WithSubType(npcSubType string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if npcSubType == "" {
			b.npcSubtype = rand.RandomSubtype(b.npcType)
		} else {
			b.npcSubtype = npcSubType
		}
	}
}

// WithFaction sets the NPC's faction; if empty, a random faction is chosen.
func WithFaction(faction string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if faction == "" {
			b.faction = rand.RandomFaction()
		} else {
			b.faction = faction
		}
	}
}

// WithSpecies sets the NPC's species; if empty, a random species is chosen.
func WithSpecies(species string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if species == "" {
			b.species = rand.RandomSpecies()
		} else {
			b.species = species
		}
	}
}

// WithName sets the NPC's name; if empty, a name is generated based on the species.
func WithName(name string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if name == "" {
			b.name = rand.GenerateName(b.species)
		} else {
			b.name = name
		}
	}
}

// WithTrait sets the NPC's trait; if empty, a random trait is generated.
func WithTrait(trait string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if trait == "" {
			b.trait = rand.GenerateTraitDescription()
		} else {
			b.trait = trait
		}
	}
}

// WithAbilities assigns a map of abilities to the NPC.
func WithAbilities(abilities map[string]string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if abilities != nil {
			b.abilities = abilities
		}
	}
}

// WithItems sets the NPC's items; if nil, items are generated based on the NPC's subtype.
func WithItems(items map[string]string) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if items == nil {
			b.items = rand.GenerateEquipment(b.npcSubtype)
		} else {
			b.items = items
		}
	}
}

// WithStats sets the NPC's stats; if nil, stats are generated based on the NPC's subtype.
func WithStats(stats map[string]int) NPCOption {
	return func(b *npcBuilder, rand RandomizerService) {
		if stats == nil {
			b.stats = rand.ApplySubtypeStats(b.npcSubtype)
		} else {
			b.stats = stats
		}
	}
}
