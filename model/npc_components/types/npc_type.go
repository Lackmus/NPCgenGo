package types

type NPCType struct {
	Name        string
	Description string
	Stats       []string
}

func (n NPCType) GetName() string {
	return n.Name
}

func (n NPCType) GetDescription() string {
	return n.Description
}
