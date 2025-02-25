package service

import (
	"fmt"

	m "github.com/lackmus/npcgengo/model"
)

// CreateNPCWithOptions creates an NPC using the provided type and faction values.
// If npctype or faction is empty, the function uses a random value.
func CreateNPCWithOptions(npctype string, faction string, c *NPCCreationSupplier) (m.NPC, error) {
	builder := NewNPCBuilder(c)

	// Use provided npctype if available; otherwise, use random.
	if npctype == "" {
		builder = builder.WithRandomType()
	} else {
		builder = builder.WithType(npctype)
	}

	// Use provided faction if available; otherwise, use random.
	if faction == "" {
		builder = builder.WithRandomFaction()
	} else {
		builder = builder.WithFaction(faction)
	}

	npc, err := builder.
		WithRandomSubtype().
		WithRandomSubtypeStats().
		WithRandomSubtypeEquipment().
		WithRandomSpecies().
		WithRandomName().
		WithRandomTrait().
		Build()
	if err != nil {
		return m.NPC{}, fmt.Errorf("error creating NPC: %w", err)
	}
	return npc, nil
}
