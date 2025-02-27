// Description: This file contains the Trait struct which represents a trait of an NPC.
package npc_components

// Trait : Represents a trait of an NPC
// It represents a trait of an NPC
type Trait struct {
	Name        string
	Description string
	Opposes     string
}

// GetName returns the name of the trait
// It returns the name of the trait
func (t Trait) GetName() string {
	return t.Name + "\n  - Opposes: " + t.Opposes
}

// GetDescription returns the description of the trait
// It returns the description of the trait
func (t Trait) GetDescription() string {
	return t.Description
}
