package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/product/model"
	cp "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
	"github.com/lackmus/npcgengo/pkg/product/service"
)

func newWebServerForTests(t *testing.T) *Server {
	t.Helper()

	baseDir := t.TempDir()
	creationDir := filepath.Join(baseDir, "creation_data")
	dbDir := filepath.Join(baseDir, "npc_database")

	if err := loader.CreateSampleCreationData(creationDir); err != nil {
		t.Fatalf("failed to create sample creation data: %v", err)
	}

	creationSupplier, err := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader(creationDir))
	if err != nil {
		t.Fatalf("failed to create creation supplier: %v", err)
	}

	npcService, err := service.NewNPCService(context.Background(), loader.NewJSONNPCStorage(dbDir))
	if err != nil {
		t.Fatalf("failed to create NPC service: %v", err)
	}

	controller := controllers.NewNPCListController(creationSupplier, npcService)
	return NewServer(controller)
}

func TestServer_NpcsHandler_GetReturnsTypedDTOs(t *testing.T) {
	srv := newWebServerForTests(t)

	npc := *model.NewNPC()
	npc.ID = "n-1"
	npc.SetComponent(cp.CompName, "Alice Smith")
	npc.SetComponent(cp.CompType, "Civilian")
	npc.SetComponent(cp.CompSubtype, "someCivilianSubtypeID")
	npc.SetComponent(cp.CompSpecies, "someSpeciesID")
	npc.SetComponent(cp.CompFaction, "someFactionID")
	npc.SetComponent(cp.CompTrait, "someTraitID")
	npc.SetComponent(cp.CompStats, "STR: 2")
	npc.SetComponent(cp.CompItems, "Weapon: Fists")
	srv.npcController.AddNPC(npc)

	req := httptest.NewRequest(http.MethodGet, "/api/npcs", nil)
	res := httptest.NewRecorder()

	srv.npcsHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	var payload []map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(payload) != 1 {
		t.Fatalf("expected 1 NPC, got %d", len(payload))
	}

	item := payload[0]
	if item["id"] != "n-1" || item["name"] != "Alice Smith" || item["type"] != "Civilian" {
		t.Fatalf("unexpected DTO fields: %+v", item)
	}
	if _, hasComponents := item["Components"]; hasComponents {
		t.Fatalf("response leaked internal Components map: %+v", item)
	}
	if _, hasLowerComponents := item["components"]; hasLowerComponents {
		t.Fatalf("response leaked internal components map: %+v", item)
	}
}

func TestServer_GenerateHandler_ReturnsTypedDTO(t *testing.T) {
	srv := newWebServerForTests(t)

	req := httptest.NewRequest(http.MethodPost, "/api/generate", nil)
	res := httptest.NewRecorder()

	srv.generateHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	requiredFields := []string{"id", "name", "type", "subtype", "species", "faction", "trait", "stats", "items", "notes"}
	for _, field := range requiredFields {
		if _, ok := payload[field]; !ok {
			t.Fatalf("missing required DTO field %q in payload: %+v", field, payload)
		}
	}

	if _, hasComponents := payload["Components"]; hasComponents {
		t.Fatalf("response leaked internal Components map: %+v", payload)
	}
	if _, hasLowerComponents := payload["components"]; hasLowerComponents {
		t.Fatalf("response leaked internal components map: %+v", payload)
	}
}

func TestServer_NpcByIDHandler_GetReturnsTypedDTO(t *testing.T) {
	srv := newWebServerForTests(t)

	npc := *model.NewNPC()
	npc.ID = "n-2"
	npc.SetComponent(cp.CompName, "Bob Smith")
	npc.SetComponent(cp.CompType, "Civilian")
	npc.SetComponent(cp.CompSubtype, "someCivilianSubtypeID")
	npc.SetComponent(cp.CompSpecies, "someSpeciesID")
	npc.SetComponent(cp.CompFaction, "someFactionID")
	npc.SetComponent(cp.CompTrait, "someTraitID")
	npc.SetComponent(cp.CompStats, "STR: 4")
	npc.SetComponent(cp.CompItems, "Weapon: Fists")
	srv.npcController.AddNPC(npc)

	req := httptest.NewRequest(http.MethodGet, "/api/npcs/n-2", nil)
	res := httptest.NewRecorder()

	srv.npcByIDHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if payload["id"] != "n-2" || payload["name"] != "Bob Smith" || payload["type"] != "Civilian" {
		t.Fatalf("unexpected DTO fields: %+v", payload)
	}
	if _, hasComponents := payload["Components"]; hasComponents {
		t.Fatalf("response leaked internal Components map: %+v", payload)
	}
	if _, hasLowerComponents := payload["components"]; hasLowerComponents {
		t.Fatalf("response leaked internal components map: %+v", payload)
	}
}
