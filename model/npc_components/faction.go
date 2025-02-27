package npc_components

// Faction : A faction in the game has a name, a lsit of strings and a description.
type Faction struct {
	Name        string
	SpeciesList []string
	Description []string
}

func (f Faction) GetName() string {
	return f.Name
}
