// Description: This file contains the Species struct and its methods.
package npc_components

import (
	"github.com/lackmus/npcgengo/helper"
)

type Species struct {
	Name        string
	NameSource  string
	Description []string
}

func (s Species) GetName() string {
	return s.Name
}

func (s Species) GetDescription() string {
	return helper.GetRandomElement(s.Description)
}
