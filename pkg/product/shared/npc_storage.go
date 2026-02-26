package shared

import (
	"context"

	"github.com/lackmus/npcgengo/pkg/product/model"
)

type NPCStorage interface {

	// LoadNPC loads an NPC from the storage by its ID.
	// It returns the NPC and an error if the NPC cannot be loaded.
	LoadNPC(ctx context.Context, id string) (model.NPC, error)

	// LoadAllNPCs loads all NPCs from the storage.
	// It returns a map of NPC IDs to NPCs and an error if the NPCs cannot be loaded.
	LoadAllNPCs(ctx context.Context) (map[string]model.NPC, error)

	// SaveAllNPCs saves all NPCs to the storage.
	// It takes a map of NPC IDs to NPCs and returns an error if the NPCs cannot be saved.
	SaveAllNPCs(ctx context.Context, npcs map[string]model.NPC) error

	// SaveNPC saves a single NPC to the storage.
	// It takes the NPC as a parameter and returns an error if the NPC cannot be saved.
	SaveNPC(ctx context.Context, npc model.NPC) error

	// DeleteNPC deletes an NPC from the storage by its ID.
	// It returns an error if the NPC cannot be deleted.
	DeleteNPC(ctx context.Context, id string) error

	// DeleteAllNPCs deletes all NPCs from the storage.
	// It returns an error if the NPCs cannot be deleted.
	DeleteAllNPCs(ctx context.Context) error
}
