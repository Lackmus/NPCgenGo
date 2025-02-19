package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo/model"
)

// JsonDatabaseLoader speichert den Pfad zum JSON-Ordner
type JSONNPCStorage struct {
	Dir string
}

// NewJsonDatabaseLoader erstellt eine neue Instanz
func NewJSONNPCStorage(dir string) *JSONNPCStorage {
	return &JSONNPCStorage{Dir: dir}
}

// LoadNpc lädt einen NPC aus einer JSON-Datei
func (j *JSONNPCStorage) LoadNPC(id string) (model.NPC, error) {
	filename := filepath.Join(j.Dir, id+".json")
	file, err := os.Open(filename)
	if err != nil {
		return model.NPC{}, err
	}
	defer file.Close()

	var data model.NPC
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return model.NPC{}, err
	}
	return data, nil
}

// LoadAllNpc lädt alle NPCs aus dem Verzeichnis
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

		id := file.Name()[:len(file.Name())-5] // ".json" entfernen
		data, err := j.LoadNPC(id)
		if err != nil {
			fmt.Printf("Warnung: Konnte %s nicht laden: %v\n", file.Name(), err)
			continue
		}
		dataMap[id] = data
	}
	return dataMap, nil
}

// SaveNpc speichert einen NPC in eine JSON-Datei
func (j *JSONNPCStorage) SaveNPC(data model.NPC) error {
	if err := os.MkdirAll(j.Dir, os.ModePerm); err != nil {
		return err
	}

	filename := filepath.Join(j.Dir, data.ID+".json")

	// Datei öffnen mit O_WRONLY (schreiben), O_CREATE (erstellen, falls nicht existiert), O_TRUNC (alte Daten überschreiben)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// SaveAllNpc speichert alle NPCs aus einer Map
func (j *JSONNPCStorage) SaveAllNPC(dataMap map[string]model.NPC) error {
	for _, data := range dataMap {
		if err := j.SaveNPC(data); err != nil {
			return err
		}
	}
	return nil
}

// DeleteNpc löscht eine NPC-Datei
func (j *JSONNPCStorage) DeleteNPC(id string) error {
	filename := filepath.Join(j.Dir, id+".json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("NPC %s nicht gefunden", id)
	}
	return os.Remove(filename)
}

// DeleteAllNpc löscht das gesamte Verzeichnis (vorsichtig nutzen!)
func (j *JSONNPCStorage) DeleteAllNPC() error {
	return os.RemoveAll(j.Dir)
}
