package main

import (
	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/loader"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/view"
)

func main() {

	creationSupplier := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader())
	npcController := controller.NewNPCListController(
		loader.NewJSONNPCStorage("data/npc_database"),
		creationSupplier,
		view.NewConsoleView(),
	)

	npcController.InitView()

	editView := view.NewConsoleEditView()
	editController := npcController.InitEditController(editView)

	for i := 0; i < 5; i++ {
		npc := service.CreateNPCWithOptions("", "", creationSupplier.RandomizerService)
		npcController.AddNpc(npc)
	}

	editController.LoadNPC(npcController.GetAllNpcs()[0])
	editController.RandomizeField("name")
	npcController.AddNpc(editController.SaveNPC())

	npcController.DeleteAllNPC()
}
