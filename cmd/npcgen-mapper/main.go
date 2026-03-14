package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lackmus/npcgengo/internal/platform/loader"
	"github.com/lackmus/npcgengo/pkg/mapper"
	"github.com/lackmus/npcgengo/pkg/service"
)

func main() {
	dataDir := flag.String("data-dir", "", "path to data or creation_data directory (overrides NPCGEN_DATA env)")
	inputPath := flag.String("in", "", "path to a JSON file containing mapper.NPCInput (defaults to stdin)")
	pretty := flag.Bool("pretty", true, "pretty-print JSON output")
	flag.Parse()

	creationDir, err := resolveCreationDataDir(*dataDir)
	if err != nil {
		log.Fatal(err)
	}

	rawInput, err := readInput(*inputPath)
	if err != nil {
		log.Fatal(err)
	}

	input, err := decodeNPCInput(rawInput)
	if err != nil {
		log.Fatal(err)
	}

	configLoader := loader.NewJSONNPCConfigLoader(creationDir)
	creationSupplier, err := service.NewNPCCreationSupplier(configLoader)
	if err != nil {
		log.Fatal(err)
	}

	creationData, err := service.NewCreationDataService(context.Background(), configLoader)
	if err != nil {
		log.Fatal(err)
	}

	builder := service.NewNPCBuilder(creationSupplier)
	npc, err := mapper.ToModelNPC(input, builder)
	if err != nil {
		log.Fatal(err)
	}

	validator := service.NewNPCValidationService(creationData)
	if err := validator.ValidateNPC(npc); err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	if *pretty {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(mapper.ToNPCInput(npc)); err != nil {
		log.Fatal(err)
	}
}

func readInput(inputPath string) ([]byte, error) {
	if strings.TrimSpace(inputPath) == "" {
		return io.ReadAll(os.Stdin)
	}

	content, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("read input file %q: %w", inputPath, err)
	}
	return content, nil
}

func decodeNPCInput(raw []byte) (mapper.NPCInput, error) {
	var input mapper.NPCInput

	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return mapper.NPCInput{}, fmt.Errorf("decode NPCInput JSON: %w", err)
	}

	return input, nil
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
