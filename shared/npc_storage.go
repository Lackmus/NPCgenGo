package shared

import (
	"context"

	"github.com/lackmus/npcgengo/model"
)

type NPCStorage interface {

	// LoadNPC loads an NPC from the storage by its ID.
	// It returns the NPC and an error if the NPC cannot be loaded.
	LoadNPC(ctx context.Context, id string) (model.NPC, error)

	// LoadAllNPC loads all NPCs from the storage.
	// It returns a map of NPC IDs to NPCs and an error if the NPCs cannot be loaded.
	LoadAllNPC(ctx context.Context) (map[string]model.NPC, error)

	// SaveAllNPC saves all NPCs to the storage.
	// It takes a map of NPC IDs to NPCs and returns an error if the NPCs cannot be saved.
	SaveAllNPC(ctx context.Context, npcs map[string]model.NPC) error

	// SaveNPC saves a single NPC to the storage.
	// It takes the NPC as a parameter and returns an error if the NPC cannot be saved.
	SaveNPC(ctx context.Context, npc model.NPC) error

	// DeleteNPC deletes an NPC from the storage by its ID.
	// It returns an error if the NPC cannot be deleted.
	DeleteNPC(ctx context.Context, id string) error

	// DeleteAllNPC deletes all NPCs from the storage.
	// It returns an error if the NPCs cannot be deleted.
	DeleteAllNPC(ctx context.Context) error
}
