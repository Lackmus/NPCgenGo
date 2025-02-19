// model/npc.go
package model

// NPC represents an immutable non-player character.
type NPC struct {
	id          string
	name        string
	faction     string
	species     string
	npcType     string
	npcSubtype  string
	trait       string
	drive       string
	stats       map[string]int
	items       map[string]string
	abilities   map[string]string
	description string
}

// Getters to access the fields
func (n NPC) ID() string         { return n.id }
func (n NPC) Name() string       { return n.name }
func (n NPC) Faction() string    { return n.faction }
func (n NPC) Species() string    { return n.species }
func (n NPC) NpcType() string    { return n.npcType }
func (n NPC) NpcSubtype() string { return n.npcSubtype }
func (n NPC) Trait() string      { return n.trait }
func (n NPC) Drive() string      { return n.drive }
func (n NPC) Stats() map[string]int {
	cp := make(map[string]int)
	for k, v := range n.stats {
		cp[k] = v
	}
	return cp
}
func (n NPC) Items() map[string]string {
	cp := make(map[string]string)
	for k, v := range n.items {
		cp[k] = v
	}
	return cp
}
func (n NPC) Abilities() map[string]string {
	cp := make(map[string]string)
	for k, v := range n.abilities {
		cp[k] = v
	}
	return cp
}
func (n NPC) Description() string { return n.description }

// NewNPC is the constructor to create an immutable NPC.
func NewNPC(
	id, name, faction, species, npcType, npcSubtype, trait, drive, description string,
	stats map[string]int,
	items map[string]string,
	abilities map[string]string,
) NPC {
	// Optionally perform deep copies of the maps here.
	return NPC{
		id:          id,
		name:        name,
		faction:     faction,
		species:     species,
		npcType:     npcType,
		npcSubtype:  npcSubtype,
		trait:       trait,
		drive:       drive,
		description: description,
		stats:       copyIntMap(stats),
		items:       copyStringMap(items),
		abilities:   copyStringMap(abilities),
	}
}

func copyIntMap(m map[string]int) map[string]int {
	if m == nil {
		return nil
	}
	cp := make(map[string]int)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func copyStringMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	cp := make(map[string]string)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
