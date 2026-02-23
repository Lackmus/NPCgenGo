package npcgengo

import (
	"github.com/lackmus/npcgengo/internal/app"
)

type NPCGen = app.NPCGen

// NewNPCGen initializes a new NPCGen instance using the default on-disk `data` directory.
func NewNPCGen() (*NPCGen, error) {
	return app.NewNPCGen()
}

// NewNPCGenWithDataDir initializes a new NPCGen instance using the provided data directory.
// If dataDir is empty, it defaults to the repository-relative "data" directory.
func NewNPCGenWithDataDir(dataDir string) (*NPCGen, error) {
	return app.NewNPCGenWithDataDir(dataDir)
}
