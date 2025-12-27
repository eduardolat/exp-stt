// Package clipboard provides output functionality for transcription results.
// It supports three modes: copy only, copy and paste, and ghost paste.
package clipboard

import (
	"context"
	"fmt"
	"time"

	atclip "github.com/atotto/clipboard"
	"github.com/varavelio/tribar/internal/logger"
)

// OutputMode defines how transcription results are delivered.
type OutputMode int

const (
	// OutputModeCopyOnly copies text to clipboard only (no paste).
	OutputModeCopyOnly OutputMode = iota
	// OutputModeCopyPaste copies and pastes, keeping the text in clipboard.
	OutputModeCopyPaste
	// OutputModeGhostPaste copies and pastes, then restores original clipboard content.
	OutputModeGhostPaste
)

// Settings configures the write behavior.
type Settings struct {
	Mode OutputMode
}

// DefaultSettings returns the default write settings.
func DefaultSettings() Settings {
	return Settings{
		Mode: OutputModeCopyPaste,
	}
}

// Instance handles output of transcription results.
type Instance struct {
	logger   logger.Logger
	settings Settings
}

// New creates a new write instance.
func New(logger logger.Logger, settings Settings) *Instance {
	return &Instance{
		logger:   logger,
		settings: settings,
	}
}

// UpdateSettings updates the write settings.
func (w *Instance) UpdateSettings(settings Settings) {
	w.settings = settings
}

// GetSettings returns the current write settings.
func (w *Instance) GetSettings() Settings {
	return w.settings
}

// Write outputs the transcription result based on the configured mode.
func (w *Instance) Write(ctx context.Context, text string) error {
	if text == "" {
		return nil
	}

	switch w.settings.Mode {
	case OutputModeCopyOnly:
		return w.copyToClipboard(ctx, text)
	case OutputModeCopyPaste:
		return w.pasteWorkflow(ctx, text, false)
	case OutputModeGhostPaste:
		return w.pasteWorkflow(ctx, text, true)
	default:
		return w.copyToClipboard(ctx, text)
	}
}

// copyToClipboard copies text to the system clipboard.
func (w *Instance) copyToClipboard(ctx context.Context, text string) error {
	if err := atclip.WriteAll(text); err != nil {
		w.logger.Error(ctx, "failed to copy to clipboard", "err", err)
		return fmt.Errorf("clipboard error: %w", err)
	}
	return nil
}

// pasteWorkflow handles the copy-paste workflow with optional clipboard restoration.
func (w *Instance) pasteWorkflow(ctx context.Context, text string, restore bool) error {
	var originalContent string

	if restore {
		originalContent, _ = atclip.ReadAll()
	}

	if err := w.copyToClipboard(ctx, text); err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	if err := triggerPastePlatform(); err != nil {
		w.logger.Warn(ctx, "paste trigger failed, text remains in clipboard", "err", err)
		return err
	}

	if restore {
		go func() {
			// Wait for the OS to process the paste before restoring
			time.Sleep(250 * time.Millisecond)
			_ = atclip.WriteAll(originalContent)
		}()
	}

	return nil
}
