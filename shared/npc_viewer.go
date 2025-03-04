package shared

// ViewPlugin defines the interface for UI plugins
type NPCListViewer interface {
	NPCObserver
	Run()
}
