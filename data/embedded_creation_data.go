package data

import (
	"embed"
	"io/fs"
)

//go:embed creation_data/**
var creationData embed.FS

// CreationDataFS returns the embedded creation data filesystem.
func CreationDataFS() fs.FS {
	return creationData
}
