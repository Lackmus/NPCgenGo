package model

// Trait represents a trait of an NPC
type Trait struct {
	Name        string
	Description string
	Opposes     string
}

// print the struct in a human readable format
func (t Trait) String() string {
	return t.GetName() + ": " + t.Description + " Opposes: " + t.Opposes
}

// Name returns the name of the trait
func (t Trait) GetName() string {
	return t.Name
}
