//go:build embeddata

package app

import (
	"path/filepath"

	npcgendata "github.com/lackmus/npcgengo/data"
	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/shared"
)

func resolveCreationConfigLoader(base string) shared.NPCConfigLoader {
	if hasCreationData(base) {
		return loader.NewJSONNPCConfigLoader(filepath.Join(base, "creation_data"))
	}

	return loader.NewFSNPCConfigLoader(npcgendata.CreationDataFS(), "creation_data")
}
