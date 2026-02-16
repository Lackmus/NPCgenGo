package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components"
)

// Simple HTTP server that serves the web demo and exposes API endpoints
// Routes:
//  GET  /api/npcs         -> list all NPCs
//  GET  /api/npcs/:id     -> get NPC by id
//  POST /api/generate     -> create an NPC server-side and store it
//  DELETE /api/npcs/:id   -> delete NPC by id

type Server struct {
	npcController *NPCListController
}

func NewServer(nc *NPCListController) *Server {
	return &Server{npcController: nc}
}

func (s *Server) Routes() {
	http.Handle("/web_demo/", http.StripPrefix("/web_demo/", http.FileServer(http.Dir("web_demo"))))

	http.HandleFunc("/api/npcs", s.npcsHandler)
	http.HandleFunc("/api/npcs/", s.npcByIDHandler)
	http.HandleFunc("/api/generate", s.generateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *Server) npcsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodGet:
		npcs := s.npcController.GetAllNpcs()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(npcs)
	case http.MethodPost:
		m, err := parseNPCFromBody(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.npcController.AddNpc(m)
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) npcByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := strings.TrimPrefix(r.URL.Path, "/api/npcs/")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		npc, err := s.npcController.GetNpcByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(npc)
	case http.MethodDelete:
		s.npcController.DeleteNPC(id)
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPut:
		m, err := parseNPCFromBody(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.npcController.UpdateNpc(m)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) generateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	npc, err := s.npcController.CreateRandomNPC()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(npc)
}

func parseNPCFromBody(body io.ReadCloser) (model.NPC, error) {
	defer body.Close()
	var p struct {
		ID, Name, Type, Subtype, Species, Faction, Stats, Items, Description, LocationID string
		Traits                                                                           []string `json:"traits"`
	}
	if err := json.NewDecoder(body).Decode(&p); err != nil {
		return model.NPC{}, err
	}

	var m model.NPC
	m.ID = p.ID
	m.LocationID = p.LocationID
	m.Components = make(map[cp.CompEnum]string)
	if p.Name != "" {
		m.Components[cp.CompName] = p.Name
	}
	if p.Type != "" {
		m.Components[cp.CompType] = p.Type
	}
	if p.Subtype != "" {
		m.Components[cp.CompSubtype] = p.Subtype
	}
	if p.Species != "" {
		m.Components[cp.CompSpecies] = p.Species
	}
	if p.Faction != "" {
		m.Components[cp.CompFaction] = p.Faction
	}
	if len(p.Traits) > 0 {
		m.Components[cp.CompTrait] = strings.Join(p.Traits, ", ")
	}
	if p.Stats != "" {
		m.Components[cp.CompStats] = p.Stats
	}
	if p.Items != "" {
		m.Components[cp.CompItems] = p.Items
	}
	if p.Description != "" {
		m.Components[cp.CompDescription] = p.Description
	}
	return m, nil
}
