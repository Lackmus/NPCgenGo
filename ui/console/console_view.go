package console

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/lackmus/npcgengo/internal/app/controllers"
	"github.com/lackmus/npcgengo/pkg/mapper"
	"github.com/lackmus/npcgengo/pkg/model"
	"github.com/lackmus/npcgengo/pkg/shared"
)

type ConsoleView struct {
	controller *controllers.NPCListController
	lastNPCs   []model.NPC
}

func NewConsoleView(ctrl *controllers.NPCListController) shared.NPCListViewer {
	return &ConsoleView{
		controller: ctrl,
	}
}

func (v *ConsoleView) Update(npcs []model.NPC) {
	fmt.Printf("ConsoleView: received update with %d NPCs\n", len(npcs))
	v.lastNPCs = npcs
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
		fmt.Println("Commands: list | details <id> | random | create | edit <id> | delete <id> | delete-all | quit")
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
		case "details":
			if len(parts) < 2 {
				fmt.Println("usage: details <id>")
				continue
			}
			v.showDetails(parts[1])
		case "random":
			v.controller.CreateRandomNPC()
		case "create":
			v.createNPC(reader)
		case "edit":
			if len(parts) < 2 {
				fmt.Println("usage: edit <id>")
				continue
			}
			v.editNPC(reader, parts[1])
		case "delete":
			if len(parts) < 2 {
				fmt.Println("usage: delete <id>")
				continue
			}
			id := parts[1]
			v.controller.DeleteNPC(id)
		case "delete-all":
			v.controller.DeleteAllNPCs()
		case "quit", "exit":
			fmt.Println("Exiting console view")
			return
		default:
			fmt.Println("unknown command")
		}
	}
}

func (v *ConsoleView) showDetails(id string) {
	npc, err := v.controller.GetNPCByID(id)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println("\n--- NPC Details ---")
	fmt.Printf("Name: %s\n", npc.Name())
	fmt.Printf("Type: %s\n", npc.Type())
	fmt.Printf("Subtype: %s\n", npc.Subtype())
	fmt.Printf("Species: %s\n", npc.Species())
	fmt.Printf("Faction: %s\n", npc.Faction())
	fmt.Printf("Trait: %s\n", npc.Trait())
	fmt.Printf("Stats: %s\n", npc.Stats())
	fmt.Printf("Items: %s\n", npc.Items())
	fmt.Printf("Notes: %s\n", npc.Notes())
}

func (v *ConsoleView) createNPC(reader *bufio.Reader) {
	input := mapper.NPCInput{}
	v.collectNPCInput(reader, &input, nil)

	if !v.confirm(reader, "Save new NPC?") {
		fmt.Println("create cancelled")
		return
	}

	v.saveInput(input)
}

func (v *ConsoleView) editNPC(reader *bufio.Reader, id string) {
	npc, err := v.controller.GetNPCByID(id)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	input := mapper.NPCInput{
		ID:      npc.ID,
		Name:    npc.Name(),
		Type:    npc.Type(),
		Subtype: npc.Subtype(),
		Species: npc.Species(),
		Faction: npc.Faction(),
		Trait:   npc.Trait(),
		Notes:   npc.Notes(),
	}

	v.collectNPCInput(reader, &input, &npc)

	if !v.confirm(reader, "Save changes?") {
		fmt.Println("edit cancelled")
		return
	}

	v.saveInputWithOriginal(input, &npc)
}

func (v *ConsoleView) collectNPCInput(reader *bufio.Reader, input *mapper.NPCInput, existing *model.NPC) {
	options := v.controller.GetCreationOptions()

	input.Type = v.promptOption(reader, "Type", options.NpcTypes, input.Type)
	input.Subtype = v.promptOption(reader, "Subtype", options.NpcSubtypeForTypeMap[input.Type], input.Subtype)
	input.Faction = v.promptOption(reader, "Faction", options.Factions, input.Faction)
	input.Species = v.promptOption(reader, "Species", options.NpcSpeciesForFactionMap[input.Faction], input.Species)
	input.Trait = v.promptOption(reader, "Trait", options.Traits, input.Trait)

	nameDefault := input.Name
	if strings.TrimSpace(nameDefault) == "" {
		if generated, err := v.rollSpeciesName(input.Species); err == nil {
			nameDefault = generated
		}
	}
	input.Name = v.promptText(reader, "Name", nameDefault)
	input.Notes = v.promptRaw(reader, "Notes (optional)", input.Notes)

	_ = existing
}

func (v *ConsoleView) promptOption(reader *bufio.Reader, label string, options []string, current string) string {
	if len(options) > 0 {
		fmt.Printf("%s options: %s\n", label, strings.Join(options, ", "))
	}
	for {
		value := v.promptRaw(reader, label, current)
		if value == "" {
			fmt.Printf("%s is required\n", label)
			continue
		}
		if len(options) > 0 && !slices.Contains(options, value) {
			fmt.Printf("invalid %s: %s\n", strings.ToLower(label), value)
			continue
		}
		return value
	}
}

func (v *ConsoleView) promptText(reader *bufio.Reader, label string, current string) string {
	for {
		value := v.promptRaw(reader, label, current)
		if value == "" {
			fmt.Printf("%s is required\n", label)
			continue
		}
		return value
	}
}

func (v *ConsoleView) promptRaw(reader *bufio.Reader, label string, current string) string {
	current = strings.TrimSpace(current)
	if current != "" {
		fmt.Printf("%s [%s]: ", label, current)
	} else {
		fmt.Printf("%s: ", label)
	}
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return current
	}
	return line
}

func (v *ConsoleView) confirm(reader *bufio.Reader, question string) bool {
	for {
		fmt.Printf("%s (y/n): ", question)
		line, _ := reader.ReadString('\n')
		switch strings.ToLower(strings.TrimSpace(line)) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
	}
}

func (v *ConsoleView) saveInput(input mapper.NPCInput) {
	missing := v.requiredMissing(input)
	if len(missing) > 0 {
		fmt.Printf("cannot save: missing required fields: %s\n", strings.Join(missing, ", "))
		return
	}

	npc, err := mapper.ToModelNPC(input, v.controller.GetNPCBuilder())
	if err != nil {
		fmt.Printf("error creating NPC: %v\n", err)
		return
	}
	if input.ID != "" {
		npc.ID = input.ID
	}

	if err := v.controller.ValidateNPC(npc); err != nil {
		fmt.Printf("cannot save: %v\n", err)
		return
	}

	v.controller.UpdateNPC(npc)
	if input.ID == "" {
		fmt.Println("created NPC")
	} else {
		fmt.Printf("updated NPC %s\n", input.ID)
	}
}

// saveInputWithOriginal saves an NPC allowing unchanged fields to be preserved
func (v *ConsoleView) saveInputWithOriginal(input mapper.NPCInput, original *model.NPC) {
	missing := v.requiredMissing(input)
	if len(missing) > 0 {
		fmt.Printf("cannot save: missing required fields: %s\n", strings.Join(missing, ", "))
		return
	}

	npc, err := mapper.ToModelNPCWithOriginal(input, v.controller.GetNPCBuilder(), original)
	if err != nil {
		fmt.Printf("error creating NPC: %v\n", err)
		return
	}
	if input.ID != "" {
		npc.ID = input.ID
	}

	if err := v.controller.ValidateNPC(npc); err != nil {
		fmt.Printf("cannot save: %v\n", err)
		return
	}

	v.controller.UpdateNPC(npc)
	if input.ID == "" {
		fmt.Println("created NPC")
	} else {
		fmt.Printf("updated NPC %s\n", input.ID)
	}
}

func (v *ConsoleView) requiredMissing(input mapper.NPCInput) []string {
	checks := []struct {
		label string
		value string
	}{
		{label: "name", value: input.Name},
		{label: "type", value: input.Type},
		{label: "subtype", value: input.Subtype},
		{label: "species", value: input.Species},
		{label: "faction", value: input.Faction},
		{label: "trait", value: input.Trait},
	}

	missing := make([]string, 0)
	for _, check := range checks {
		if strings.TrimSpace(check.value) == "" {
			missing = append(missing, check.label)
		}
	}
	return missing
}

func (v *ConsoleView) rollSubtypeFields(subtype string) (string, string, error) {
	return v.controller.GetSubtypeFields(subtype)
}

func (v *ConsoleView) rollSpeciesName(species string) (string, error) {
	return v.controller.GetSpeciesName(species)
}
