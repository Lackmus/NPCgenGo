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
	FactionMap           map[string]model.Faction
	SpeciesMap           map[string]model.Species
	TraitMap             map[string]model.Trait
	NameMap              map[string]model.NameData
	NpcTypeMap           map[string]types.NPCType
	NpcSubtypeMap        map[string]types.NPCSubtype
	SpeciesNameMap       map[string]string   // register species name to species key
	NpcSubtypeForTypeMap map[string][]string // register npc subtype names to npc type key
	CivilianSubtypeMap   map[string]types.NPCSubtype
	MilitarySubtypeMap   map[string]types.NPCSubtype
	NpcConfigLoader      shared.NPCConfigLoader
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

	service := &CreationDataService{
		FactionMap:         factionMap,
		SpeciesMap:         speciesMap,
		TraitMap:           traitMap,
		NameMap:            nameMap,
		NpcTypeMap:         make(map[string]types.NPCType),
		CivilianSubtypeMap: civilianSubtypeMap,
		MilitarySubtypeMap: militarySubtypeMap,
		NpcConfigLoader:    npcConfigLoader,
	}

	service.NpcTypeMap = service.loadNpcTypeMap()
	service.SpeciesNameMap = service.buildSpeciesNameMap()
	service.NpcSubtypeForTypeMap = service.buildNpcTypeNameMap()
	service.NpcSubtypeMap = service.mergeNpcSubtypeMaps()

	return service, nil
}

// buildSpeciesNameMap builds a map of species keys to species names.
func (c *CreationDataService) buildSpeciesNameMap() map[string]string {
	speciesNameMap := make(map[string]string)
	for key, species := range c.SpeciesMap {
		if nameData, ok := c.NameMap[species.NameSource]; ok {
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
		"Civilian": slices.Collect(maps.Keys(c.CivilianSubtypeMap)),
		"Military": slices.Collect(maps.Keys(c.MilitarySubtypeMap)),
	}
}

// mergeNpcSubtypeMaps merges the civilian and military subtype maps into a single map.
func (c *CreationDataService) mergeNpcSubtypeMaps() map[string]types.NPCSubtype {
	merged := maps.Clone(c.CivilianSubtypeMap)
	for key, subtype := range c.MilitarySubtypeMap {
		merged[key] = subtype
	}
	return merged
}

// GetFactionData returns the faction data for the given key.
func (c *CreationDataService) GetFactionData(key string) model.Faction {
	faction, ok := c.FactionMap[key]
	if !ok {
		panic(fmt.Sprintf("faction not found: %s", key))
	}
	return faction
}

// GetTraitData returns the trait data for the given key.
func (c *CreationDataService) GetTraitData(key string) model.Trait {
	trait, ok := c.TraitMap[key]
	if !ok {
		panic(fmt.Sprintf("trait not found: %s", key))
	}
	return trait
}

// GetNameData returns the name data for the given key.
func (c *CreationDataService) GetNameData(key string) model.NameData {
	nameData, ok := c.NameMap[key]
	if !ok {
		panic(fmt.Sprintf("name not found: %s", key))
	}
	return nameData
}

// GetSpeciesData returns the species data for the given key.
func (c *CreationDataService) GetSpeciesData(key string) model.Species {
	species, ok := c.SpeciesMap[key]
	if !ok {
		panic(fmt.Sprintf("species not found: %s", key))
	}
	return species
}

// GetNpcTypeData returns the npc type data for the given key.
func (c *CreationDataService) GetNpcTypeData(key string) types.NPCType {
	npcType, ok := c.NpcTypeMap[key]
	if !ok {
		panic(fmt.Sprintf("npc type not found: %s", key))
	}
	return npcType
}

// GetNpcSubtypeData returns the npc subtype data for the given key.
func (c *CreationDataService) GetNpcSubtypeData(key string) types.NPCSubtype {
	npcSubtype, ok := c.NpcSubtypeMap[key]
	if !ok {
		fmt.Printf("npc subtype not found: %s", key)
		//panic(fmt.Sprintf("npc subtype not found: %s", key))
	}
	return npcSubtype
}
