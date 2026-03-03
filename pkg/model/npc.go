// Description: This file contains the implementation of the NPC struct, which represents a non-player character built using a set of components.
package model

import (
	"fmt"
	"strings"

	cp "github.com/lackmus/npcgengo/pkg/model/npc_components"
)

type NPC struct {
	ID         string
	LocationID string
	Components map[cp.CompEnum]string
}

func NewNPC() *NPC {
	return &NPC{
		Components: make(map[cp.CompEnum]string),
	}
}

func (n *NPC) AddComponent(c cp.NPCComponent) {
	if n.Components == nil {
		n.Components = make(map[cp.CompEnum]string)
	}
	n.Components[c.Name] = c.Value
}

func (n *NPC) GetComponent(name cp.CompEnum) string {
	if n.HasComponent(name) {
		return n.Components[name]
	}
	return ""
}

func (n *NPC) RemoveComponent(name cp.CompEnum) {
	delete(n.Components, name)
}

func (n *NPC) String() string {
	var sb strings.Builder
	for i := range cp.CompEnumValues() {
		c := cp.CompEnum(i + 1)
		if comp, ok := n.Components[c]; ok {
			sb.WriteString(fmt.Sprintf("\n  %s: %s", c, comp))
		}
	}
	return sb.String()
}

func (n *NPC) ShortString() string {
	var sb strings.Builder
	for i := range 5 {
		c := cp.CompEnum(i + 1)
		if comp, ok := n.Components[c]; ok {
			sb.WriteString(fmt.Sprintf("%s: [%s] ", c, comp))
		}
	}
	return sb.String()
}

func (n *NPC) HasComponent(name cp.CompEnum) bool {
	_, ok := n.Components[name]
	return ok
}

