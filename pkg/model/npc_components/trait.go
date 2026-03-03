package npc_components

type Trait struct {
	Name    string
	Opposes string
}

func (t Trait) GetName() string {
	return t.Name
}

func (t Trait) GetDisplayName() string {
	if t.Opposes == "" {
		return t.Name
	}
	return t.Name + "\nOpposes: " + t.Opposes
}
