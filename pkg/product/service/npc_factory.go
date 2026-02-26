package service

import (
	h "github.com/lackmus/npcgengo/internal/platform/helpers"
	m "github.com/lackmus/npcgengo/pkg/product/model"
)

func CreateNPCWithOptions(npcType string, faction string, c *NPCCreationSupplier) (m.NPC, error) {
	builder := NewNPCBuilder(c)

	if npcType == h.Random {
		builder = builder.WithRandomType()
	} else {
		builder = builder.WithType(npcType)
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
		Build()
}
