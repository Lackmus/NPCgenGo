package service

import (
	h "github.com/lackmus/npcgengo/internal/platform/helpers"
	m "github.com/lackmus/npcgengo/pkg/product/model"
)

// CreateNPCWithOptionsAndSeed creates an NPC with the given type and faction, using the provided seed for randomness.
func CreateNPCWithOptionsAndSeed(npcType string, faction string, seed int64, c *NPCCreationSupplier) (m.NPC, error) {
	// Use the helper function to set the seed for randomness during NPC creation.
	return h.WithSeed(seed, func() (m.NPC, error) {
		// Call the main NPC creation function with the provided options.
		return CreateNPCWithOptions(npcType, faction, c)
	})
}

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
