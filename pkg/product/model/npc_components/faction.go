package npc_components

type Faction struct {
	Name        string
	SpeciesList []string
}

func (f Faction) GetName() string {
	return f.Name
}
