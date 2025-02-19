package types

// NpcSubtype is a type of npc.
type NPCSubtype struct {
	Name             string
	NpcTypeName      string
	Description      string
	Stats            []string
	EquipmentOptions map[string][]string
}
