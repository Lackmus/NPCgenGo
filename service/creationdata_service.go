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
	factionMap           map[string]model.Faction
	speciesMap           map[string]model.Species
	traitMap             map[string]model.Trait
	nameMap              map[string]model.NameData
	npcTypeMap           map[string]types.NPCType
	npcSubtypeMap        map[string]types.NPCSubtype
	speciesNameMap       map[string]string   // register species name to species key
	npcSubtypeForTypeMap map[string][]string // register npc subtype names to npc type key
	civilianSubtypeMap   map[string]types.NPCSubtype
	militarySubtypeMap   map[string]types.NPCSubtype
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

// buildSpeciesNameMap builds a map of species keys to species names.
func (c *CreationDataService) buildSpeciesNameMap() map[string]string {
	speciesNameMap := make(map[string]string)
	for key, species := range c.speciesMap {
		if nameData, ok := c.nameMap[species.NameSource]; ok {
			speciesNameMap[key] = nameData.Name
		}
	}
	return speciesNameMap
}

// loadNpcTypeMap initializes the NPC type map.
func (c *CreationDataService) loadNpcTypeMap() map[string]types.NPCType {
	return map[string]types.NPCType{
		"Civilian": types.GetCivilianInstance().NPCType, // dereference the embedded field
		"Military": types.GetMilitaryInstance().NPCType,
	}
}

// buildNpcTypeNameMap builds a map of npc type keys to npc type names.
func (c *CreationDataService) buildNpcTypeNameMap() map[string][]string {
	return map[string][]string{
		"Civilian": slices.Collect(maps.Keys(c.civilianSubtypeMap)),
		"Military": slices.Collect(maps.Keys(c.militarySubtypeMap)),
	}
}

// mergeNpcSubtypeMaps merges the civilian and military subtype maps into a single map.
func (c *CreationDataService) mergeNpcSubtypeMaps() map[string]types.NPCSubtype {
	merged := maps.Clone(c.civilianSubtypeMap)
	for key, subtype := range c.militarySubtypeMap {
		merged[key] = subtype
	}
	return merged
}

// GetFactionData returns the faction data for the given key.
func (c *CreationDataService) GetFactionData(key string) model.Faction {
	faction, ok := c.factionMap[key]
	if !ok {
		panic(fmt.Sprintf("faction not found: %s", key))
	}
	return faction
}

// GetTraitData returns the trait data for the given key.
func (c *CreationDataService) GetTraitData(key string) model.Trait {
	trait, ok := c.traitMap[key]
	if !ok {
		panic(fmt.Sprintf("trait not found: %s", key))
	}
	return trait
}

// GetNameData returns the name data for the given key.
func (c *CreationDataService) GetNameData(key string) model.NameData {
	nameData, ok := c.nameMap[key]
	if !ok {
		panic(fmt.Sprintf("name not found: %s", key))
	}
	return nameData
}

// GetSpeciesData returns the species data for the given key.
func (c *CreationDataService) GetSpeciesData(key string) model.Species {
	species, ok := c.speciesMap[key]
	if !ok {
		panic(fmt.Sprintf("species not found: %s", key))
	}
	return species
}

// GetNpcTypeData returns the npc type data for the given key.
func (c *CreationDataService) GetNpcTypeData(key string) types.NPCType {
	npcType, ok := c.npcTypeMap[key]
	if !ok {
		panic(fmt.Sprintf("npc type not found: %s", key))
	}
	return npcType
}

// GetNpcSubtypeData returns the npc subtype data for the given key.
func (c *CreationDataService) GetNpcSubtypeData(key string) types.NPCSubtype {
	npcSubtype, ok := c.npcSubtypeMap[key]
	if !ok {
		fmt.Printf("npc subtype not found: %s", key)
		//panic(fmt.Sprintf("npc subtype not found: %s", key))
	}
	return npcSubtype
}

func (c *CreationDataService) GetFactionMap() map[string]model.Faction {
	return maps.Clone(c.factionMap) // Returning a shallow copy of the map to avoid mutation
}

func (c *CreationDataService) GetSpeciesMap() map[string]model.Species {
	return maps.Clone(c.speciesMap) // Returning a shallow copy of the map to avoid mutation
}

func (c *CreationDataService) GetNpcSubtypeForTypeMap() map[string][]string {
	return maps.Clone(c.npcSubtypeForTypeMap) // Returning a shallow copy of the map to avoid mutation
}

// GetTraitMap returns a copy of the trait map.
func (c *CreationDataService) GetTraitMap() map[string]model.Trait {
	return maps.Clone(c.traitMap) // Returning a shallow copy of the map to avoid mutation
}

// GetNpcTypeMap returns a copy of the npc type map.
func (c *CreationDataService) GetNpcTypeMap() map[string]types.NPCType {
	return maps.Clone(c.npcTypeMap) // Returning a shallow copy of the map to avoid mutation
}

// GetSpeciesNameMap returns a copy of the species name map.
func (c *CreationDataService) GetSpeciesNameMap() map[string]string {
	return maps.Clone(c.speciesNameMap) // Returning a shallow copy of the map to avoid mutation
}
