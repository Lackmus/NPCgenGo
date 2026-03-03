package npc_components

type NPCType struct {
	Name  string
	Stats []string
}

func (n NPCType) GetName() string {
	return n.Name
}
