package view

import (
	"fmt"

	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

type ConsoleView struct {
}

func NewConsoleView() shared.NPCListViewer {
	return &ConsoleView{}
}

func (v *ConsoleView) Update(npcs []model.NPC) {
	fmt.Println("\n=== NPC Generator Console View ===")
	if len(npcs) == 0 {
		fmt.Println("No NPCs available.")
		return
	}

	fmt.Println("  Available NPCs:")
	for _, npc := range npcs {
		fmt.Printf("\n  %s\n", npc.ShortString())
	}
}

func (v *ConsoleView) Run() {
}
