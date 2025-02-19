package service

import (
	"github.com/lackmus/npcgengo/model"
)

// NPCOption defines a function that configures an NPC.
// The RandomizerService is passed along so options can generate defaults.
type NPCOption func(npc *model.NPC, rand RandomizerService)

// NewNPC creates a new NPC by applying the provided options.
// It uses the given RandomizerService to supply default values where needed.
func NewNPC(rand RandomizerService, opts ...NPCOption) model.NPC {
	npc := model.NPC{}
	for _, opt := range opts {
		opt(&npc, rand)
	}
	return npc
}

// WithID sets the NPC's ID; if an empty string is provided, a random ID is generated.
func WithID(id string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if id == "" {
			npc.ID = rand.GenerateID()
		} else {
			npc.ID = id
		}
	}
}

// WithType sets the NPC's type; if empty, a random type is chosen.
func WithType(npcType string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if npcType == "" {
			npc.NpcType = rand.RandomType()
		} else {
			npc.NpcType = npcType
		}
	}
}

// WithSubType sets the NPC's subtype. Note: It uses the already set NpcType.
// (So ensure WithType is applied before WithSubType.)
func WithSubType(npcSubType string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if npcSubType == "" {
			npc.NpcSubtype = rand.RandomSubtype(npc.NpcType)
		} else {
			npc.NpcSubtype = npcSubType
		}
	}
}

// WithFaction sets the NPC's faction; if empty, a random faction is chosen.
func WithFaction(faction string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if faction == "" {
			npc.Faction = rand.RandomFaction()
		} else {
			npc.Faction = faction
		}
	}
}

// WithSpecies sets the NPC's species; if empty, a random species is chosen.
func WithSpecies(species string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if species == "" {
			npc.Species = rand.RandomSpecies()
		} else {
			npc.Species = species
		}
	}
}

// WithName sets the NPC's name; if empty, a name is generated based on the species.
func WithName(name string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if name == "" {
			npc.Name = rand.GenerateName(npc.Species)
		} else {
			npc.Name = name
		}
	}
}

// WithTrait sets the NPC's trait; if empty, a random trait is generated.
func WithTrait(trait string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if trait == "" {
			npc.Trait = rand.GenerateTraitDescription()
		} else {
			npc.Trait = trait
		}
	}
}

// WithAbilities assigns a map of abilities to the NPC.
// If nil is provided, the NPC.Abilities remains as the empty map initialized in NewNPC.
func WithAbilities(abilities map[string]string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if abilities != nil {
			npc.Abilities = abilities
		}
	}
}

// WithItems sets the NPC's items; if nil, items are generated based on the NPC's subtype.
func WithItems(items map[string]string) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if items == nil {
			npc.Items = rand.GenerateEquipment(npc.NpcSubtype)
		} else {
			npc.Items = items
		}
	}
}

// WithStats sets the NPC's stats; if nil, stats are generated based on the NPC's subtype.
func WithStats(stats map[string]int) NPCOption {
	return func(npc *model.NPC, rand RandomizerService) {
		if stats == nil {
			npc.Stats = rand.ApplySubtypeStats(npc.NpcSubtype)
		} else {
			npc.Stats = stats
		}
	}
}
