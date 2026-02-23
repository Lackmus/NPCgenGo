package service

import (
	"fmt"
	"strings"

	"github.com/lackmus/npcgengo/pkg/product/model"
	cp "github.com/lackmus/npcgengo/pkg/product/model/npc_components"
)

type NPCValidationService struct {
	creationData *CreationDataService
}

func NewNPCValidationService(creationData *CreationDataService) *NPCValidationService {
	return &NPCValidationService{creationData: creationData}
}

func (v *NPCValidationService) ValidateNPC(npc model.NPC) error {
	if v == nil || v.creationData == nil {
		return nil
	}

	npcType := strings.TrimSpace(npc.GetComponent(cp.CompType))
	if npcType != "" {
		if _, err := v.creationData.GetNpcTypeData(npcType); err != nil {
			return fmt.Errorf("invalid type %q: %w", npcType, err)
		}
	}

	subtype := strings.TrimSpace(npc.GetComponent(cp.CompSubtype))
	if subtype != "" {
		subtypeData, err := v.creationData.GetNpcSubtypeData(subtype)
		if err != nil {
			return fmt.Errorf("invalid subtype %q: %w", subtype, err)
		}
		if npcType != "" && strings.TrimSpace(subtypeData.NpcTypeName) != "" && subtypeData.NpcTypeName != npcType {
			return fmt.Errorf("subtype %q does not belong to type %q", subtype, npcType)
		}
	}

	species := strings.TrimSpace(npc.GetComponent(cp.CompSpecies))
	if species != "" {
		if _, err := v.creationData.GetSpeciesData(species); err != nil {
			return fmt.Errorf("invalid species %q: %w", species, err)
		}
	}

	faction := strings.TrimSpace(npc.GetComponent(cp.CompFaction))
	if faction != "" {
		if _, err := v.creationData.GetFactionData(faction); err != nil {
			return fmt.Errorf("invalid faction %q: %w", faction, err)
		}
	}

	traitValue := strings.TrimSpace(npc.GetComponent(cp.CompTrait))
	if traitValue != "" {
		for _, rawTrait := range strings.Split(traitValue, ",") {
			trait := strings.TrimSpace(rawTrait)
			if trait == "" {
				continue
			}
			if _, err := v.creationData.GetTraitData(trait); err != nil {
				return fmt.Errorf("invalid trait %q: %w", trait, err)
			}
		}
	}

	return nil
}
