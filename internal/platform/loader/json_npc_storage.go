// Description: JSON NPC storage loader implementation.
package loader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo/pkg/model"
	"github.com/lackmus/npcgengo/pkg/shared"
)

type JSONNPCStorage struct {
	Dir string
}

func NewJSONNPCStorage(dir string) shared.NPCStorage {
	return &JSONNPCStorage{Dir: dir}
}

func (j *JSONNPCStorage) LoadNPC(ctx context.Context, id string) (model.NPC, error) {
	filename := filepath.Join(j.Dir, id+".json")
	var npc model.NPC = *model.NewNPC()
	npc.ID = id
	if err := ctx.Err(); err != nil {
		return model.NPC{}, err
	}
	raw, err := os.ReadFile(filename)
	if err != nil {
		return model.NPC{}, err
	}
	if err := json.Unmarshal(raw, &npc); err != nil {
		return model.NPC{}, err
	}
	return npc, nil
}

func (j *JSONNPCStorage) LoadAllNPC(ctx context.Context) (map[string]model.NPC, error) {
	dataMap := make(map[string]model.NPC)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	files, err := os.ReadDir(j.Dir)
	if err != nil {
		return nil, err
	}

	var errs []error
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return dataMap, err
		}
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		id := file.Name()[:len(file.Name())-5]
		data, err := j.LoadNPC(ctx, id)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", id, err))
			continue
		}
		dataMap[id] = data
	}
	if len(errs) > 0 {
		return dataMap, errors.Join(errs...)
	}
	return dataMap, nil
}

func (j *JSONNPCStorage) SaveNPC(ctx context.Context, npc model.NPC) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.MkdirAll(j.Dir, os.ModePerm); err != nil {
		return err
	}

	filename := filepath.Join(j.Dir, npc.ID+".json")

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(npc)
}

func (j *JSONNPCStorage) SaveAllNPC(ctx context.Context, dataMap map[string]model.NPC) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, data := range dataMap {
		if err := j.SaveNPC(ctx, data); err != nil {
			return err
		}
	}
	return nil
}

func (j *JSONNPCStorage) DeleteNPC(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	filename := filepath.Join(j.Dir, id+".json")
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			// Nothing to delete
			return nil
		}
		return err
	}
	return os.Remove(filename)
}

func (j *JSONNPCStorage) DeleteAllNPC(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	files, err := os.ReadDir(j.Dir)
	if err != nil {
		return err
	}
	var errs []error
	for _, file := range files {
		if err := os.Remove(filepath.Join(j.Dir, file.Name())); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

