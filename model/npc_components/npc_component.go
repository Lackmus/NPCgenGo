package npc_components

type NPCComponent struct {
	Name  CompEnum
	Value string
}

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
	CompDrive       CompEnum = 7
	CompStats       CompEnum = 8
	CompAbilities   CompEnum = 9
	CompItems       CompEnum = 10
	CompDescription CompEnum = 11
)

// compenumvalues
func CompEnumValues() []CompEnum {
	return []CompEnum{
		CompName,
		CompType,
		CompSubtype,
		CompSpecies,
		CompFaction,
		CompTrait,
		CompDrive,
		CompStats,
		CompAbilities,
		CompItems,
		CompDescription,
	}
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
	case CompDescription:
		return "Description"
	default:
		return "UNKNOWN_COMPONENT"
	}
}
