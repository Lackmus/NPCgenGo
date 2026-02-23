// Description: This file contains the JSONNPCConfigLoader struct and its methods.
package loader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	c "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

const (
	factionDir            = "factiondata"
	speciesDir            = "speciesdata"
	traitDir              = "traitdata"
	nameDir               = "namedata"
	npcCivilianSubtypeDir = "npctypedata/civilian"
	npcMilitarySubtypeDir = "npctypedata/military"
)

type JSONNPCConfigLoader struct {
	dir string
}

func NewJSONNPCConfigLoader(dir string) shared.NPCConfigLoader {
	return &JSONNPCConfigLoader{
		dir: dir,
	}
}
func (j *JSONNPCConfigLoader) LoadFactionMap(ctx context.Context) (map[string]c.Faction, error) {
	return loadJSONMap[c.Faction](ctx, filepath.Join(j.dir, factionDir))
}

func (j *JSONNPCConfigLoader) LoadSpeciesMap(ctx context.Context) (map[string]c.Species, error) {
	return loadJSONMap[c.Species](ctx, filepath.Join(j.dir, speciesDir))
}

func (j *JSONNPCConfigLoader) LoadTraitMap(ctx context.Context) (map[string]c.Trait, error) {
	return loadJSONMap[c.Trait](ctx, filepath.Join(j.dir, traitDir))
}

func (j *JSONNPCConfigLoader) LoadNameMap(ctx context.Context) (map[string]c.NameData, error) {
	return loadJSONMap[c.NameData](ctx, filepath.Join(j.dir, nameDir))
}

func (j *JSONNPCConfigLoader) LoadNpcCivilianSubtypeMap(ctx context.Context) (map[string]c.NPCSubtype, error) {
	return loadJSONMap[c.NPCSubtype](ctx, filepath.Join(j.dir, npcCivilianSubtypeDir))
}

func (j *JSONNPCConfigLoader) LoadNpcMilitarySubtypeMap(ctx context.Context) (map[string]c.NPCSubtype, error) {
	return loadJSONMap[c.NPCSubtype](ctx, filepath.Join(j.dir, npcMilitarySubtypeDir))
}

func loadJSONMap[T shared.Nameable](ctx context.Context, dir string) (map[string]T, error) {
	dataMap := make(map[string]T)
	if err := ctx.Err(); err != nil {
		return dataMap, err
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return dataMap, fmt.Errorf("error reading directory %s: %w", dir, err)
	}
	var errs []error
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return dataMap, err
		}
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if !strings.EqualFold(ext, ".json") {
			continue
		}
		data, err := loadJSON[T](ctx, filepath.Join(dir, file.Name()))
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

// loadJSON reads a JSON file and unmarshals it into the provided type.
// loadJSON reads a JSON file and unmarshals it into the provided type.
func loadJSON[T any](ctx context.Context, filePath string) (T, error) {
	var result T
	if err := ctx.Err(); err != nil {
		return result, err
	}
	// read full file into memory (config files are small)
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return result, fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return result, fmt.Errorf("error decoding JSON from %s: %w", filePath, err)
	}
	if validatable, ok := any(result).(interface{ Validate() error }); ok {
		if err = validatable.Validate(); err != nil {
			return result, fmt.Errorf("validation failed for %s: %w", filePath, err)
		}
	}
	return result, nil
}

