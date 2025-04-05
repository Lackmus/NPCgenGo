package main

import (
	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/loader"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/view"
)

// constants for the application
const (
	defaultLocationID = "default"
	npcDatabasePath   = "data/npc_database"
	creationDataPath  = "data/creation_data"
)

func main() {
	n := NewNPCGen()
	n.InitNPCListView(defaultLocationID)
}

type NPCGen struct {
	creationSupplier *service.NPCCreationSupplier
	npcService       *service.NPCService
}

// NewNPCGen initializes a new NPCGen instance with the provided parameters.
// It creates a new NPCCreationSupplier and NPCListController, and sets up the Fyne view.
func NewNPCGen() *NPCGen {
	// Create the supplier once and reuse it
	creationSupplier := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader(creationDataPath))
	npcService := service.NewNPCService(loader.NewJSONNPCStorage(npcDatabasePath))

	return &NPCGen{
		creationSupplier: creationSupplier, // This assumes supplier implements the interface
		npcService:       npcService,
	}
}

// New Controller and view instance
func (n *NPCGen) InitNPCListView(locationID string) {
	npcListController := controller.NewNPCListController(n.creationSupplier, n.npcService, locationID)
	npcListView := view.NewFyneListView(npcListController)
	npcListView.Run()
}

// GetCreationOptions returns the creation options from the supplier.
func (n *NPCGen) GetFactions() []string {
	return n.creationSupplier.CreationOptions.Factions
}
