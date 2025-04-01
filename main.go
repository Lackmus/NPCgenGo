package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo/controller"
	h "github.com/lackmus/npcgengo/helper"
	"github.com/lackmus/npcgengo/loader"
	"github.com/lackmus/npcgengo/service"
	"github.com/lackmus/npcgengo/view"
)

func main() {
	creationSupplier := service.NewNPCCreationSupplier(loader.NewJSONNPCConfigLoader("data/creation_data"))
	npcController := controller.NewNPCListController(
		loader.NewJSONNPCStorage("data/npc_database"),
		creationSupplier,
	)
	fyneView := view.NewFyneListView(npcController)

	// delete one NPC if there are any
	if len(npcController.GetAllNpcs()) > 0 {
		npcController.DeleteNPC(npcController.GetAllNpcs()[0].ID)
	}

	for range 5 {
		npc, err := service.CreateNPCWithOptions(h.Random, h.Random, creationSupplier)
		if err != nil {
			fmt.Println(err)
			return
		}
		npcController.AddNpc(npc)
	}

	fyneView.Run()
	tidyUp()
	// delete all NPCs
	npcController.DeleteAllNPC()

}

func tidyUp() {
	fmt.Println("Exited")
}

func printDirStructure(root string, indent string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		fmt.Println(indent + "|-- " + entry.Name())
		if entry.IsDir() {
			printDirStructure(filepath.Join(root, entry.Name()), indent+"    ")
		}
	}
}
