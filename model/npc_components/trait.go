package npc_components

// Trait represents a trait of an NPC
type Trait struct {
	Name        string
	Description string
	Opposes     string
}

// Name returns the name of the trait
func (t Trait) GetName() string {
	return t.Name
}

// NewTraitComponent creates a new component for the trait
func (t Trait) NewTraitComponent() *Component {
	return &Component{
		Name:  CompTrait,
		Value: t.Name + ": " + t.Description + " Opposes: " + t.Opposes,
	}
}
