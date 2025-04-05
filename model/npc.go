// Description: This file contains the implementation of the NPC struct, which represents a non-player character built using a set of components.
package model

import (
	"fmt"
	"strings"

	cp "github.com/lackmus/npcgengo/model/npc_components"
)

// NPC represents a non-player character built using a set of components.
type NPC struct {
	ID         string
	LocationID string
	Components map[cp.CompEnum]string
}

// NewNPC creates a new NPC with the given ID.
// It returns a pointer to the new NPC.
func NewNPC() *NPC {
	return &NPC{
		Components: make(map[cp.CompEnum]string),
	}
}

// AddComponent attaches a component to the NPC.
// It adds the component to the NPC's components map.
func (n *NPC) AddComponent(c cp.NPCComponent) {
	if n.Components == nil {
		n.Components = make(map[cp.CompEnum]string)
	}
	n.Components[c.Name] = c.Value
}

// GetComponent returns the value of a component attached to the NPC.
// It returns the value of the component with the given name.
func (n *NPC) GetComponent(name cp.CompEnum) string {
	if n.HasComponent(name) {
		return n.Components[name]
	}
	return ""
}

// RemoveComponent removes a component from the NPC.
// It removes the component with the given name from the NPC's components map.
func (n *NPC) RemoveComponent(name cp.CompEnum) {
	delete(n.Components, name)
}

// String returns a string representation of the NPC and its components.
// It returns a string representation of the NPC and its components.
func (n *NPC) String() string {
	var sb strings.Builder
	// for each component, append the name and value
	for i := range cp.CompEnumValues() {
		c := cp.CompEnum(i + 1)
		if comp, ok := n.Components[c]; ok {
			sb.WriteString(fmt.Sprintf("\n  %s: %s", c, comp))
		}
	}
	return sb.String()
}

// ShortString returns a short string representation of the NPC and its components.
// It returns a short string representation of the NPC and its components.
func (n *NPC) ShortString() string {
	var sb strings.Builder
	// for each component, append the name and value
	for i := range 5 {
		c := cp.CompEnum(i + 1)
		if comp, ok := n.Components[c]; ok {
			sb.WriteString(fmt.Sprintf("%s: [%s] ", c, comp))
		}
	}
	return sb.String()
}

// HasComponent checks if the NPC has a component with the given name.
// It returns true if the NPC has a component with the given name.
func (n *NPC) HasComponent(name cp.CompEnum) bool {
	_, ok := n.Components[name]
	return ok
}
