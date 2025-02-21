package service

import (
	"fmt"
	"maps"
	"slices"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/model/types"
	"github.com/lackmus/npcgengo/shared"
)

// ==========================================
// CreationDataService
// ==========================================

// CreationDataService provides access to the creation data for NPCs. It loads the data from the creation data file into maps.
type CreationDataService struct {
	factionMap                                            map[string]model.Faction
	speciesMap                                            map[string]model.Species
	traitMap                                              map[string]model.Trait
	nameMap                                               map[string]model.NameData
	npcTypeMap                                            map[string]types.NPCType
	npcSubtypeMap, civilianSubtypeMap, militarySubtypeMap map[string]types.NPCSubtype
	speciesNameMap                                        map[string]string
	npcSubtypeForTypeMap                                  map[string][]string
}

// NewCreationDataService creates a new CreationDataService. It loads the data from the creation data file into maps.
func NewCreationDataService(npcConfigLoader shared.NPCConfigLoader) (*CreationDataService, error) {
	factionMap, err := npcConfigLoader.LoadFactionMap()
	if err != nil {
		return nil, fmt.Errorf("failed to load faction map: %w", err)
	}
	speciesMap, err := npcConfigLoader.LoadSpeciesMap()
	if err != nil {
		return nil, fmt.Errorf("failed to load species map: %w", err)
	}
	traitMap, err := npcConfigLoader.LoadTraitMap()
	if err != nil {
		return nil, fmt.Errorf("failed to load trait map: %w", err)
	}
	nameMap, err := npcConfigLoader.LoadNameMap()
	if err != nil {
		return nil, fmt.Errorf("failed to load name map: %w", err)
	}
	civilianSubtypeMap, err := npcConfigLoader.LoadNpcCivilianSubtypeMap()
	if err != nil {
		return nil, fmt.Errorf("failed to load civilian subtype map: %w", err)
	}
	militarySubtypeMap, err := npcConfigLoader.LoadNpcMilitarySubtypeMap()
	if err != nil {
		return nil, fmt.Errorf("failed to load military subtype map: %w", err)
	}

	cds := &CreationDataService{
		factionMap:         factionMap,
		speciesMap:         speciesMap,
		traitMap:           traitMap,
		nameMap:            nameMap,
		civilianSubtypeMap: civilianSubtypeMap,
		militarySubtypeMap: militarySubtypeMap,
	}
	cds.npcTypeMap = cds.loadNpcTypeMap()
	cds.npcSubtypeMap = cds.mergeNpcSubtypeMaps()
	cds.speciesNameMap = cds.buildSpeciesNameMap()
	cds.npcSubtypeForTypeMap = cds.buildNpcTypeNameMap()
	return cds, nil
}

func (c *CreationDataService) buildSpeciesNameMap() map[string]string {
	snm := make(map[string]string)
	for key, species := range c.speciesMap {
		if nameData, ok := c.nameMap[species.NameSource]; ok {
			snm[key] = nameData.Name
		}
	}
	return snm
}

func (c *CreationDataService) loadNpcTypeMap() map[string]types.NPCType {
	return map[string]types.NPCType{
		"Civilian": types.GetCivilianInstance().NPCType, // dereference the embedded field
		"Military": types.GetMilitaryInstance().NPCType,
	}
}

func (c *CreationDataService) buildNpcTypeNameMap() map[string][]string {
	return map[string][]string{
		"Civilian": slices.Collect(maps.Keys(c.civilianSubtypeMap)),
		"Military": slices.Collect(maps.Keys(c.militarySubtypeMap)),
	}
}

func (c *CreationDataService) mergeNpcSubtypeMaps() map[string]types.NPCSubtype {
	merged := maps.Clone(c.civilianSubtypeMap)
	for key, subtype := range c.militarySubtypeMap {
		merged[key] = subtype
	}
	return merged
}

func (c *CreationDataService) GetFactionData(key string) model.Faction {
	faction, ok := c.factionMap[key]
	if !ok {
		panic(fmt.Sprintf("faction not found: %s", key))
	}
	return faction
}

func (c *CreationDataService) GetTraitData(key string) model.Trait {
	trait, ok := c.traitMap[key]
	if !ok {
		panic(fmt.Sprintf("trait not found: %s", key))
	}
	return trait
}

func (c *CreationDataService) GetNameData(key string) model.NameData {
	nd, ok := c.nameMap[key]
	if ok {
		return nd
	}
	panic(fmt.Sprintf("name not found: %s", key))
}

func (c *CreationDataService) GetSpeciesData(key string) model.Species {
	s, ok := c.speciesMap[key]
	if ok {
		return s
	}
	panic(fmt.Sprintf("species not found: %s", key))
}

func (c *CreationDataService) GetNpcTypeData(key string) types.NPCType {
	nt, ok := c.npcTypeMap[key]
	if ok {
		return nt
	}
	panic(fmt.Sprintf("npc type not found: %s", key))
}

func (c *CreationDataService) GetNpcSubtypeData(key string) types.NPCSubtype {
	ns, ok := c.npcSubtypeMap[key]
	if ok {
		return ns
	}
	panic(fmt.Sprintf("npc subtype not found: %s", key))
}

func (c *CreationDataService) GetFactionMap() map[string]model.Faction {
	return maps.Clone(c.factionMap)
}

func (c *CreationDataService) GetSpeciesMap() map[string]model.Species {
	return maps.Clone(c.speciesMap)
}

func (c *CreationDataService) GetNpcSubtypeForTypeMap() map[string][]string {
	return maps.Clone(c.npcSubtypeForTypeMap)
}

func (c *CreationDataService) GetTraitMap() map[string]model.Trait {
	return maps.Clone(c.traitMap)
}

func (c *CreationDataService) GetNpcTypeMap() map[string]types.NPCType {
	return maps.Clone(c.npcTypeMap)
}

func (c *CreationDataService) GetSpeciesNameMap() map[string]string {
	return maps.Clone(c.speciesNameMap)
}
