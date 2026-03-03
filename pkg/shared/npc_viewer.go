package shared

// ViewPlugin defines the interface for UI plugins
type NPCListViewer interface {

	// NPCObserver : The interface for the NPC observer.
	// It defines the methods that the observer must implement to receive updates from the NPC viewer.
	NPCObserver

	// Run : The method to run the NPC viewer.
	// It takes no parameters and returns nothing.
	Run()
}
