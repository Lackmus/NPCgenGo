// Description: This file contains the Trait struct which represents a trait of an NPC.
package npc_components

type Trait struct {
	Name        string
	Description string
	Opposes     string
}

func (t Trait) GetName() string {
	return t.Name + "\nOpposes: " + t.Opposes
}

func (t Trait) GetDescription() string {
	return t.Description
}
