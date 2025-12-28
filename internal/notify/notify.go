// Package notify provides desktop notification functionality for the application.
// It uses the beeep library to display native desktop notifications across platforms.
package notify

import (
	"context"

	"github.com/gen2brain/beeep"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

// Settings configures notification behavior.
type Settings struct {
	NotifyOnError  bool // Always notify on errors (default: true)
	NotifyOnStart  bool // Notify when transcription starts
	NotifyOnFinish bool // Notify when transcription completes
}

// DefaultSettings returns the default notification settings.
func DefaultSettings() Settings {
	return Settings{
		NotifyOnError:  true,
		NotifyOnStart:  false,
		NotifyOnFinish: false,
	}
}

// Instance handles desktop notifications.
type Instance struct {
	logger   logger.Logger
	settings Settings
}

// New creates a new notification instance.
func New(logger logger.Logger, settings Settings) *Instance {
	return &Instance{
		logger:   logger,
		settings: settings,
	}
}

// UpdateSettings updates the notification settings.
func (n *Instance) UpdateSettings(settings Settings) {
	n.settings = settings
}

// GetSettings returns the current notification settings.
func (n *Instance) GetSettings() Settings {
	return n.settings
}

// Error displays an error notification if error notifications are enabled.
func (n *Instance) Error(ctx context.Context, title, message string) {
	if !n.settings.NotifyOnError {
		return
	}

	n.send(ctx, title, message)
}

// TranscriptionStarted displays a notification when transcription starts.
func (n *Instance) TranscriptionStarted(ctx context.Context) {
	if !n.settings.NotifyOnStart {
		return
	}

	n.send(ctx, config.AppName, "Recording started...")
}

// TranscriptionFinished displays a notification when transcription completes.
func (n *Instance) TranscriptionFinished(ctx context.Context, text string) {
	if !n.settings.NotifyOnFinish {
		return
	}

	message := text
	if len(message) > 100 {
		message = message[:97] + "..."
	}

	n.send(ctx, "Transcription Complete", message)
}

// send dispatches a notification to the desktop.
func (n *Instance) send(ctx context.Context, title, message string) {
	if err := beeep.Notify(title, message, ""); err != nil {
		n.logger.Error(ctx, "failed to send desktop notification",
			"title", title,
			"message", message,
			"err", err,
		)
	}
}
