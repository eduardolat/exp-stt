package state

import (
	"context"
	"sync"

	"github.com/varavelio/tribar/internal/history"
)

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
	historyManager *history.Manager

	statusMu       sync.RWMutex
	statusPrevious Status
	statusCurrent  Status
}

// New creates a new Instance with the initial status set to StatusUnloaded.
func New(historyManager *history.Manager) *Instance {
	return &Instance{
		historyManager: historyManager,
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

// GetHistory returns a copy of the transcription history.
// This refreshes the cache from disk to ensure data is up-to-date.
func (i *Instance) GetHistory(ctx context.Context) []history.Entry {
	return i.historyManager.GetAll(ctx)
}

// GetHistoryEntry retrieves a specific history entry by ID.
func (i *Instance) GetHistoryEntry(id string) (history.Entry, bool) {
	return i.historyManager.GetByID(id)
}

// GetHistoryAudioPath returns the full path to the audio file for an entry.
func (i *Instance) GetHistoryAudioPath(id string) string {
	return i.historyManager.GetAudioPath(id)
}

// DeleteHistoryEntry removes a history entry by ID.
func (i *Instance) DeleteHistoryEntry(ctx context.Context, id string) error {
	return i.historyManager.Delete(ctx, id)
}

// ClearHistory removes all entries from the history.
func (i *Instance) ClearHistory(ctx context.Context) error {
	return i.historyManager.Clear(ctx)
}

// HistoryCount returns the number of entries in the history.
func (i *Instance) HistoryCount() int {
	return i.historyManager.Count()
}
