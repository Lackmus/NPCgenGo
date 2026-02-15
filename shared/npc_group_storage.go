package shared

import "github.com/lackmus/npcgengo/model"

// NPCGroupStorage interface for storing and retrieving NPC groups
type NPCGroupStorage interface {
	// LoadNPCGroup loads an NPC group from storage by its ID
	LoadNPCGroup(id string) (model.NPCGroup, error)

	// LoadAllNPCGroups loads all NPC groups from storage
	LoadAllNPCGroups() (map[string]model.NPCGroup, error)

	// SaveNPCGroup saves an NPC group to storage
	SaveNPCGroup(group model.NPCGroup) error

	// DeleteNPCGroup deletes an NPC group from storage
	DeleteNPCGroup(id string) error

	// DeleteAllNPCGroups deletes all NPC groups from storage
	DeleteAllNPCGroups() error
}

// NPCGroupObserver interface for objects that need to be notified of NPC group changes
type NPCGroupObserver interface {
	// UpdateGroups is called when the list of NPC groups changes
	UpdateGroups(groups []model.NPCGroup)
}
