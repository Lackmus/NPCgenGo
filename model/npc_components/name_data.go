package npc_components

import (
	"github.com/lackmus/npcgengo/helper"
)

// NameData : A name has a name, a list of forenames and a list of surnames.
type NameData struct {
	Name      string
	Forenames []string
	Surnames  []string
}

// Name : Return the name of the NameData
func (n NameData) GetName() string {
	return n.Name
}

// NewNameComponent : Create a new component for the name
func (n NameData) NewNameComponent() *NPCComponent {
	forname := helper.GetRandomElement(n.Forenames)
	surname := helper.GetRandomElement(n.Surnames)
	return &NPCComponent{
		Name:  CompName,
		Value: forname + " " + surname,
	}
}
