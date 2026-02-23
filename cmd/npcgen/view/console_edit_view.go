// Description: This file contains the ConsoleEditView struct and its methods. This struct is used to display an NPC in the console.
package view

import (
	"fmt"

	"github.com/lackmus/npcgengo/pkg/product/model"
	cp "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

type ConsoleEditView struct {
}

func NewConsoleEditView() shared.NPCEditViewer {
	return &ConsoleEditView{}
}

func (v *ConsoleEditView) UpdateNPC(npc model.NPC) {
	fmt.Println("\n=== NPC Generator Console Edit View ===")
	fmt.Printf("  %s\n", npc.String())

}

func (v *ConsoleEditView) Run() {

}

func (v *ConsoleEditView) UpdateField(field cp.CompEnum, value any) {
	fmt.Println("Field: ", field, " Value: ", value)
}

func (v *ConsoleEditView) OnNPCEditError(err error) {
	fmt.Println("Error: ", err)
}

