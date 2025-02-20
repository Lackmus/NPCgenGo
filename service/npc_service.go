package service

import (
	"fmt"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

// ==============================
// NPC SERVICE
// ==============================

// NpcService verwaltet NPCs und lädt sie aus einer Datenbank
type NPCService struct {
	loader shared.NPCStorage
	npcs   map[string]model.NPC // Statt Slice jetzt ein Map für schnelleren Zugriff
}

// NewNpcService erstellt einen neuen NpcService und lädt vorhandene NPCs
func NewNPCService(loader shared.NPCStorage) *NPCService {
	npcMap, err := loader.LoadAllNPC()
	if err != nil {
		fmt.Println("Fehler beim Laden der NPCs:", err)
		npcMap = make(map[string]model.NPC) // Falls Fehler, leere Map erstellen
	}

	return &NPCService{
		loader: loader,
		npcs:   npcMap,
	}
}

// AddNpc fügt einen neuen NPC hinzu und speichert ihn direkt
func (s *NPCService) AddNpc(npc model.NPC) {
	s.npcs[npc.ID()] = npc
	s.loader.SaveNPC(npc)
}

// GetAllNPCs gibt alle NPCs als Slice zurück
func (s *NPCService) GetAllNpcs() []model.NPC {
	npcList := make([]model.NPC, 0, len(s.npcs))
	for _, npc := range s.npcs {
		npcList = append(npcList, npc)
	}
	return npcList
}

// GetNpcByID gibt den NPC mit der angegebenen ID zurück
func (s *NPCService) GetNpcByID(id string) (model.NPC, error) {
	npc, found := s.npcs[id]
	if !found {
		return model.NPC{}, fmt.Errorf("NPC mit ID %s nicht gefunden", id)
	}
	return npc, nil
}

// UpdateNpc aktualisiert einen bestehenden NPC und speichert ihn
func (s *NPCService) UpdateNPC(updatedNpc model.NPC) error {
	id := updatedNpc.ID()
	if _, found := s.npcs[id]; !found {
		return fmt.Errorf("NPC mit ID %s nicht gefunden", id)
	}
	s.npcs[id] = updatedNpc
	return s.loader.SaveNPC(updatedNpc)
}

// DeleteNpc entfernt einen NPC aus der Map und aus der JSON-Datei
func (s *NPCService) DeleteNPC(id string) error {
	if _, found := s.npcs[id]; !found {
		return fmt.Errorf("NPC mit ID %s nicht gefunden", id)
	}
	delete(s.npcs, id)
	return s.loader.DeleteNPC(id)
}

// DeleteAllNPCs löscht alle NPCs
func (s *NPCService) DeleteAllNPC() {
	s.loader.DeleteAllNPC()
	s.npcs = make(map[string]model.NPC)
}

// CountNPCs gibt die Anzahl der NPCs zurück
func (s *NPCService) CountNPC() int {
	return len(s.npcs)
}

// PrintNPCs gibt alle NPCs aus
func (s *NPCService) PrintAllNPC() {
	for _, npc := range s.npcs {
		fmt.Println(npc)
	}
}
