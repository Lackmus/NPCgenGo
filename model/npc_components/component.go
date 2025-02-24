package npc_components

type Component struct {
	Name  CompEnum
	Value string
}

type CompEnum int

const (
	CompName      CompEnum = 1
	CompType      CompEnum = 2
	CompSubtype   CompEnum = 3
	CompSpecies   CompEnum = 4
	CompFaction   CompEnum = 5
	CompTrait     CompEnum = 6
	CompDrive     CompEnum = 7
	CompStats     CompEnum = 8
	CompAbilities CompEnum = 9
	CompItems     CompEnum = 10
)

// compenumvalues
func CompEnumValues() []CompEnum {
	return []CompEnum{CompName, CompType, CompSubtype, CompSpecies, CompFaction, CompTrait, CompDrive, CompStats, CompAbilities, CompItems}
}

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
	case CompDrive:
		return "Drive"
	case CompStats:
		return "Stats"
	case CompAbilities:
		return "Abilities"
	case CompItems:
		return "Items"
	default:
		return "UNKNOWN_COMPONENT"
	}
}
