// Package clipboard provides output functionality for transcription results.
// It supports three modes: copy only, copy and paste, and ghost paste.
package clipboard

import (
	"context"
	"fmt"
	"time"

	atclip "github.com/atotto/clipboard"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

// Instance handles output of transcription results.
type Instance struct {
	logger logger.Logger
}

// New creates a new clipboard instance.
func New(logger logger.Logger) *Instance {
	return &Instance{
		logger: logger,
	}
}

// Write outputs the transcription result based on the configured mode.
func (w *Instance) Write(ctx context.Context, mode config.OutputMode, text string) error {
	if text == "" {
		return nil
	}

	switch mode {
	case config.OutputModeCopyOnly:
		return w.copyToClipboard(ctx, text)
	case config.OutputModeCopyPaste:
		return w.pasteWorkflow(ctx, text, false)
	case config.OutputModeGhostPaste:
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
