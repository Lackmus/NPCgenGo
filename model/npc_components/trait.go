package npc_components

// Trait represents a trait of an NPC
type Trait struct {
	Name        string
	Description string
	Opposes     string
}

// Name returns the name of sthe trait
func (t Trait) GetName() string {
	return t.Name + "\n  - Opposes: " + t.Opposes
}

// GetDescription returns the description of the trait
func (t Trait) GetDescription() string {
	return t.Description
}
