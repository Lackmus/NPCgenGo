package npc_components

type Species struct {
	Name       string
	NameSource string
}

func (s Species) GetName() string {
	return s.Name
}
