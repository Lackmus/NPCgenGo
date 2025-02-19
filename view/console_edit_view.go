package view

import (
	"fmt"

	"github.com/lackmus/npcgengo/model"
)

// =============================================================================
// DefaultNpcView
// =============================================================================

// DefaultNpcView is a view that displays an NPC.
type ConsoleEditView struct {
}

// NewDefaultNpcView creates a new instance of DefaultNpcView.
func NewConsoleEditView() *ConsoleEditView {
	return &ConsoleEditView{}
}

// UpdateWithNPC updates the view with a new NPC.
func (v *ConsoleEditView) UpdateNPC(npc model.NPC) {
	fmt.Println("\n=== NPC Generator Console Edit View ===")
	fmt.Println(npc)

}

// Render manually displays an NPC (e.g., for an initial view)
func (v *ConsoleEditView) Render() {

}
