package main

import (
	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/loader"
	"github.com/lackmus/npcgengo/view"
)

func main() {

	npcController := controller.NewNPCListController(
		loader.NewJSONNPCStorage("data/npc_database"),
		loader.NewJSONNpcConfigLoader(),
		view.NewConsoleView(),
	)

	npcController.InitView()

	// new editview
	editView := view.NewConsoleEditView()
	editController := npcController.InitEditController(editView)

	for i := 0; i < 2; i++ {
		editController.CreateNPC("", "")
		npcController.AddNpc(editController.SaveNPC())
	}

}
