package loader

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

// JSONNPCGroupStorage is a JSON file-based storage for NPC groups
type JSONNPCGroupStorage struct {
	Dir string
}

// NewJSONNPCGroupStorage creates a new JSONNPCGroupStorage
func NewJSONNPCGroupStorage(dir string) shared.NPCGroupStorage {
	return &JSONNPCGroupStorage{Dir: dir}
}

// LoadNPCGroup loads an NPC group from a JSON file
func (j *JSONNPCGroupStorage) LoadNPCGroup(name string) (model.NPCGroup, error) {
	filename := filepath.Join(j.Dir, "groups", name+".json")
	file, err := os.Open(filename)
	group := *model.NewNPCGroup(name)
	group.Name = name
	if err != nil {
		return model.NPCGroup{}, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&group); err != nil {
		return model.NPCGroup{}, err
	}
	return group, nil
}

// LoadAllNPCGroups loads all NPC groups from JSON files in the directory
func (j *JSONNPCGroupStorage) LoadAllNPCGroups() (map[string]model.NPCGroup, error) {
	dataMap := make(map[string]model.NPCGroup)

	groupsDir := filepath.Join(j.Dir, "groups")
	if err := os.MkdirAll(groupsDir, os.ModePerm); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(groupsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		id := file.Name()[:len(file.Name())-5]
		data, err := j.LoadNPCGroup(id)
		if err != nil {
			log.Printf("Error loading NPC group %s: %v", id, err)
			continue
		}
		dataMap[id] = data
	}
	return dataMap, nil
}

// SaveNPCGroup saves an NPC group to a JSON file
func (j *JSONNPCGroupStorage) SaveNPCGroup(group model.NPCGroup) error {
	groupsDir := filepath.Join(j.Dir, "groups")
	if err := os.MkdirAll(groupsDir, os.ModePerm); err != nil {
		return err
	}

	filename := filepath.Join(groupsDir, group.Name+".json")

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(group)
}

// DeleteNPCGroup deletes an NPC group JSON file
func (j *JSONNPCGroupStorage) DeleteNPCGroup(id string) error {
	filename := filepath.Join(j.Dir, "groups", id+".json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("File %s does not exist: %v", filename, err)
		return nil
	}
	return os.Remove(filename)
}

// DeleteAllNPCGroups deletes all NPC group JSON files
func (j *JSONNPCGroupStorage) DeleteAllNPCGroups() error {
	groupsDir := filepath.Join(j.Dir, "groups")
	if err := os.MkdirAll(groupsDir, os.ModePerm); err != nil {
		return err
	}

	dir, err := os.ReadDir(groupsDir)
	if err != nil {
		log.Printf("Error reading directory %s: %v", groupsDir, err)
		return err
	}

	for _, file := range dir {
		if err := os.Remove(filepath.Join(groupsDir, file.Name())); err != nil {
			log.Printf("Error deleting file %s: %v", file.Name(), err)
		}
	}
	return nil
}
