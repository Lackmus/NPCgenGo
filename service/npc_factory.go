// Description: This file contains the CreateNPCWithOptions function, which creates an NPC using the provided type and faction values. If npctype or faction is empty, the function uses a random value.
package service

import (
	h "github.com/lackmus/npcgengo/helper"
	m "github.com/lackmus/npcgengo/model"
)

// CreateNPCWithOptions creates an NPC using the provided type and faction values.
// If npctype or faction is empty, the function uses a random value.
func CreateNPCWithOptions(npctype string, faction string, c *NPCCreationSupplier, locationID string) (m.NPC, error) {
	builder := NewNPCBuilder(c, locationID)

	// Use provided npctype if available; otherwise, use random.
	if npctype == h.Random {
		builder = builder.WithRandomType()
	} else {
		builder = builder.WithType(npctype)
	}

	// Use provided faction if available; otherwise, use random.
	if faction == h.Random {
		builder = builder.WithRandomFaction()
	} else {
		builder = builder.WithFaction(faction)
	}

	return builder.
		WithRandomSubtype().
		WithRandomSubtypeStats().
		WithRandomSubtypeEquipment().
		WithRandomSpecies().
		WithRandomName().
		WithRandomTrait().
		WithRandomDescription().
		Build()
}
