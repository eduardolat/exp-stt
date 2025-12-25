package app

type Status int

const (
	StatusUnknown Status = iota
	StatusUnloaded
	StatusLoading
	StatusLoaded
	StatusListening
	StatusTranscribing
	StatusPostProcessing
)

// Instance represents an application status instance, this status is used in all other
// packages to react to the current state of the application.
type Instance struct {
	StatusCurrent  Status
	StatusPrevious Status
}

// New creates a new Instance with the initial status set to StatusUnloaded.
func New() *Instance {
	return &Instance{
		StatusCurrent: StatusUnloaded,
	}
}

// SetStatus changes the current status of the application instance. It also updates
// the previous status to the current one before the change.
func (i *Instance) SetStatus(newStatus Status) {
	i.StatusPrevious = i.StatusCurrent
	i.StatusCurrent = newStatus
}
