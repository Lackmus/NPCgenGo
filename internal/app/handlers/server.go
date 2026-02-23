package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lackmus/npcgengo/pkg/product/model"
	cp "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
)

// Simple HTTP server that serves the web demo and exposes API endpoints
// Routes:
//  GET  /api/npcs         -> list all NPCs
//  GET  /api/npcs/:id     -> get NPC by id
//  POST /api/generate     -> create an NPC server-side and store it
//  DELETE /api/npcs/:id   -> delete NPC by id

type Server struct {
	npcController *NPCListController
	httpServer    *http.Server
}

func NewServer(nc *NPCListController) *Server {
	return &Server{npcController: nc}
}

// Routes registers HTTP handlers for the server. It is called by main() to set up the server before starting it.
// This method is deprecated in favor of Start(), which also registers handlers but returns any error instead of exiting the process.
func (s *Server) Routes() {
	// Deprecated compatibility method: keep old behavior (process-exiting)
	// Register handlers and start listening on :8080; any error will be fatal.
	http.Handle("/web_demo/", http.StripPrefix("/web_demo/", http.FileServer(http.Dir("web_demo"))))

	http.HandleFunc("/api/npcs", s.npcsHandler)
	http.HandleFunc("/api/npcs/", s.npcByIDHandler)
	http.HandleFunc("/api/generate", s.generateHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// Start registers handlers and starts the HTTP server on the provided address.
// It returns any error to the caller instead of exiting the process, making
// it safe to use this package as a library.
func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/web_demo/", http.StripPrefix("/web_demo/", http.FileServer(http.Dir("web_demo"))))

	mux.HandleFunc("/api/npcs", s.npcsHandler)
	mux.HandleFunc("/api/npcs/", s.npcByIDHandler)
	mux.HandleFunc("/api/generate", s.generateHandler)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server using the provided context.
// If Start has not been called or the server is not running, Shutdown is a no-op.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

// Handler for GET/POST /api/npcs - list all NPCs or create a new NPC.
// For GET, returns a JSON array of all NPCs. For POST, expects a JSON body with NPC data, creates it, and returns the created NPC as JSON.
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

// Handler for GET/DELETE/PUT /api/npcs/:id - get, delete, or update an NPC by ID.
// For GET, returns the NPC as JSON. For DELETE, removes it from the service. For PUT, updates it with the provided data.
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

// Handler for POST /api/generate - creates a new random NPC server-side and returns it as JSON.
// The NPC is also stored in the service, so it will be included in subsequent GET /api/npcs responses.
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

// Helper function to parse NPC data from request body and convert it to model.NPC struct.
// Expects JSON with fields like ID, Name, Type, Subtype, Species, Faction, Stats, Items, Description, LocationID, and Traits (array of strings).
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
