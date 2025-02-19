package shared

// ViewPlugin defines the interface for UI plugins
type NPCViewer interface {
	NPCObserver
	Render()
}
