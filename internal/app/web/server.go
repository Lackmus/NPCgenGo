package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/internal/app/mapper"
	"github.com/lackmus/npcgengo/pkg/product/model"
)

// Simple HTTP server that serves the web UI and exposes API endpoints
// Routes:
//  GET  /api/npcs         -> list all NPCs
//  GET  /api/npcs/:id     -> get NPC by id
//  POST /api/generate     -> create an NPC server-side and store it
//  DELETE /api/npcs/:id   -> delete NPC by id

type Server struct {
	npcController *controllers.NPCListController
	httpServer    *http.Server
}

func NewServer(nc *controllers.NPCListController) *Server {
	return &Server{
		npcController: nc,
	}
}

// Start registers handlers and starts the HTTP server on the provided address.
// It returns any error to the caller instead of exiting the process, making
// it safe to use this package as a library.
func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui/web"))))

	mux.HandleFunc("/api/npcs", s.npcsHandler)
	mux.HandleFunc("/api/npcs/", s.npcByIDHandler)
	mux.HandleFunc("/api/species/", s.speciesNameRollHandler)
	mux.HandleFunc("/api/subtypes/", s.subtypeRollHandler)
	mux.HandleFunc("/api/generate", s.generateHandler)
	mux.HandleFunc("/api/options", s.optionsHandler)

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
		m, err := parseNPCFromBody(r.Body, s.npcController)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.npcController.ValidateNPC(m); err != nil {
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
		var original *model.NPC
		if npcData, err := s.npcController.GetNpcByID(id); err == nil {
			original = &npcData
		}
		m, err := parseNPCFromBodyWithOriginal(r.Body, s.npcController, original)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.npcController.ValidateNPC(m); err != nil {
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

func (s *Server) optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.npcController.GetCreationOptions())
}

type subtypeRollResponse struct {
	Stats string `json:"stats"`
	Items string `json:"items"`
}

type speciesNameRollResponse struct {
	Name string `json:"name"`
}

func (s *Server) subtypeRollHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/subtypes/")
	if !strings.HasSuffix(path, "/roll") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	rawSubtype := strings.TrimSuffix(path, "/roll")
	rawSubtype = strings.TrimSuffix(rawSubtype, "/")
	if rawSubtype == "" {
		http.Error(w, "missing subtype", http.StatusBadRequest)
		return
	}

	subtype, err := url.PathUnescape(rawSubtype)
	if err != nil {
		http.Error(w, "invalid subtype", http.StatusBadRequest)
		return
	}

	stats, items, err := s.npcController.GetSubtypeFields(subtype)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subtypeRollResponse{
		Stats: stats,
		Items: items,
	})
}

func (s *Server) speciesNameRollHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/species/")
	if !strings.HasSuffix(path, "/name") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	rawSpecies := strings.TrimSuffix(path, "/name")
	rawSpecies = strings.TrimSuffix(rawSpecies, "/")
	if rawSpecies == "" {
		http.Error(w, "missing species", http.StatusBadRequest)
		return
	}

	species, err := url.PathUnescape(rawSpecies)
	if err != nil {
		http.Error(w, "invalid species", http.StatusBadRequest)
		return
	}

	name, err := s.npcController.GetSpeciesName(species)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(speciesNameRollResponse{Name: name})
}

// Helper function to parse NPC data from request body and convert it to model.NPC struct.
// Expects JSON with fields like ID, Name, Type, Subtype, Species, Faction, Trait, Stats, and Items.
func parseNPCFromBody(body io.ReadCloser, controller *controllers.NPCListController) (model.NPC, error) {
	defer body.Close()
	var p mapper.NPCInput
	if err := json.NewDecoder(body).Decode(&p); err != nil {
		return model.NPC{}, err
	}

	return mapper.ToModelNPC(p, controller.GetNPCBuilder())
}

// parseNPCFromBodyWithOriginal is like parseNPCFromBody but passes the original NPC
// so unchanged values can be retained by the builder flow.
func parseNPCFromBodyWithOriginal(body io.ReadCloser, controller *controllers.NPCListController, original *model.NPC) (model.NPC, error) {
	defer body.Close()
	var p mapper.NPCInput
	if err := json.NewDecoder(body).Decode(&p); err != nil {
		return model.NPC{}, err
	}

	return mapper.ToModelNPCWithOriginal(p, controller.GetNPCBuilder(), original)
}
