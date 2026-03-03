package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/service"
)

func main() {
	dataDir := flag.String("data-dir", "", "path to data or creation_data directory (overrides NPCGEN_DATA env)")
	flag.Parse()

	creationDir, err := resolveCreationDataDir(*dataDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := service.ValidateCreationData(context.Background(), loader.NewJSONNPCConfigLoader(creationDir)); err != nil {
		fmt.Fprintln(os.Stderr, "creation data validation failed:")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("creation data validation passed")
}

func resolveCreationDataDir(dataDirFlag string) (string, error) {
	base := strings.TrimSpace(dataDirFlag)
	if base == "" {
		base = strings.TrimSpace(os.Getenv("NPCGEN_DATA"))
	}
	if base == "" {
		base = "data"
	}

	if hasFactionDir(base) {
		return base, nil
	}

	candidate := filepath.Join(base, "creation_data")
	if hasFactionDir(candidate) {
		return candidate, nil
	}

	return "", fmt.Errorf("unable to locate creation data under %q or %q", base, candidate)
}

func hasFactionDir(path string) bool {
	info, err := os.Stat(filepath.Join(path, "factiondata"))
	if err != nil {
		return false
	}
	return info.IsDir()
}
