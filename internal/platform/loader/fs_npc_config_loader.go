package loader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"

	c "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

type FSNPCConfigLoader struct {
	baseFS fs.FS
	base   string
}

//
func NewFSNPCConfigLoader(baseFS fs.FS, base string) shared.NPCConfigLoader {
	return &FSNPCConfigLoader{
		baseFS: baseFS,
		base:   strings.Trim(path.Clean(base), "/"),
	}
}

func (f *FSNPCConfigLoader) LoadFactionMap(ctx context.Context) (map[string]c.Faction, error) {
	return loadJSONMapFromFS[c.Faction](ctx, f.baseFS, path.Join(f.base, factionDir))
}

func (f *FSNPCConfigLoader) LoadSpeciesMap(ctx context.Context) (map[string]c.Species, error) {
	return loadJSONMapFromFS[c.Species](ctx, f.baseFS, path.Join(f.base, speciesDir))
}

func (f *FSNPCConfigLoader) LoadTraitMap(ctx context.Context) (map[string]c.Trait, error) {
	return loadJSONMapFromFS[c.Trait](ctx, f.baseFS, path.Join(f.base, traitDir))
}

func (f *FSNPCConfigLoader) LoadNameMap(ctx context.Context) (map[string]c.NameData, error) {
	return loadJSONMapFromFS[c.NameData](ctx, f.baseFS, path.Join(f.base, nameDir))
}

func (f *FSNPCConfigLoader) LoadNPCSubtypeMaps(ctx context.Context) (map[string]map[string]c.NPCSubtype, error) {
	dataMap := make(map[string]map[string]c.NPCSubtype)
	if err := ctx.Err(); err != nil {
		return dataMap, err
	}

	baseDir := path.Join(f.base, npcSubtypeDir)
	entries, err := fs.ReadDir(f.baseFS, baseDir)
	if err != nil {
		return dataMap, fmt.Errorf("error reading directory %s: %w", baseDir, err)
	}

	var errs []error
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return dataMap, err
		}
		if !entry.IsDir() {
			continue
		}

		subtypes, loadErr := loadJSONMapFromFS[c.NPCSubtype](ctx, f.baseFS, path.Join(baseDir, entry.Name()))
		if loadErr != nil {
			errs = append(errs, fmt.Errorf("failed to load subtype directory %s: %w", entry.Name(), loadErr))
			continue
		}

		typeName := resolveNPCTypeName(entry.Name(), subtypes)
		if _, ok := dataMap[typeName]; !ok {
			dataMap[typeName] = make(map[string]c.NPCSubtype)
		}

		for key, subtype := range subtypes {
			if strings.TrimSpace(subtype.NpcTypeName) == "" {
				subtype.NpcTypeName = typeName
			}
			dataMap[typeName][key] = subtype
		}
	}

	if len(errs) > 0 {
		return dataMap, errors.Join(errs...)
	}

	return dataMap, nil
}

func loadJSONMapFromFS[T shared.Nameable](ctx context.Context, baseFS fs.FS, dir string) (map[string]T, error) {
	dataMap := make(map[string]T)
	if err := ctx.Err(); err != nil {
		return dataMap, err
	}

	files, err := fs.ReadDir(baseFS, dir)
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
		if !strings.EqualFold(path.Ext(file.Name()), ".json") {
			continue
		}

		data, loadErr := loadJSONFromFS[T](ctx, baseFS, path.Join(dir, file.Name()))
		if loadErr != nil {
			errs = append(errs, fmt.Errorf("failed to load file %s: %w", file.Name(), loadErr))
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

func loadJSONFromFS[T any](ctx context.Context, baseFS fs.FS, filePath string) (T, error) {
	var result T
	if err := ctx.Err(); err != nil {
		return result, err
	}

	raw, err := fs.ReadFile(baseFS, filePath)
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
