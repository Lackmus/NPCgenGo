// Description: Defines the Military type of npc.
package types

type Military struct {
	NPCType
}

var militaryInstance *Military

func GetMilitaryInstance() *Military {
	if militaryInstance == nil {
		militaryInstance = &Military{
			NPCType: NPCType{
				Name:        "Military",
				Description: "A military npc",
				Stats:       []string{"health", "speed", "strength"},
			},
		}
	}
	return militaryInstance
}

