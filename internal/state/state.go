// Package state manages the global application state in a thread-safe way.
package state

import (
	"context"
	"sync"

	"github.com/gen2brain/malgo"
	"github.com/varavelio/tribar/internal/history"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusUnloaded
	StatusDownloading
	StatusLoading
	StatusLoaded
	StatusListening
	StatusTranscribing
	StatusPostProcessing
)

// DownloadProgress represents the current download progress.
type DownloadProgress struct {
	FileName   string  `json:"file_name"`
	Downloaded int64   `json:"downloaded"`
	Total      int64   `json:"total"`
	Percent    float64 `json:"percent"`
}

// AudioDevice represents an available audio device.
type AudioDevice struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// AvailableDevices holds lists of input/output devices.
type AvailableDevices struct {
	InputDevices  []AudioDevice `json:"input_devices"`
	OutputDevices []AudioDevice `json:"output_devices"`
}

// Instance represents the application state, this state is used in all other
// packages to react to the current state of the application.
type Instance struct {
	historyManager *history.Manager

	statusMu       sync.RWMutex
	statusPrevious Status
	statusCurrent  Status

	downloadMu       sync.RWMutex
	downloadProgress DownloadProgress

	audioCtx *malgo.AllocatedContext
}

// New creates a new Instance with the initial status set to StatusUnloaded.
func New(historyManager *history.Manager) *Instance {
	audioCtx, _ := malgo.InitContext(nil, malgo.ContextConfig{}, nil)

	return &Instance{
		historyManager: historyManager,
		statusMu:       sync.RWMutex{},
		statusPrevious: StatusUnknown,
		statusCurrent:  StatusUnloaded,
		audioCtx:       audioCtx,
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

// SetDownloadProgress updates the current download progress.
func (i *Instance) SetDownloadProgress(fileName string, downloaded, total int64, percent float64) {
	i.downloadMu.Lock()
	defer i.downloadMu.Unlock()
	i.downloadProgress = DownloadProgress{
		FileName:   fileName,
		Downloaded: downloaded,
		Total:      total,
		Percent:    percent,
	}
}

// GetDownloadProgress returns the current download progress.
func (i *Instance) GetDownloadProgress() DownloadProgress {
	i.downloadMu.RLock()
	defer i.downloadMu.RUnlock()
	return i.downloadProgress
}

// ClearDownloadProgress clears the download progress.
func (i *Instance) ClearDownloadProgress() {
	i.downloadMu.Lock()
	defer i.downloadMu.Unlock()
	i.downloadProgress = DownloadProgress{}
}

// GetAvailableDevices returns lists of available input and output audio devices.
func (i *Instance) GetAvailableDevices() AvailableDevices {
	result := AvailableDevices{
		InputDevices:  []AudioDevice{},
		OutputDevices: []AudioDevice{},
	}

	if i.audioCtx == nil {
		return result
	}

	// Get capture (input) devices
	captureDevices, err := i.audioCtx.Devices(malgo.Capture)
	if err == nil {
		for idx, dev := range captureDevices {
			result.InputDevices = append(result.InputDevices, AudioDevice{
				ID:        dev.ID.String(),
				Name:      dev.Name(),
				IsDefault: idx == 0,
			})
		}
	}

	// Get playback (output) devices
	playbackDevices, err := i.audioCtx.Devices(malgo.Playback)
	if err == nil {
		for idx, dev := range playbackDevices {
			result.OutputDevices = append(result.OutputDevices, AudioDevice{
				ID:        dev.ID.String(),
				Name:      dev.Name(),
				IsDefault: idx == 0,
			})
		}
	}

	return result
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
