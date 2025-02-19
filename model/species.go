package model

import (
	"fmt"
)

// Species : A species in the game has a name and a name source.
type Species struct {
	Name       string
	NameSource string
}

// print the struct in a human readable format
func (s Species) String() string {
	return fmt.Sprintf("Species: %s", s.Name)
}
