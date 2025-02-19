package shared

import (
	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/model/types"
)

type NPCConfigLoader interface {
	LoadFactionMap() (map[string]model.Faction, error)
	LoadSpeciesMap() (map[string]model.Species, error)
	LoadTraitMap() (map[string]model.Trait, error)
	LoadNameMap() (map[string]model.NameData, error)
	LoadNpcCivilianSubtypeMap() (map[string]types.NPCSubtype, error)
	LoadNpcMilitarySubtypeMap() (map[string]types.NPCSubtype, error)
}
