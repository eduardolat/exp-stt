package state

import (
	"sync"
	"time"
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

// HistoryEntry represents a single transcription record.
type HistoryEntry struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	AudioPath string    `json:"audio_path"`
	Timestamp time.Time `json:"timestamp"`
}

// Instance represents the application state, this state is used in all other
// packages to react to the current state of the application.
type Instance struct {
	statusMu       sync.RWMutex
	statusPrevious Status
	statusCurrent  Status

	historyMu    sync.RWMutex
	history      []HistoryEntry
	historyLimit int
	nextID       int
}

// New creates a new Instance with the initial status set to StatusUnloaded.
func New(historyLimit int) *Instance {
	return &Instance{
		statusMu:       sync.RWMutex{},
		statusPrevious: StatusUnknown,
		statusCurrent:  StatusUnloaded,
		historyMu:      sync.RWMutex{},
		history:        make([]HistoryEntry, 0),
		historyLimit:   historyLimit,
		nextID:         1,
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

// AddHistoryEntry adds a new transcription to the history.
func (i *Instance) AddHistoryEntry(text, audioPath string) {
	i.historyMu.Lock()
	defer i.historyMu.Unlock()

	entry := HistoryEntry{
		ID:        i.nextID,
		Text:      text,
		AudioPath: audioPath,
		Timestamp: time.Now(),
	}
	i.nextID++

	i.history = append([]HistoryEntry{entry}, i.history...)

	if len(i.history) > i.historyLimit {
		i.history = i.history[:i.historyLimit]
	}
}

// GetHistory returns a copy of the transcription history.
func (i *Instance) GetHistory() []HistoryEntry {
	i.historyMu.RLock()
	defer i.historyMu.RUnlock()

	result := make([]HistoryEntry, len(i.history))
	copy(result, i.history)
	return result
}

// GetHistoryEntry retrieves a specific history entry by ID.
func (i *Instance) GetHistoryEntry(id int) (HistoryEntry, bool) {
	i.historyMu.RLock()
	defer i.historyMu.RUnlock()

	for _, entry := range i.history {
		if entry.ID == id {
			return entry, true
		}
	}
	return HistoryEntry{}, false
}

// ClearHistory removes all entries from the history.
func (i *Instance) ClearHistory() {
	i.historyMu.Lock()
	defer i.historyMu.Unlock()
	i.history = make([]HistoryEntry, 0)
}

// SetHistoryLimit updates the maximum number of history entries.
func (i *Instance) SetHistoryLimit(limit int) {
	i.historyMu.Lock()
	defer i.historyMu.Unlock()

	i.historyLimit = limit
	if len(i.history) > limit {
		i.history = i.history[:limit]
	}
}
