package shared

// NPCViewer : The interface for the NPC viewer.
type NPCEditViewer interface {

	// NPCEditObserver : The interface for the NPC edit observer.
	// It defines the methods that the observer must implement to receive updates from the NPC viewer.
	NPCEditObserver

	// Run : The method to run the NPC viewer.
	// It takes no parameters and returns nothing.
	Run()
}

