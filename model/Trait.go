package model

// Trait represents a trait of an NPC
type Trait struct {
	Name        string
	Description string
	Opposes     string
	Stats       map[string]int
}

// print the struct in a human readable format
func (t Trait) String() string {
	return t.Name + ": " + t.Description + " Opposes: " + t.Opposes
}
