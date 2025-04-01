// Description: This file contains the NPCComponent struct and its methods.
package npc_components

// NPCComponent represents a component of an NPC.
// It represents a component of an NPC.
type NPCComponent struct {
	Name  CompEnum
	Value string
}

// NewComponent creates a new NPCComponent.
// It returns a new NPCComponent with the given name and value.
func NewComponent(name CompEnum, value string) NPCComponent {
	return NPCComponent{
		Name:  name,
		Value: value,
	}
}

type CompEnum int

const (
	CompName        CompEnum = 1
	CompType        CompEnum = 2
	CompSubtype     CompEnum = 3
	CompSpecies     CompEnum = 4
	CompFaction     CompEnum = 5
	CompTrait       CompEnum = 6
	CompStats       CompEnum = 7
	CompItems       CompEnum = 8
	CompDescription CompEnum = 9
)

// CompEnumValues returns a slice of all the CompEnum values.
// It returns a slice of all the CompEnum values.
func CompEnumValues() []CompEnum {
	return []CompEnum{
		CompName,
		CompType,
		CompSubtype,
		CompSpecies,
		CompFaction,
		CompTrait,
		CompStats,
		CompItems,
		CompDescription,
	}
}

// String returns the string representation of the CompEnum.
// It returns the string representation of the CompEnum.
func (c CompEnum) String() string {
	switch c {
	case CompName:
		return "Name"
	case CompType:
		return "Type"
	case CompSubtype:
		return "Subtype"
	case CompSpecies:
		return "Species"
	case CompFaction:
		return "Faction"
	case CompTrait:
		return "Trait"
	case CompStats:
		return "Stats"
	case CompItems:
		return "Items"
	case CompDescription:
		return "Description"
	default:
		return "UNKNOWN_COMPONENT"
	}
}
