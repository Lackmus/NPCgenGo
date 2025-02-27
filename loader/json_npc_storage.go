// Description: JSON NPC storage loader implementation.
package loader

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

// JSONNPCStorage is a JSON file-based storage for NPCs.
// It implements the shared.NPCStorage interface.
type JSONNPCStorage struct {
	Dir string
}

// NewJSONNPCStorage creates a new JSONNPCStorage.
func NewJSONNPCStorage(dir string) shared.NPCStorage {
	return &JSONNPCStorage{Dir: dir}
}

// LoadNPC loads an NPC from a JSON file.
// It returns an error if the file cannot be opened or decoded.
func (j *JSONNPCStorage) LoadNPC(id string) (model.NPC, error) {
	filename := filepath.Join(j.Dir, id+".json")
	file, err := os.Open(filename)
	npc := *model.NewNPC(id)
	if err != nil {
		return model.NPC{}, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&npc); err != nil {
		return model.NPC{}, err
	}
	return npc, nil
}

// LoadAllNPC loads all NPCs from JSON files in the directory.
// It returns a map of NPC IDs to NPCs.
func (j *JSONNPCStorage) LoadAllNPC() (map[string]model.NPC, error) {
	dataMap := make(map[string]model.NPC)

	files, err := os.ReadDir(j.Dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		id := file.Name()[:len(file.Name())-5]
		data, err := j.LoadNPC(id)
		if err != nil {
			log.Printf("Error loading NPC %s: %v", id, err)
			continue
		}
		dataMap[id] = data
	}
	return dataMap, nil
}

// SaveNPC saves an NPC to a JSON file.
// It returns an error if the file cannot be created or encoded.
func (j *JSONNPCStorage) SaveNPC(npc model.NPC) error {
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

// SaveAllNPC saves all NPCs to JSON files.
// It returns an error if any NPC cannot be saved.
func (j *JSONNPCStorage) SaveAllNPC(dataMap map[string]model.NPC) error {
	for _, data := range dataMap {
		if err := j.SaveNPC(data); err != nil {
			return err
		}
	}
	return nil
}

// DeleteNPC deletes an NPC JSON file.
// It returns an error if the file cannot be deleted.
func (j *JSONNPCStorage) DeleteNPC(id string) error {
	filename := filepath.Join(j.Dir, id+".json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("File %s does not exist: %v", filename, err)
		return nil
	}
	return os.Remove(filename)
}

// DeleteAllNPC deletes all NPC JSON files.
// It returns an error if any file cannot be deleted.
func (j *JSONNPCStorage) DeleteAllNPC() error {
	dir, err := os.ReadDir(j.Dir)
	if err != nil {
		log.Printf("Error reading directory %s: %v", j.Dir, err)
	}
	for _, file := range dir {
		if err := os.Remove(filepath.Join(j.Dir, file.Name())); err != nil {
			log.Printf("Error deleting file %s: %v", file.Name(), err)
		}
	}
	return nil
}
