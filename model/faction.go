package model

// Faction : A faction in the game has a name, a lsit of strings and a description.
type Faction struct {
	Name        string
	SpeciesList []string
	Description []string
}
