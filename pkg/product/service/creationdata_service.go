// Description: This file contains the CreationDataService struct and its methods. The CreationDataService struct provides access to the creation data for NPCs. It is used to load the data from the creation data file into maps.
package service

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	cp "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
	t "github.com/lackmus/npcgengo/pkg/product/model/npc_components/types"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

// CreationDataService provides access to the creation data for NPCs.
// It is used to load the data from the creation data file into maps.
type CreationDataService struct {
	factionMap           map[string]cp.Faction
	speciesMap           map[string]cp.Species
	traitMap             map[string]cp.Trait
	nameMap              map[string]cp.NameData
	npcTypeMap           map[string]t.NPCType
	npcSubtypeMap        map[string]cp.NPCSubtype
	npcSubtypeByTypeMap  map[string]map[string]cp.NPCSubtype
	speciesNameMap       map[string]string
	npcSubtypeForTypeMap map[string][]string
}

func NewCreationDataService(ctx context.Context, npcConfigLoader shared.NPCConfigLoader) (*CreationDataService, error) {
	cds := &CreationDataService{}
	if err := cds.initConfigLoaderMaps(ctx, npcConfigLoader); err != nil {
		return nil, err
	}
	cds.npcSubtypeMap = cds.mergeNpcSubtypeMaps(cds.npcSubtypeByTypeMap)
	cds.npcSubtypeForTypeMap = cds.buildNpcSubtypeForTypeMap()
	cds.npcTypeMap = cds.buildNpcTypeMap()
	cds.speciesNameMap = cds.buildSpeciesNameMap()
	return cds, nil
}

func (c *CreationDataService) initConfigLoaderMaps(ctx context.Context, npcConfigLoader shared.NPCConfigLoader) error {
	factionMap, err := npcConfigLoader.LoadFactionMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to load faction map: %w", err)
	}
	speciesMap, err := npcConfigLoader.LoadSpeciesMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to load species map: %w", err)
	}
	traitMap, err := npcConfigLoader.LoadTraitMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to load trait map: %w", err)
	}
	nameMap, err := npcConfigLoader.LoadNameMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to load name map: %w", err)
	}
	npcSubtypeByTypeMap, err := npcConfigLoader.LoadNpcSubtypeMaps(ctx)
	if err != nil {
		return fmt.Errorf("failed to load npc subtype maps: %w", err)
	}

	c.factionMap = factionMap
	c.speciesMap = speciesMap
	c.traitMap = traitMap
	c.nameMap = nameMap
	c.npcSubtypeByTypeMap = npcSubtypeByTypeMap
	return nil
}

func (c *CreationDataService) buildSpeciesNameMap() map[string]string {
	snm := make(map[string]string)
	for key, species := range c.speciesMap {
		if nameData, ok := c.nameMap[species.NameSource]; ok {
			snm[key] = nameData.GetName()
		}
	}
	return snm
}

func (c *CreationDataService) buildNpcTypeMap() map[string]t.NPCType {
	typeMap := make(map[string]t.NPCType)
	for typeName := range c.npcSubtypeForTypeMap {
		typeMap[typeName] = t.NPCType{
			Name:        typeName,
			Description: fmt.Sprintf("A %s npc", strings.ToLower(typeName)),
			Stats:       []string{"health", "speed", "strength"},
		}
	}
	return typeMap
}

func (c *CreationDataService) buildNpcSubtypeForTypeMap() map[string][]string {
	result := make(map[string][]string)
	for typeName, subtypeMap := range c.npcSubtypeByTypeMap {
		result[typeName] = slices.Collect(maps.Keys(subtypeMap))
	}
	return result
}

func (c *CreationDataService) buildSpeciesForFactionMap() map[string][]string {
	result := make(map[string][]string)
	for factionName, faction := range c.factionMap {
		result[factionName] = faction.SpeciesList
	}
	return result
}

func (c *CreationDataService) mergeNpcSubtypeMaps(subtypeByTypeMap map[string]map[string]cp.NPCSubtype) map[string]cp.NPCSubtype {
	merged := make(map[string]cp.NPCSubtype)
	for _, subtypeMap := range subtypeByTypeMap {
		for key, subtype := range subtypeMap {
			merged[key] = subtype
		}
	}
	return merged
}

func (c *CreationDataService) GetFactionData(key string) (cp.Faction, error) {
	faction, ok := c.factionMap[key]
	if !ok {
		return cp.Faction{}, fmt.Errorf("faction not found: %s", key)
	}
	return faction, nil
}

func (c *CreationDataService) GetTraitData(key string) (cp.Trait, error) {
	trait, ok := c.traitMap[key]
	if !ok {
		return cp.Trait{}, fmt.Errorf("trait not found: %s", key)
	}
	return trait, nil
}

func (c *CreationDataService) GetNameData(key string) (cp.NameData, error) {
	nd, ok := c.nameMap[key]
	if ok {
		return nd, nil
	}
	return cp.NameData{}, fmt.Errorf("name not found: %s", key)
}

func (c *CreationDataService) GetSpeciesData(key string) (cp.Species, error) {
	s, ok := c.speciesMap[key]
	if ok {
		return s, nil
	}
	return cp.Species{}, fmt.Errorf("species not found: %s", key)
}

func (c *CreationDataService) GetNpcTypeData(key string) (t.NPCType, error) {
	nt, ok := c.npcTypeMap[key]
	if ok {
		return nt, nil
	}
	return t.NPCType{}, fmt.Errorf("npc type not found: %s", key)
}

func (c *CreationDataService) GetNpcSubtypeData(key string) (cp.NPCSubtype, error) {
	ns, ok := c.npcSubtypeMap[key]
	if ok {
		return ns, nil
	}
	return cp.NPCSubtype{}, fmt.Errorf("npc subtype not found: %s", key)
}

func (c *CreationDataService) GetFactionMap() map[string]cp.Faction {
	return maps.Clone(c.factionMap)
}

func (c *CreationDataService) GetSpeciesMap() map[string]cp.Species {
	return maps.Clone(c.speciesMap)
}

func (c *CreationDataService) GetNpcSubtypeForTypeMap() map[string][]string {
	return maps.Clone(c.npcSubtypeForTypeMap)
}

func (c *CreationDataService) GetTraitMap() map[string]cp.Trait {
	return maps.Clone(c.traitMap)
}

func (c *CreationDataService) GetNpcTypeMap() map[string]t.NPCType {
	return maps.Clone(c.npcTypeMap)
}

func (c *CreationDataService) GetSpeciesNameMap() map[string]string {
	return maps.Clone(c.speciesNameMap)
}

func (c *CreationDataService) GetSpeciesForFactionMap() map[string][]string {
	return c.buildSpeciesForFactionMap()
}
