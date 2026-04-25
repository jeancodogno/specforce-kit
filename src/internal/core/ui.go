package core

// UI defines the interface for interacting with the user and logging progress.
type UI interface {
	Log(message string)
	Warn(message string)
	Error(message string)
	Success(message string)
	
	// SubTask logs a smaller unit of work.
	SubTask(message string)
	
	// Spinner management
	StartSpinner(message string)
	StopSpinner()

	// Interactivity
	Confirm(question string) bool
}
