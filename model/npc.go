package model

import (
	"fmt"

	cp "github.com/lackmus/npcgengo/model/npc_components"
)

// NPC represents a non-player character built using a set of components.
type NPC struct {
	ID         string
	Components map[cp.CompEnum]string
}

// NewNPC creates a new NPC with the given ID.
func NewNPC(id string) *NPC {
	return &NPC{
		ID:         id,
		Components: make(map[cp.CompEnum]string),
	}
}

// AddComponent attaches a new component to the NPC.
func (n *NPC) AddComponent(c cp.Component) {
	n.Components[c.Name] = c.Value
}

// GetComponent retrieves a component by its name.
// It returns the component and a boolean indicating whether it was found.
func (n *NPC) GetComponent(name cp.CompEnum) string {
	if n.HasComponent(name) {
		return n.Components[name]
	}
	return ""
}

// RemoveComponent detaches a component from the NPC.
func (n *NPC) RemoveComponent(name cp.CompEnum) {
	delete(n.Components, name)
}

// String returns a string representation of the NPC and its components.
func (n *NPC) String() string {
	result := ""
	//for first 5 components of map print key and value. print all components in one line seperated by,trim last comma
	for i := range cp.CompEnumValues() {
		c := cp.CompEnum(i + 1)
		if comp, ok := n.Components[c]; ok {
			result += fmt.Sprintf("\n  %s: [%s]", c, comp)
		}
	}
	return result
}

// shortstring returns a string representation of the NPC and its components. It is shorter than the String() method. if comp = Name Type Subtype faction species
func (n *NPC) ShortString() string {
	result := ""
	//for first 5 components of map print key and value. print all components in one line seperated by,trim last comma
	for i := 1; i < 6; i++ {
		if comp, ok := n.Components[cp.CompEnum(i)]; ok {
			result += fmt.Sprintf("%s: [%s] ", cp.CompEnum(i), comp)
		}
	}
	return result
}

// Hascomponent returns true if the NPC has the component
func (n *NPC) HasComponent(name cp.CompEnum) bool {
	_, ok := n.Components[name]
	return ok
}
