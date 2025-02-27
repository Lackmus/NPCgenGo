// Description: This file contains the JSONNPCConfigLoader struct and its methods.
package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	c "github.com/lackmus/npcgengo/model/npc_components"
	"github.com/lackmus/npcgengo/shared"
)

const (
	factionDir            = "factiondata"
	speciesDir            = "speciesdata"
	traitDir              = "traitdata"
	nameDir               = "namedata"
	npcCivilianSubtypeDir = "npctypedata/civilian"
	npcMilitarySubtypeDir = "npctypedata/military"
)

// JSONNPCConfigLoader loads NPC configuration data from JSON files
// in a directory.
type JSONNPCConfigLoader struct {
	dir string
}

// JSONNPCConfigLoader loads NPC configuration data from JSON files in a directory.
// It implements the shared.NPCConfigLoader interface.
func NewJSONNPCConfigLoader(dir string) shared.NPCConfigLoader {
	return &JSONNPCConfigLoader{
		dir: dir,
	}
}

// LoadFactionMap loads the faction data from the JSON files in the factiondata directory.
// It returns a map of faction IDs to faction data.
func (j *JSONNPCConfigLoader) LoadFactionMap() (map[string]c.Faction, error) {
	return loadJSONMap[c.Faction](filepath.Join(j.dir, factionDir))
}

// LoadSpeciesMap loads the species data from the JSON files in the speciesdata directory.
// It returns a map of species IDs to species data.
func (j *JSONNPCConfigLoader) LoadSpeciesMap() (map[string]c.Species, error) {
	return loadJSONMap[c.Species](filepath.Join(j.dir, speciesDir))
}

// LoadTraitMap loads the trait data from the JSON files in the traitdata directory.
// It returns a map of trait IDs to trait data.
func (j *JSONNPCConfigLoader) LoadTraitMap() (map[string]c.Trait, error) {
	return loadJSONMap[c.Trait](filepath.Join(j.dir, traitDir))
}

// LoadNameMap loads the name data from the JSON files in the namedata directory.
// It returns a map of name IDs to name data.
func (j *JSONNPCConfigLoader) LoadNameMap() (map[string]c.NameData, error) {
	return loadJSONMap[c.NameData](filepath.Join(j.dir, nameDir))
}

// LoadNpcCivilianSubtypeMap loads the civilian NPC subtype data from the JSON files in the npctypedata/civilian directory.
// It returns a map of civilian NPC subtype IDs to civilian NPC subtype data.
func (j *JSONNPCConfigLoader) LoadNpcCivilianSubtypeMap() (map[string]c.NPCSubtype, error) {
	return loadJSONMap[c.NPCSubtype](filepath.Join(j.dir, npcCivilianSubtypeDir))
}

// LoadNpcMilitarySubtypeMap loads the military NPC subtype data from the JSON files in the npctypedata/military directory.
// It returns a map of military NPC subtype IDs to military NPC subtype data.
func (j *JSONNPCConfigLoader) LoadNpcMilitarySubtypeMap() (map[string]c.NPCSubtype, error) {
	return loadJSONMap[c.NPCSubtype](filepath.Join(j.dir, npcMilitarySubtypeDir))
}

// loadJSONMap loads JSON files from a directory into a map.
// It returns a map of IDs to data and an error if one occurred.
func loadJSONMap[T shared.Nameable](dir string) (map[string]T, error) {
	dataMap := make(map[string]T)
	files, err := os.ReadDir(dir)
	if err != nil {
		return dataMap, fmt.Errorf("error reading directory %s: %w", dir, err)
	}
	var errs []error
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if !strings.EqualFold(ext, ".json") {
			continue
		}
		data, err := loadJSON[T](filepath.Join(dir, file.Name()))
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to load file %s: %w", file.Name(), err))
			continue
		}
		id := data.GetName()
		dataMap[id] = data
	}
	if len(errs) > 0 {
		return dataMap, errors.Join(errs...)
	}
	return dataMap, nil
}

// loadJSON loads a JSON file into a struct.
// It returns the struct and an error if one occurred.
func loadJSON[T any](filePath string) (T, error) {
	var result T
	file, err := os.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&result); err != nil {
		return result, fmt.Errorf("error decoding JSON from %s: %w", filePath, err)
	}
	if validatable, ok := any(result).(interface{ Validate() error }); ok {
		if err = validatable.Validate(); err != nil {
			return result, fmt.Errorf("validation failed for %s: %w", filePath, err)
		}
	}
	return result, nil
}
