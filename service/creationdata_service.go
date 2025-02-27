// Description: This file contains the CreationDataService struct and its methods. The CreationDataService struct provides access to the creation data for NPCs. It is used to load the data from the creation data file into maps.
package service

import (
	"fmt"
	"maps"
	"slices"

	cp "github.com/lackmus/npcgengo/model/npc_components"
	t "github.com/lackmus/npcgengo/model/npc_components/types"
	"github.com/lackmus/npcgengo/shared"
)

// CreationDataService provides access to the creation data for NPCs.
// It is used to load the data from the creation data file into maps.
type CreationDataService struct {
	factionMap                                            map[string]cp.Faction
	speciesMap                                            map[string]cp.Species
	traitMap                                              map[string]cp.Trait
	nameMap                                               map[string]cp.NameData
	npcTypeMap                                            map[string]t.NPCType
	npcSubtypeMap, civilianSubtypeMap, militarySubtypeMap map[string]cp.NPCSubtype
	speciesNameMap                                        map[string]string
	npcSubtypeForTypeMap                                  map[string][]string
}

// NewCreationDataService creates a new CreationDataService. It loads the data from the creation data file into maps.
// It returns an error if the data cannot be loaded.
func NewCreationDataService(npcConfigLoader shared.NPCConfigLoader) (*CreationDataService, error) {
	cds := &CreationDataService{}
	cds.initConfigLoaderMaps(npcConfigLoader)
	cds.npcTypeMap = cds.loadNpcTypeMap()
	cds.npcSubtypeMap = cds.mergeNpcSubtypeMaps(cds.civilianSubtypeMap, cds.militarySubtypeMap)
	cds.speciesNameMap = cds.buildSpeciesNameMap()
	cds.npcSubtypeForTypeMap = cds.buildNpcTypeNameMap()
	return cds, nil
}

// initConfigLoaderMaps loads the data from the creation data file into maps.
// It panics if the data cannot be loaded.
func (c *CreationDataService) initConfigLoaderMaps(npcConfigLoader shared.NPCConfigLoader) {
	factionMap, err := npcConfigLoader.LoadFactionMap()
	if err != nil {
		panic(fmt.Errorf("failed to load faction map: %w", err))
	}
	speciesMap, err := npcConfigLoader.LoadSpeciesMap()
	if err != nil {
		panic(fmt.Errorf("failed to load species map: %w", err))
	}
	traitMap, err := npcConfigLoader.LoadTraitMap()
	if err != nil {
		panic(fmt.Errorf("failed to load trait map: %w", err))
	}
	nameMap, err := npcConfigLoader.LoadNameMap()
	if err != nil {
		panic(fmt.Errorf("failed to load name map: %w", err))
	}
	civilianSubtypeMap, err := npcConfigLoader.LoadNpcCivilianSubtypeMap()
	if err != nil {
		panic(fmt.Errorf("failed to load civilian subtype map: %w", err))
	}
	militarySubtypeMap, err := npcConfigLoader.LoadNpcMilitarySubtypeMap()
	if err != nil {
		panic(fmt.Errorf("failed to load military subtype map: %w", err))
	}

	c.factionMap = factionMap
	c.speciesMap = speciesMap
	c.traitMap = traitMap
	c.nameMap = nameMap
	c.civilianSubtypeMap = civilianSubtypeMap
	c.militarySubtypeMap = militarySubtypeMap
}

// buildSpeciesNameMap creates a map of species keys to species names.
// It uses the name data from the name map to get the species names.
func (c *CreationDataService) buildSpeciesNameMap() map[string]string {
	snm := make(map[string]string)
	for key, species := range c.speciesMap {
		if nameData, ok := c.nameMap[species.NameSource]; ok {
			snm[key] = nameData.GetName()
		}
	}
	return snm
}

// loadNpcTypeMap creates a map of NPC type keys to NPC types.
// It uses the GetCivilianInstance and GetMilitaryInstance functions from the npc_types package to get the NPC types.
func (c *CreationDataService) loadNpcTypeMap() map[string]t.NPCType {
	return map[string]t.NPCType{
		"Civilian": t.GetCivilianInstance().NPCType, // dereference the embedded field
		"Military": t.GetMilitaryInstance().NPCType,
	}
}

// buildNpcTypeNameMap creates a map of NPC type keys to NPC type names.
// It uses the subtype maps to get the NPC type names.
func (c *CreationDataService) buildNpcTypeNameMap() map[string][]string {
	return map[string][]string{
		"Civilian": slices.Collect(maps.Keys(c.civilianSubtypeMap)),
		"Military": slices.Collect(maps.Keys(c.militarySubtypeMap)),
	}
}

// mergeNpcSubtypeMaps merges the subtype maps into a single map.
// It returns the merged map.
func (c *CreationDataService) mergeNpcSubtypeMaps(subtypeMaps ...map[string]cp.NPCSubtype) map[string]cp.NPCSubtype {
	merged := make(map[string]cp.NPCSubtype)
	for _, subtypeMap := range subtypeMaps {
		for key, subtype := range subtypeMap {
			merged[key] = subtype
		}
	}
	return merged
}

// GetFactionData returns the faction data for the given key.
// It panics if the faction data cannot be found.
func (c *CreationDataService) GetFactionData(key string) cp.Faction {
	faction, ok := c.factionMap[key]
	if !ok {
		panic(fmt.Sprintf("faction not found: %s", key))
	}
	return faction
}

// GetTraitData returns the trait data for the given key.
// It panics if the trait data cannot be found.
func (c *CreationDataService) GetTraitData(key string) cp.Trait {
	trait, ok := c.traitMap[key]
	if !ok {
		panic(fmt.Sprintf("trait not found: %s", key))
	}
	return trait
}

// GetNameData returns the name data for the given key.
// It panics if the name data cannot be found.
func (c *CreationDataService) GetNameData(key string) cp.NameData {
	nd, ok := c.nameMap[key]
	if ok {
		return nd
	}
	panic(fmt.Sprintf("name not found: %s", key))
}

// GetSpeciesData returns the species data for the given key.
// It panics if the species data cannot be found.
func (c *CreationDataService) GetSpeciesData(key string) cp.Species {
	s, ok := c.speciesMap[key]
	if ok {
		return s
	}
	panic(fmt.Sprintf("species not found: %s", key))
}

// GetNpcTypeData returns the NPC type data for the given key.
// It panics if the NPC type data cannot be found.
func (c *CreationDataService) GetNpcTypeData(key string) t.NPCType {
	nt, ok := c.npcTypeMap[key]
	if ok {
		return nt
	}
	panic(fmt.Sprintf("npc type not found: %s", key))
}

// GetSpeciesName returns the name of the species for the given key.
// It panics if the species name cannot be found.
func (c *CreationDataService) GetNpcSubtypeData(key string) cp.NPCSubtype {
	ns, ok := c.npcSubtypeMap[key]
	if ok {
		return ns
	}
	panic(fmt.Sprintf("npc subtype not found: %s", key))
}

// GetSpeciesName returns the name of the species for the given key.
// It panics if the species name cannot be found.
func (c *CreationDataService) GetFactionMap() map[string]cp.Faction {
	return maps.Clone(c.factionMap)
}

// GetSpeciesMap returns the species map.
// It panics if the species map cannot be found.
func (c *CreationDataService) GetSpeciesMap() map[string]cp.Species {
	return maps.Clone(c.speciesMap)
}

// GetNpcSubtypeForTypeMap returns the NPC subtype for type map.
// It panics if the NPC subtype for type map cannot be found.
func (c *CreationDataService) GetNpcSubtypeForTypeMap() map[string][]string {
	return maps.Clone(c.npcSubtypeForTypeMap)
}

// GetTraitMap returns the trait map.
// It panics if the trait map cannot be found.
func (c *CreationDataService) GetTraitMap() map[string]cp.Trait {
	return maps.Clone(c.traitMap)
}

// GetNameMap returns the name map.
// It panics if the name map cannot be found.
func (c *CreationDataService) GetNpcTypeMap() map[string]t.NPCType {
	return maps.Clone(c.npcTypeMap)
}

// GetNpcSubtypeMap returns the NPC subtype map.
// It panics if the NPC subtype map cannot be found.
func (c *CreationDataService) GetSpeciesNameMap() map[string]string {
	return maps.Clone(c.speciesNameMap)
}
