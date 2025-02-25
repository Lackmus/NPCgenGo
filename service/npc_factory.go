package service

import (
	m "github.com/lackmus/npcgengo/model"
)

func CreateNPCWithOptions(c *NPCCreationSupplier) m.NPC {
	return NewNPCBuilder(c).
		WithRandomType().
		WithRandomSubtype().          // Base subtype component.
		WithRandomSubtypeStats().     // Separate stats component.
		WithRandomSubtypeEquipment(). // Separate equipment component.
		WithRandomSpecies().
		WithName().
		WithRandomFaction().
		WithRandomTrait().
		Build()
}
