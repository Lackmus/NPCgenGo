package main

import (
	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/loader"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/view"
)

func main() {

	// service.creationsuppöier
	creationSupplier := service.NewNPCCreationSupplier(loader.NewJSONNpcConfigLoader())
	npcController := controller.NewNPCListController(
		loader.NewJSONNPCStorage("data/npc_database"),
		*creationSupplier,
		view.NewConsoleView(),
	)

	npcController.InitView()

	// new editview
	editView := view.NewConsoleEditView()
	editController := npcController.InitEditController(editView)

	for i := 0; i < 5; i++ {
		npc := service.CreateNPCWithOptions("", "", creationSupplier.RandomizerService)
		npcController.AddNpc(npc)
	}

	editController.LoadNPC(npcController.GetAllNpcs()[0])
	// change field with new value
	editController.RandomizeField("name")
	npcController.AddNpc(editController.SaveNPC())

	//delete all npcs
	npcController.DeleteAllNPC()
}
