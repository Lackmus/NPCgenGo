// Description: Defines the Caravan type.
package types

type Caravan struct {
	NPCType
}

var caravanInstance *Caravan

func GetCaravanInstance() *Caravan {
	if caravanInstance == nil {
		caravanInstance = &Caravan{
			NPCType: NPCType{
				Name:        "Caravan",
				Description: "A caravan",
				Stats:       []string{"health", "speed", "strength"},
			},
		}
	}
	return caravanInstance
}

