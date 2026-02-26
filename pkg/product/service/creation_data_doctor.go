package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lackmus/npcgengo/pkg/product/shared"
)

// ValidateCreationData checks creation data integrity across maps and cross-references.
// Returns nil when data is valid, otherwise an aggregated error describing all issues.
func ValidateCreationData(ctx context.Context, npcConfigLoader shared.NPCConfigLoader) error {
	creationData, err := NewCreationDataService(ctx, npcConfigLoader)
	if err != nil {
		return fmt.Errorf("failed to initialize creation data: %w", err)
	}

	var validationErrors []error

	speciesMap := creationData.GetSpeciesMap()
	for factionName, faction := range creationData.GetFactionMap() {
		if len(faction.SpeciesList) == 0 {
			validationErrors = append(validationErrors, fmt.Errorf("faction %q has no species", factionName))
		}
		for _, speciesID := range faction.SpeciesList {
			if _, ok := speciesMap[speciesID]; !ok {
				validationErrors = append(validationErrors, fmt.Errorf("faction %q references unknown species %q", factionName, speciesID))
			}
		}
	}

	for speciesID, species := range speciesMap {
		nameSource := strings.TrimSpace(species.NameSource)
		if nameSource == "" {
			validationErrors = append(validationErrors, fmt.Errorf("species %q has empty name source", speciesID))
			continue
		}
		if _, err := creationData.GetNameData(nameSource); err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("species %q references unknown name source %q", speciesID, nameSource))
		}
	}

	for traitID, trait := range creationData.GetTraitMap() {
		opposes := strings.TrimSpace(trait.Opposes)
		if opposes == "" {
			continue
		}
		if opposes == traitID {
			validationErrors = append(validationErrors, fmt.Errorf("trait %q cannot oppose itself", traitID))
			continue
		}
		if _, err := creationData.GetTraitData(opposes); err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("trait %q references unknown opposing trait %q", traitID, opposes))
		}
	}

	for npcType, subtypeIDs := range creationData.GetNpcSubtypeForTypeMap() {
		if len(subtypeIDs) == 0 {
			validationErrors = append(validationErrors, fmt.Errorf("npc type %q has no subtypes", npcType))
			continue
		}
		for _, subtypeID := range subtypeIDs {
			subtype, subtypeErr := creationData.GetNpcSubtypeData(subtypeID)
			if subtypeErr != nil {
				validationErrors = append(validationErrors, fmt.Errorf("npc type %q references unknown subtype %q", npcType, subtypeID))
				continue
			}
			if subtypeType := strings.TrimSpace(subtype.NpcTypeName); subtypeType != "" && subtypeType != npcType {
				validationErrors = append(validationErrors, fmt.Errorf("subtype %q has type %q but is listed under %q", subtypeID, subtypeType, npcType))
			}
			if len(subtype.Stats) == 0 {
				validationErrors = append(validationErrors, fmt.Errorf("subtype %q has no stats", subtypeID))
			}
			if len(subtype.EquipmentOptions) == 0 {
				validationErrors = append(validationErrors, fmt.Errorf("subtype %q has no equipment options", subtypeID))
			}
			for equipmentSlot, options := range subtype.EquipmentOptions {
				if len(options) == 0 {
					validationErrors = append(validationErrors, fmt.Errorf("subtype %q has empty equipment option list for %q", subtypeID, equipmentSlot))
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return errors.Join(validationErrors...)
	}
	return nil
}
