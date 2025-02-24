package types

//"github.com/lackmus/npcgengo/movementstrategy"

type Civilian struct {
	NPCType
}

// GetCivilianInstance returns a singleton instance of Civilian.
var civilianInstance *Civilian

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
