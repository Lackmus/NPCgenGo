package types

type NPCType struct {
	Name  string
	Stats []string
}

func (n NPCType) GetName() string {
	return n.Name
}
