// Description: Defines the Caravan type.
package types

// Caravan is a type of NPC.
// It has a name, description, and stats.
type Caravan struct {
	NPCType
}

var caravanInstance *Caravan

// GetCaravanInstance returns the Caravan instance.
// It returns the Caravan instance.
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
