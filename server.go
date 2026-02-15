package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components"
)

// Simple HTTP server that serves the web demo and exposes API endpoints
// Routes:
//  GET  /api/npcs         -> list all NPCs
//  GET  /api/npcs/:id     -> get NPC by id
//  POST /api/generate     -> create an NPC server-side and store it
//  DELETE /api/npcs/:id   -> delete NPC by id

func main() {
	n := NewNPCGen()

	// create controller and use it for handlers
	npcController := controller.NewNPCListController(n.creationSupplier, n.npcService, defaultLocationID)

	// Serve static files: web_demo and data
	http.Handle("/web_demo/", http.StripPrefix("/web_demo/", http.FileServer(http.Dir("web_demo"))))
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))

	// API: list and create (GET /api/npcs, POST /api/npcs)
	http.HandleFunc("/api/npcs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch r.Method {
		case http.MethodGet:
			npcs := npcController.GetAllNpcs()
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(npcs); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case http.MethodPost:
			var p struct {
				ID          string   `json:"id"`
				Name        string   `json:"name"`
				Type        string   `json:"type"`
				Subtype     string   `json:"subtype"`
				Species     string   `json:"species"`
				Faction     string   `json:"faction"`
				Traits      []string `json:"traits"`
				Stats       string   `json:"stats"`
				Items       string   `json:"items"`
				Description string   `json:"description"`
				LocationID  string   `json:"locationID"`
			}
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
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
			npcController.AddNpc(m)
			w.WriteHeader(http.StatusCreated)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	// API: get/delete/update by id
	http.HandleFunc("/api/npcs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := strings.TrimPrefix(r.URL.Path, "/api/npcs/")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			npc, err := npcController.GetNpcByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(npc)
			return
		case http.MethodDelete:
			npcController.DeleteNPC(id)
			w.WriteHeader(http.StatusNoContent)
			return
		case http.MethodPut:
			// decode simplified client payload and map into model.NPC
			var p struct {
				ID          string   `json:"id"`
				Name        string   `json:"name"`
				Type        string   `json:"type"`
				Subtype     string   `json:"subtype"`
				Species     string   `json:"species"`
				Faction     string   `json:"faction"`
				Traits      []string `json:"traits"`
				Stats       string   `json:"stats"`
				Items       string   `json:"items"`
				Description string   `json:"description"`
				LocationID  string   `json:"locationID"`
			}
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
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
			npcController.UpdateNpc(m)
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	// API: generate server-side
	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		npc, err := npcController.CreateRandomNPC()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(npc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Println("Starting server on :8080 — open http://localhost:8080/web_demo/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
