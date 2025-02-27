// Description: Define the NPCType struct type.
package types

// NPCType is a type of NPC.
// It has a name, description, and stats.
type NPCType struct {
	Name        string
	Description string
	Stats       []string
}

// GetName returns the name of the NPCType.
// It returns the name of the NPCType.
func (n NPCType) GetName() string {
	return n.Name
}

// GetStats returns the stats of the NPCType.
// It returns the stats of the NPCType.
func (n NPCType) GetDescription() string {
	return n.Description
}
