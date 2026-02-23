// Description: This package contains the civilian type.
package types

type Civilian struct {
	NPCType
}

var civilianInstance *Civilian

func GetCivilianInstance() *Civilian {
	if civilianInstance == nil {
		civilianInstance = &Civilian{
			NPCType: NPCType{
				Name:  "Civilian",
				Stats: []string{"health", "speed", "strength"},
			},
		}
	}
	return civilianInstance
}
