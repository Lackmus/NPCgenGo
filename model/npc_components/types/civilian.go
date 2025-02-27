// Description: This package contains the civilian type.
package types

//"github.com/lackmus/npcgengo/movementstrategy"

// Civilian is a type of npc.
// It represents a civilian npc.
type Civilian struct {
	NPCType
}

var civilianInstance *Civilian

// GetCivilianInstance returns the Civilian instance.
// It returns the Civilian instance.
func GetCivilianInstance() *Civilian {
	if civilianInstance == nil {
		civilianInstance = &Civilian{
			NPCType: NPCType{
				Name:        "Civilian",
				Description: "A regular civilian",
				Stats:       []string{"health", "speed", "strength"},
			},
		}
	}
	return civilianInstance
}
