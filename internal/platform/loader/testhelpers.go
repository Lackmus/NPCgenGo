package loader

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// CreateSampleCreationData writes minimal JSON fixtures into dir to satisfy loader tests.
func CreateSampleCreationData(baseDir string) error {
	// directories
	dirs := []string{
		filepath.Join(baseDir, "factiondata"),
		filepath.Join(baseDir, "speciesdata"),
		filepath.Join(baseDir, "traitdata"),
		filepath.Join(baseDir, "namedata"),
		filepath.Join(baseDir, "npctypedata", "civilian"),
		filepath.Join(baseDir, "npctypedata", "military"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}

	// faction
	faction := map[string]any{
		"Name":        "someFactionID",
		"SpeciesList": []string{"someSpeciesID"},
	}
	if err := writeJSON(filepath.Join(baseDir, "factiondata", "someFactionID.json"), faction); err != nil {
		return err
	}

	// species
	species := map[string]any{
		"Name":       "someSpeciesID",
		"NameSource": "someNameID",
	}
	if err := writeJSON(filepath.Join(baseDir, "speciesdata", "someSpeciesID.json"), species); err != nil {
		return err
	}

	// trait
	trait := map[string]any{
		"Name":    "someTraitID",
		"Opposes": "",
	}
	if err := writeJSON(filepath.Join(baseDir, "traitdata", "someTraitID.json"), trait); err != nil {
		return err
	}

	// name data
	name := map[string]any{
		"Name":      "someNameID",
		"Forenames": []string{"Alice"},
		"Surnames":  []string{"Smith"},
	}
	if err := writeJSON(filepath.Join(baseDir, "namedata", "someNameID.json"), name); err != nil {
		return err
	}

	// civilian subtype
	civ := map[string]any{
		"Name":             "someCivilianSubtypeID",
		"NpcTypeName":      "Civilian",
		"Stats":            []string{"Str", "Dex"},
		"EquipmentOptions": map[string][]string{"Weapon": {"Fists"}},
	}
	if err := writeJSON(filepath.Join(baseDir, "npctypedata", "civilian", "someCivilianSubtypeID.json"), civ); err != nil {
		return err
	}

	// military subtype
	mil := map[string]any{
		"Name":             "someMilitarySubtypeID",
		"NpcTypeName":      "Military",
		"Stats":            []string{"Str", "Dex"},
		"EquipmentOptions": map[string][]string{"Weapon": {"Sword"}},
	}
	if err := writeJSON(filepath.Join(baseDir, "npctypedata", "military", "someMilitarySubtypeID.json"), mil); err != nil {
		return err
	}

	return nil
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
