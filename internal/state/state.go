// Package state manages the global application state in a thread-safe way.
package state

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/gen2brain/malgo"
	"github.com/varavelio/tribar/internal/eventbus"
	"github.com/varavelio/tribar/internal/history"
)

// AudioContext abstracts the audio context for testing.
type AudioContext interface {
	Devices(kind malgo.DeviceType) ([]malgo.DeviceInfo, error)
}

// HistoryManager abstracts the history manager for testing.
type HistoryManager interface {
	GetAll(ctx context.Context) []history.Entry
	GetByID(id string) (history.Entry, bool)
	GetAudioPath(id string) string
	Delete(ctx context.Context, id string) error
	Clear(ctx context.Context) error
	Count() int
}

// RuntimeInfo holds system information detected at startup.
// This is populated by InitRuntime and should not be modified afterwards.
var RuntimeInfo = SystemInfo{}

var runtimeOnce sync.Once

// InitRuntime detects and initializes the runtime environment information.
// This should be called once at application startup. It is safe to call multiple times.
func InitRuntime() {
	runtimeOnce.Do(func() {
		RuntimeInfo = SystemInfo{
			OS:            runtime.GOOS,
			Arch:          runtime.GOARCH,
			DisplayServer: getDisplayServer(),
		}
	})
}

// getDisplayServer detects the current display server environment.
func getDisplayServer() string {
	switch runtime.GOOS {
	case "windows":
		return "win32"
	case "darwin":
		return "cocoa"
	case "linux":
		WAYLAND_DISPLAY := strings.ToLower(os.Getenv("WAYLAND_DISPLAY"))
		XDG_SESSION_TYPE := strings.ToLower(os.Getenv("XDG_SESSION_TYPE"))
		DISPLAY := strings.ToLower(os.Getenv("DISPLAY"))

		if WAYLAND_DISPLAY != "" {
			return "wayland"
		}

		if strings.Contains(XDG_SESSION_TYPE, "wayland") {
			return "wayland"
		}
		if strings.Contains(XDG_SESSION_TYPE, "x11") {
			return "x11"
		}
		if DISPLAY != "" {
			return "x11"
		}

		return "unknown"
	default:
		return "unknown"
	}
}

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

// SystemInfo contains runtime environment information.
type SystemInfo struct {
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	DisplayServer string `json:"display_server"` // "x11", "wayland", "cocoa", "win32"
}

// Instance represents the application state, this state is used in all other
// packages to react to the current state of the application.
type Instance struct {
	eventBus       *eventbus.EventBus
	historyManager HistoryManager

	statusMu       sync.RWMutex
	statusPrevious Status
	statusCurrent  Status

	downloadMu       sync.RWMutex
	downloadProgress DownloadProgress

	audioCtx AudioContext
}

// New creates a new Instance with the initial status set to StatusUnloaded.
func New(eventBus *eventbus.EventBus, historyManager *history.Manager) *Instance {
	audioCtx, _ := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	return NewWithDependencies(eventBus, historyManager, audioCtx)
}

// NewWithDependencies creates a new Instance with injected dependencies.
func NewWithDependencies(eventBus *eventbus.EventBus, historyManager HistoryManager, audioCtx AudioContext) *Instance {
	return &Instance{
		eventBus:       eventBus,
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
	i.eventBus.PublishStateChanged()
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

	// Publish state changes only for every exact 5%
	percentsToPublish := []float64{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100}
	if slices.Contains(percentsToPublish, percent) {
		i.eventBus.PublishStateChanged()
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
	i.eventBus.PublishStateChanged()
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
	err := i.historyManager.Delete(ctx, id)
	if err == nil {
		i.eventBus.PublishStateChanged()
		return nil
	}
	return fmt.Errorf("failed to delete history entry: %w", err)
}

// ClearHistory removes all entries from the history.
func (i *Instance) ClearHistory(ctx context.Context) error {
	err := i.historyManager.Clear(ctx)
	if err == nil {
		i.eventBus.PublishStateChanged()
		return nil
	}
	return fmt.Errorf("failed to clear history: %w", err)
}

// HistoryCount returns the number of entries in the history.
func (i *Instance) HistoryCount() int {
	return i.historyManager.Count()
}
