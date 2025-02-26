package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo/controller"
	h "github.com/lackmus/npcgengo/helper"
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
		npc, err := service.CreateNPCWithOptions(h.Random, h.Random, creationSupplier)
		if err != nil {
			fmt.Println(err)
			return
		}
		npcController.AddNpc(npc)
	}

	// edit a random NPC name
	npc := npcController.GetAllNpcs()[0]
	editController.LoadNPC(npc)
	builder := service.NewNPCBuilder(creationSupplier).WithNPC(npc)
	fmt.Println("\nUpdate NPC with name: " + npc.GetComponent(cp.CompName))
	builder.WithRandomName().Build()
	fmt.Println("for Name: " + npc.GetComponent(cp.CompName))
	// get id
	npcController.UpdateNpc(editController.SaveNPC())
	npcController.DeleteAllNPC()

	// add a new NPC
	npc, err := service.CreateNPCWithOptions(h.Random, h.Random, creationSupplier)
	if err != nil {
		fmt.Println(err)
		return
	}
	npcController.AddNpc(npc)

	rootDir := "." // Change this to your app’s root directory if needed
	fmt.Println(rootDir)
	//printDirStructure(rootDir, "")
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
