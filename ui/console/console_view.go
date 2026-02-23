package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/pkg/product/model"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

type ConsoleView struct {
	controller *controllers.NPCListController
	lastNpcs   []model.NPC
}

func NewConsoleView(ctrl *controllers.NPCListController) shared.NPCListViewer {
	return &ConsoleView{controller: ctrl}
}

func (v *ConsoleView) Update(npcs []model.NPC) {
	fmt.Printf("ConsoleView: received update with %d NPCs\n", len(npcs))
	v.lastNpcs = npcs
	fmt.Println("\n=== NPC Generator Console View ===")
	if len(npcs) == 0 {
		fmt.Println("No NPCs available.")
		return
	}

	fmt.Println("Available NPCs:")
	for _, npc := range npcs {
		fmt.Printf("  %s\n", npc.ShortString())
	}
}

func (v *ConsoleView) Run() {
	if v.controller == nil {
		fmt.Println("ConsoleView: controller is not set")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println()
		fmt.Println("Commands: list | random | delete <id> | delete-all | quit")
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := strings.ToLower(parts[0])

		switch cmd {
		case "list":
			v.controller.NotifyObservers()
		case "random":
			v.controller.CreateRandomNPC()
		case "delete":
			if len(parts) < 2 {
				fmt.Println("usage: delete <id>")
				continue
			}
			id := parts[1]
			v.controller.DeleteNPC(id)
		case "delete-all":
			v.controller.DeleteAllNPC()
		case "quit", "exit":
			fmt.Println("Exiting console view")
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
