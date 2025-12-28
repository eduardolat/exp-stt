// Package notify provides desktop notification functionality for the application.
// It uses the beeep library to display native desktop notifications across platforms.
package notify

import (
	"context"
	"fmt"

	"github.com/gen2brain/beeep"
	"github.com/varavelio/tribar/assets/logo"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

// Instance handles desktop notifications.
type Instance struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager
}

// New creates a new notification instance.
func New(logger logger.Logger, settingsManager *config.SettingsManager) *Instance {
	beeep.AppName = fmt.Sprintf("%s v%s", config.AppName, config.AppVersion)

	return &Instance{
		logger:          logger,
		settingsManager: settingsManager,
	}
}

// Error displays an error notification if error notifications are enabled.
func (n *Instance) Error(ctx context.Context, title, message string) {
	settings := n.settingsManager.Get()
	if !settings.NotifyOnError {
		return
	}

	n.send(ctx, title, message)
}

// TranscriptionStarted displays a notification when transcription starts.
func (n *Instance) TranscriptionStarted(ctx context.Context) {
	settings := n.settingsManager.Get()
	if !settings.NotifyOnStart {
		return
	}

	n.send(ctx, "Recording started", "Speak now...")
}

// TranscriptionFinished displays a notification when transcription completes.
func (n *Instance) TranscriptionFinished(ctx context.Context, text string) {
	settings := n.settingsManager.Get()
	if !settings.NotifyOnFinish {
		return
	}

	message := text
	if len(message) > 100 {
		message = message[:97] + "..."
	}

	n.send(ctx, "Transcription completed", message)
}

// send dispatches a notification to the desktop.
func (n *Instance) send(ctx context.Context, title, message string) {
	if err := beeep.Notify(title, message, logo.LogoBlackWhite.PNG.Size128.Logo); err != nil {
		n.logger.Error(ctx, "failed to send desktop notification",
			"title", title,
			"message", message,
			"err", err,
		)
	}
}
