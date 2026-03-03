//go:build !embeddata

package app

import (
	"path/filepath"

	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/product/shared"
)

func resolveCreationConfigLoader(base string) shared.NPCConfigLoader {
	return loader.NewJSONNPCConfigLoader(filepath.Join(base, "creation_data"))
}
