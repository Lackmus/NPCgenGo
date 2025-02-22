package types

// NpcSubtype is a type of npc.
type NPCSubtype struct {
	Name             string
	NpcTypeName      string
	Description      string
	Stats            []string
	EquipmentOptions map[string][]string
}

// GetName returns the name of the NPCSubtype.
func (n NPCSubtype) GetName() string {
	return n.Name
}
