package types

import cp "github.com/lackmus/npcgengo/model/npc_components"

type NPCType struct {
	Name        string
	Description string
	Stats       []string
}

func (n NPCType) GetName() string {
	return n.Name
}

func (n NPCType) NewNPCTypeComponent() *cp.NPCComponent {
	return &cp.NPCComponent{
		Name:  cp.CompType,
		Value: n.Name,
	}
}

func (n NPCType) GetDescription() string {
	return n.Description
}
