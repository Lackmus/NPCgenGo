package main

import (
	"fmt"

	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/loader"
	cp "github.com/lackmus/npcgengo/model/npc_components"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/view"
)

func main() {

	creationSupplier := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader("data/creation_data"))
	npcController := controller.NewNPCListController(
		loader.NewJSONNPCStorage("data/npc_database"),
		creationSupplier,
		view.NewConsoleView(),
	)

	npcController.InitView()

	editView := view.NewConsoleEditView()
	editController := npcController.InitEditController(editView)

	// delete one NPC if there are any
	if len(npcController.GetAllNpcs()) > 0 {
		npcController.DeleteNPC(npcController.GetAllNpcs()[0].ID)
	}

	for range 5 {
		npc := service.CreateNPCWithOptions(creationSupplier)
		npcController.AddNpc(npc)
	}

	// edit a random NPC name
	npc := npcController.GetAllNpcs()[0]
	editController.LoadNPC(npc)
	builder := service.NewNPCBuilderFromNPC(creationSupplier, npc)
	fmt.Println("\nUpdate NPC with name: " + npc.GetComponent(cp.CompName))
	builder.WithName().Build()
	fmt.Println("for Name: " + npc.GetComponent(cp.CompName))
	// get id
	npcController.UpdateNpc(editController.SaveNPC())
	npcController.DeleteAllNPC()

	// add a new NPC
	npc = service.CreateNPCWithOptions(creationSupplier)
	npcController.AddNpc(npc)
}
