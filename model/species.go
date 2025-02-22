package model

// Species : A species in the game has a name and a name source.
type Species struct {
	Name       string
	NameSource string
}

// Name : Return the name of the Species
func (s Species) GetName() string {
	return s.Name
}
