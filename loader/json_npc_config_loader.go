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
)

// JsonNpcConfigLoader loads NPC generation data from JSON files
// =============================================================================
// JsonNpcConfigLoader
// =============================================================================

// JsonNpcConfigLoader loads NPC configuration data from JSON files stored in directories.
type JSONNPCConfigLoader struct {
	Dir                   string
	factionDir            string
	speciesDir            string
	traitDir              string
	nameDir               string
	npcCivilianSubtypeDir string
	npcMilitarySubtypeDir string
}

// NewJSONNpcConfigLoader creates a new instance of JSONNpcConfigLoader with default directory paths.
func NewJSONNpcConfigLoader() *JSONNPCConfigLoader {
	return &JSONNPCConfigLoader{
		Dir:                   "data/creation_data",
		factionDir:            "factiondata",
		speciesDir:            "speciesdata",
		traitDir:              "traitdata",
		nameDir:               "namedata",
		npcCivilianSubtypeDir: "npctypedata/civilian",
		npcMilitarySubtypeDir: "npctypedata/military",
	}
}

// LoadFactionMap loads faction configuration data.
func (j *JSONNPCConfigLoader) LoadFactionMap() (map[string]model.Faction, error) {
	return LoadJSONMap[model.Faction](filepath.Join(j.Dir, j.factionDir))
}

// LoadSpeciesMap loads species configuration data.
func (j *JSONNPCConfigLoader) LoadSpeciesMap() (map[string]model.Species, error) {
	return LoadJSONMap[model.Species](filepath.Join(j.Dir, j.speciesDir))
}

// LoadTraitMap loads trait configuration data.
func (j *JSONNPCConfigLoader) LoadTraitMap() (map[string]model.Trait, error) {
	return LoadJSONMap[model.Trait](filepath.Join(j.Dir, j.traitDir))
}

// LoadNameMap loads name data configuration.
func (j *JSONNPCConfigLoader) LoadNameMap() (map[string]model.NameData, error) {
	return LoadJSONMap[model.NameData](filepath.Join(j.Dir, j.nameDir))
}

// LoadNpcCivilianSubtypeMap loads civilian NPC subtype configuration data.
func (j *JSONNPCConfigLoader) LoadNpcCivilianSubtypeMap() (map[string]types.NPCSubtype, error) {
	return LoadJSONMap[types.NPCSubtype](filepath.Join(j.Dir, j.npcCivilianSubtypeDir))
}

// LoadNpcMilitarySubtypeMap loads military NPC subtype configuration data.
func (j *JSONNPCConfigLoader) LoadNpcMilitarySubtypeMap() (map[string]types.NPCSubtype, error) {
	return LoadJSONMap[types.NPCSubtype](filepath.Join(j.Dir, j.npcMilitarySubtypeDir))
}

// LoadJSONMap loads all JSON files from the given directory into a map keyed by filename (without extension).
// It returns an aggregated error if one or more files fail to load.
func LoadJSONMap[T any](dir string) (map[string]T, error) {
	dataMap := make(map[string]T)
	files, err := os.ReadDir(dir)
	if err != nil {
		return dataMap, fmt.Errorf("error reading directory %s: %w", dir, err)
	}

	var errs []error
	for _, file := range files {
		// Skip directories.
		if file.IsDir() {
			continue
		}

		// Check file extension (case-insensitive).
		ext := filepath.Ext(file.Name())
		if !strings.EqualFold(ext, ".json") {
			continue
		}

		// Use the filename (without extension) as the key.
		id := file.Name()[:len(file.Name())-len(ext)]
		data, err := LoadJSON[T](filepath.Join(dir, file.Name()))
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

// LoadJSON reads and decodes a JSON file into the provided generic type T.
// If T implements a Validate() error method, it is called to perform schema validation.
func LoadJSON[T any](filePath string) (T, error) {
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

	// If the decoded result implements a Validate() method, invoke it.
	if validatable, ok := any(result).(interface{ Validate() error }); ok {
		if err = validatable.Validate(); err != nil {
			return result, fmt.Errorf("validation failed for %s: %w", filePath, err)
		}
	}

	return result, nil
}
