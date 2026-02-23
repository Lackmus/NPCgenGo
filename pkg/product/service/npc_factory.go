// Description: This file contains the CreateNPCWithOptions function, which creates an NPC using the provided type and faction values. If npctype or faction is empty, the function uses a random value.
package service

import (
	h "github.com/lackmus/npcgengo/internal/platform/helpers"
	m "github.com/lackmus/npcgengo/pkg/product/model"
)

func CreateNPCWithOptions(npctype string, faction string, c *NPCCreationSupplier, locationID string) (m.NPC, error) {
	builder := NewNPCBuilder(c, locationID)

	if npctype == h.Random {
		builder = builder.WithRandomType()
	} else {
		builder = builder.WithType(npctype)
	}

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

