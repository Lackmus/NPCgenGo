package shared

// NPCViewer : The interface for the NPC viewer.
type NPCEditViewer interface {
	NPCEditObserver
	Run()
}
