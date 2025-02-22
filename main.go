package main

import (
	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/loader"
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

	for range 5 {
		npc := service.CreateNPCWithOptions("", "", creationSupplier.RandomizerService)
		npcController.AddNpc(npc)
	}

	for _, npc := range npcController.GetAllNpcs() {
		editController.LoadNPC(npc)
		editController.RandomizeField("name")
		npcController.UpdateNpc(editController.SaveNPC())
	}

	npcController.DeleteAllNPC()
}
