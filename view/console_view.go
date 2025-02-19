package view

import (
	"fmt"

	"github.com/lackmus/npcgengo/model"
)

// ConsoleView is an observer that prints NPC updates to the console
type ConsoleView struct {
}

// NewConsoleView creates a ConsoleView and registers it as an observer
func NewConsoleView() *ConsoleView {
	return &ConsoleView{}
}

// Update is triggered when NPCs change
func (v *ConsoleView) Update(npcs []model.NPC) {
	fmt.Println("\n=== NPC Generator Console View ===")
	if len(npcs) == 0 {
		fmt.Println("No NPCs available.")
		return
	}

	fmt.Println("Available NPCs:")
	for i, npc := range npcs {
		fmt.Printf("[%d] %s (Type: %s)\n", i+1, npc.Name(), npc.NPCType())
	}
}

// Render manually displays NPCs (e.g., for an initial view)
func (v *ConsoleView) Render() {
}
