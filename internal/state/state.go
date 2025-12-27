package state

import "sync"

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

// Instance represents the application state, this state is used in all other
// packages to react to the current state of the application.
type Instance struct {
	statusMu       sync.RWMutex
	statusPrevious Status
	statusCurrent  Status
}

// New creates a new Instance with the initial status set to StatusUnloaded.
func New() *Instance {
	return &Instance{
		statusMu:       sync.RWMutex{},
		statusPrevious: StatusUnknown,
		statusCurrent:  StatusUnloaded,
	}
}

// SetStatus changes the current status of the application instance. It also updates
// the previous status to the current one before the change.
func (i *Instance) SetStatus(newStatus Status) {
	i.statusMu.Lock()
	defer i.statusMu.Unlock()
	i.statusPrevious = i.statusCurrent
	i.statusCurrent = newStatus
}

// GetStatus retrieves the current and previous statuses of the application instance.
func (i *Instance) GetStatus() (current Status, previous Status) {
	i.statusMu.RLock()
	defer i.statusMu.RUnlock()
	return i.statusCurrent, i.statusPrevious
}
