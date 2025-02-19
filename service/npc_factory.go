package service

// NPCFactory
type NPCFactory struct {
	randomizerService RandomizerService
}

func NewNPCFactory(randomizerService RandomizerService) *NPCFactory {
	return &NPCFactory{
		randomizerService: randomizerService,
	}
}

func (of *NPCFactory) CreateNPCWithOptions(npcType string, faction string) NPCBuilder {
	return *NewBuilder().
		WithType(npcType).
		WithFaction(faction)
}
