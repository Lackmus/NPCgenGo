package types

//"github.com/lackmus/npcgengo/movementstrategy"

// Military is a type of npc.
type Military struct {
	NPCType
}

// militaryInstance is a singleton instance of Military.
var militaryInstance *Military

// GetMilitaryInstance returns a singleton instance of Military.
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
