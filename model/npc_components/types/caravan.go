package types

// Caravan is a type of npc.
type Caravan struct {
	NPCType
}

// caravanInstance is a singleton instance of Caravan.
var caravanInstance *Caravan

// GetCaravanInstance returns a singleton instance of Caravan.
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
