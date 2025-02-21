package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/model/types"
	"github.com/lackmus/npcgengo/shared"
)

const (
	dir                   = "data/creation_data"
	factionDir            = "factiondata"
	speciesDir            = "speciesdata"
	traitDir              = "traitdata"
	nameDir               = "namedata"
	npcCivilianSubtypeDir = "npctypedata/civilian"
	npcMilitarySubtypeDir = "npctypedata/military"
)

// JSONNPCConfigLoader loads the NPC configuration data from JSON files.
type JSONNPCConfigLoader struct{}

func NewJSONNPCConfigLoader() shared.NPCConfigLoader {
	return &JSONNPCConfigLoader{}
}

func (j *JSONNPCConfigLoader) LoadFactionMap() (map[string]model.Faction, error) {
	return loadJSONMap[model.Faction](filepath.Join(dir, factionDir))
}

func (j *JSONNPCConfigLoader) LoadSpeciesMap() (map[string]model.Species, error) {
	return loadJSONMap[model.Species](filepath.Join(dir, speciesDir))
}

func (j *JSONNPCConfigLoader) LoadTraitMap() (map[string]model.Trait, error) {
	return loadJSONMap[model.Trait](filepath.Join(dir, traitDir))
}

func (j *JSONNPCConfigLoader) LoadNameMap() (map[string]model.NameData, error) {
	return loadJSONMap[model.NameData](filepath.Join(dir, nameDir))
}

func (j *JSONNPCConfigLoader) LoadNpcCivilianSubtypeMap() (map[string]types.NPCSubtype, error) {
	return loadJSONMap[types.NPCSubtype](filepath.Join(dir, npcCivilianSubtypeDir))
}

func (j *JSONNPCConfigLoader) LoadNpcMilitarySubtypeMap() (map[string]types.NPCSubtype, error) {
	return loadJSONMap[types.NPCSubtype](filepath.Join(dir, npcMilitarySubtypeDir))
}

func loadJSONMap[T any](dir string) (map[string]T, error) {
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
		id := file.Name()[:len(file.Name())-len(ext)]
		data, err := loadJSON[T](filepath.Join(dir, file.Name()))
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to load file %s: %w", file.Name(), err))
			continue
		}
		dataMap[id] = data
	}
	if len(errs) > 0 {
		return dataMap, errors.Join(errs...)
	}
	return dataMap, nil
}

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
