package model

import (
	"fmt"
	"strings"

	cp "github.com/lackmus/npcgengo/pkg/model/npc_components"
)

const CurrentNPCSchemaVersion = 1

type NPC struct {
	SchemaVersion int
	ID            string
	Components    map[cp.CompEnum]string
}

func NewNPC() *NPC {
	return &NPC{
		SchemaVersion: CurrentNPCSchemaVersion,
		Components:    make(map[cp.CompEnum]string),
	}
}

func (n *NPC) EnsureSchemaVersion() {
	if n.SchemaVersion <= 0 {
		n.SchemaVersion = CurrentNPCSchemaVersion
	}
}

func (n *NPC) AddComponent(c cp.NPCComponent) {
	if n.Components == nil {
		n.Components = make(map[cp.CompEnum]string)
	}
	n.Components[c.Name] = c.Value
}

func (n *NPC) SetComponent(name cp.CompEnum, value string) {
	n.AddComponent(cp.NewComponent(name, value))
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

func (n *NPC) Name() string {
	return n.GetComponent(cp.CompName)
}

func (n *NPC) SetName(value string) {
	n.SetComponent(cp.CompName, value)
}

func (n *NPC) Type() string {
	return n.GetComponent(cp.CompType)
}

func (n *NPC) SetType(value string) {
	n.SetComponent(cp.CompType, value)
}

func (n *NPC) Subtype() string {
	return n.GetComponent(cp.CompSubtype)
}

func (n *NPC) SetSubtype(value string) {
	n.SetComponent(cp.CompSubtype, value)
}

func (n *NPC) Species() string {
	return n.GetComponent(cp.CompSpecies)
}

func (n *NPC) SetSpecies(value string) {
	n.SetComponent(cp.CompSpecies, value)
}

func (n *NPC) Faction() string {
	return n.GetComponent(cp.CompFaction)
}

func (n *NPC) SetFaction(value string) {
	n.SetComponent(cp.CompFaction, value)
}

func (n *NPC) Trait() string {
	return n.GetComponent(cp.CompTrait)
}

func (n *NPC) SetTrait(value string) {
	n.SetComponent(cp.CompTrait, value)
}

func (n *NPC) Stats() string {
	return n.GetComponent(cp.CompStats)
}

func (n *NPC) SetStats(value string) {
	n.SetComponent(cp.CompStats, value)
}

func (n *NPC) Items() string {
	return n.GetComponent(cp.CompItems)
}

func (n *NPC) SetItems(value string) {
	n.SetComponent(cp.CompItems, value)
}

func (n *NPC) Notes() string {
	return n.GetComponent(cp.CompNotes)
}

func (n *NPC) SetNotes(value string) {
	n.SetComponent(cp.CompNotes, value)
}
