// Description: This file contains the ConsoleEditView struct and its methods. This struct is used to display an NPC in the console.
package view

import (
	"fmt"

	"github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components"
	"github.com/lackmus/npcgengo/shared"
)

// DefaultNpcView is a view that displays an NPC.
type ConsoleEditView struct {
}

func NewConsoleEditView() shared.NPCEditViewer {
	return &ConsoleEditView{}
}

// UpdateWithNPC updates the view with a new NPC.
func (v *ConsoleEditView) UpdateNPC(npc model.NPC) {
	fmt.Println("\n=== NPC Generator Console Edit View ===")
	fmt.Printf("  %s\n", npc.String())

}

// Render manually displays an NPC (e.g., for an initial view)
func (v *ConsoleEditView) Run() {

}

// UpdateField updates the field of the NPC with the given value.
// It takes a field of type cp.CompEnum and a value of any type as parameters and returns nothing.
func (v *ConsoleEditView) UpdateField(field cp.CompEnum, value any) {
	fmt.Println("Field: ", field, " Value: ", value)
}

// OnNPCEditError is a new method for error reporting
// It takes an error as a parameter and returns nothing.
func (v *ConsoleEditView) OnNPCEditError(err error) {
	fmt.Println("Error: ", err)
}
