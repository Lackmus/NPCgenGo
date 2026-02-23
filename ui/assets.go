package uiassets

import "embed"

//go:embed all:wails/dist all:shared
var assets embed.FS

func Assets() embed.FS {
	return assets
}
