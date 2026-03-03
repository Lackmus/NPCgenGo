package shared

import (
	"context"

	c "github.com/lackmus/npcgengo/pkg/model/npc_components"
)

type NPCConfigLoader interface {

	// LoadFactionMap loads the faction map from the config file.
	// It returns a map of faction names to Faction objects and an error if the data cannot be loaded.
	LoadFactionMap(ctx context.Context) (map[string]c.Faction, error)

	// LoadSpeciesMap loads the species map from the config file.
	// It returns a map of species names to Species objects and an error if the data cannot be loaded.
	LoadSpeciesMap(ctx context.Context) (map[string]c.Species, error)

	// LoadTraitMap loads the trait map from the config file.
	// It returns a map of trait names to Trait objects and an error if the data cannot be loaded.
	LoadTraitMap(ctx context.Context) (map[string]c.Trait, error)

	// LoadNameMap loads the name map from the config file.
	// It returns a map of name keys to NameData values and an error if the data cannot be loaded.
	LoadNameMap(ctx context.Context) (map[string]c.NameData, error)

	// LoadNPCSubtypeMaps loads subtype maps grouped by NPC type.
	// The returned map key is NPC type name and each value is a subtype map keyed by subtype name.
	LoadNPCSubtypeMaps(ctx context.Context) (map[string]map[string]c.NPCSubtype, error)
}
