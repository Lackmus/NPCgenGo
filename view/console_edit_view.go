package view

import (
	"fmt"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

// =============================================================================
// DefaultNpcView
// =============================================================================

// DefaultNpcView is a view that displays an NPC.
type ConsoleEditView struct {
}

// NewDefaultNpcView creates a new instance of DefaultNpcView.
func NewConsoleEditView() shared.NPCEditViewer {
	return &ConsoleEditView{}
}

// UpdateWithNPC updates the view with a new NPC.
func (v *ConsoleEditView) UpdateNPC(npc model.NPC) {
	fmt.Println("\n=== NPC Generator Console Edit View ===")
	fmt.Printf("  %s\n", npc.String())

}

// Render manually displays an NPC (e.g., for an initial view)
func (v *ConsoleEditView) Render() {

}

// uodate field
func (v *ConsoleEditView) UpdateField(field string, value any) {
	fmt.Println("Field: ", field, " Value: ", value)
}
