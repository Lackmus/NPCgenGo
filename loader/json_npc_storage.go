package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo/model"
)

// JSONNPCStorage manages saving/loading NPCs as JSON.
type JSONNPCStorage struct {
	Dir string
}

// NewJSONNPCStorage creates a new instance.
func NewJSONNPCStorage(dir string) *JSONNPCStorage {
	return &JSONNPCStorage{Dir: dir}
}

// dtoNPC is a Data Transfer Object used only for JSON marshalling/unmarshalling.
type dtoNPC struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Faction     string            `json:"faction"`
	Species     string            `json:"species"`
	NpcType     string            `json:"npc_type"`
	NpcSubtype  string            `json:"npc_subtype"`
	Trait       string            `json:"trait"`
	Drive       string            `json:"drive"`
	Stats       map[string]int    `json:"stats"`
	Items       map[string]string `json:"items"`
	Abilities   map[string]string `json:"abilities"`
	Description string            `json:"description"`
}

// toDTO converts an immutable NPC into its DTO form.
func toDTO(n model.NPC) dtoNPC {
	return dtoNPC{
		ID:          n.ID(),
		Name:        n.Name(),
		Faction:     n.Faction(),
		Species:     n.Species(),
		NpcType:     n.NpcType(),
		NpcSubtype:  n.NpcSubtype(),
		Trait:       n.Trait(),
		Drive:       n.Drive(),
		Stats:       n.Stats(),
		Items:       n.Items(),
		Abilities:   n.Abilities(),
		Description: n.Description(),
	}
}

// fromDTO converts a dtoNPC to an immutable NPC.
func fromDTO(dto dtoNPC) model.NPC {
	return model.NewNPC(
		dto.ID,
		dto.Name,
		dto.Faction,
		dto.Species,
		dto.NpcType,
		dto.NpcSubtype,
		dto.Trait,
		dto.Drive,
		dto.Description,
		dto.Stats,
		dto.Items,
		dto.Abilities,
	)
}

// LoadNPC loads an NPC from a JSON file.
func (j *JSONNPCStorage) LoadNPC(id string) (model.NPC, error) {
	filename := filepath.Join(j.Dir, id+".json")
	file, err := os.Open(filename)
	if err != nil {
		return model.NPC{}, err
	}
	defer file.Close()

	var dto dtoNPC
	if err := json.NewDecoder(file).Decode(&dto); err != nil {
		return model.NPC{}, err
	}
	return fromDTO(dto), nil
}

// LoadAllNPC loads all NPCs from the directory.
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
		id := file.Name()[0 : len(file.Name())-len(".json")]
		npc, err := j.LoadNPC(id)
		if err != nil {
			fmt.Printf("Warning: Could not load %s: %v\n", file.Name(), err)
			continue
		}
		dataMap[id] = npc
	}
	return dataMap, nil
}

// SaveNPC saves an NPC to a JSON file.
func (j *JSONNPCStorage) SaveNPC(npc model.NPC) error {
	if err := os.MkdirAll(j.Dir, os.ModePerm); err != nil {
		return err
	}
	filename := filepath.Join(j.Dir, npc.ID()+".json")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(toDTO(npc))
}

// SaveAllNPC saves all NPCs from a map.
func (j *JSONNPCStorage) SaveAllNPC(dataMap map[string]model.NPC) error {
	for _, npc := range dataMap {
		if err := j.SaveNPC(npc); err != nil {
			return err
		}
	}
	return nil
}

// DeleteNPC deletes an NPC file.
func (j *JSONNPCStorage) DeleteNPC(id string) error {
	filename := filepath.Join(j.Dir, id+".json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("NPC %s not found", id)
	}
	return os.Remove(filename)
}

// DeleteAllNPC deletes the entire directory (use with caution!).
func (j *JSONNPCStorage) DeleteAllNPC() error {
	return os.RemoveAll(j.Dir)
}
