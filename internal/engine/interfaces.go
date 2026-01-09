package engine

import (
	"context"

	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/history"
	"github.com/varavelio/tribar/internal/state"
	"github.com/varavelio/tribar/internal/transcribe"
)

// SettingsManager defines the interface for retrieving application settings.
type SettingsManager interface {
	Get() config.Settings
}

// HistoryManager defines the interface for managing transcription history.
type HistoryManager interface {
	Write(ctx context.Context, req history.WriteRequest) (history.Entry, error)
}

// StateManager defines the interface for managing application state.
type StateManager interface {
	SetStatus(status state.Status)
	GetStatus() (current state.Status, previous state.Status)
	SetDownloadProgress(fileName string, downloaded, total int64, percent float64)
	ClearDownloadProgress()
}

// Recorder defines the interface for audio recording.
type Recorder interface {
	Start() error
	Stop()
	GetData() []byte
	BuildWAV() []byte
}

// Transcriber defines the interface for audio transcription.
type Transcriber interface {
	CheckModels() (bool, []transcribe.ModelFile)
	DownloadModels(progress transcribe.DownloadProgressCallback) error
	LoadModels() error
	UnloadModels()
	TranscribeWAV(wavData []byte) (string, error)
}

// PostProcess defines the interface for post-processing transcription text.
type PostProcess interface {
	IsEnabled() bool
	Process(ctx context.Context, text string) (string, error)
}

// ClipboardWriter defines the interface for writing to clipboard.
type ClipboardWriter interface {
	Write(ctx context.Context, mode config.OutputMode, pasteShortcut string, text string) error
}

// Notifier defines the interface for sending user notifications.
type Notifier interface {
	Error(ctx context.Context, title string, message string)
	TranscriptionStarted(ctx context.Context)
	TranscriptionFinished(ctx context.Context, text string)
}

// SoundManager defines the interface for playing sound effects.
type SoundManager interface {
	RecordingStarted(ctx context.Context)
	RecordingStopped(ctx context.Context)
	TranscriptionError(ctx context.Context)
	TranscriptionSuccess(ctx context.Context)
}
