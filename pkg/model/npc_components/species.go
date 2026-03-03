// Description: This file contains the Species struct and its methods.
package npc_components

import (
	helper "github.com/lackmus/npcgengo/internal/platform/helpers"
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

