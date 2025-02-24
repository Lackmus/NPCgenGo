package shared

import (
	c "github.com/lackmus/npcgengo/model/npc_components"
)

type NPCConfigLoader interface {
	LoadFactionMap() (map[string]c.Faction, error)
	LoadSpeciesMap() (map[string]c.Species, error)
	LoadTraitMap() (map[string]c.Trait, error)
	LoadNameMap() (map[string]c.NameData, error)
	LoadNpcCivilianSubtypeMap() (map[string]c.NPCSubtype, error)
	LoadNpcMilitarySubtypeMap() (map[string]c.NPCSubtype, error)
}
