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

// PasteShortcut represents a validated keyboard shortcut for pasting.
type PasteShortcut string

const (
	// PasteShortcutCtrlV is the standard Ctrl+V paste shortcut.
	PasteShortcutCtrlV PasteShortcut = "ctrl+v"
	// PasteShortcutCtrlShiftV is the Ctrl+Shift+V paste shortcut (often used for plain text paste).
	PasteShortcutCtrlShiftV PasteShortcut = "ctrl+shift+v"
	// PasteShortcutShiftInsert is the Shift+Insert paste shortcut.
	PasteShortcutShiftInsert PasteShortcut = "shift+insert"
)

// parsePasteShortcut parses a string into a PasteShortcut, returning the default if invalid.
func parsePasteShortcut(s string) PasteShortcut {
	switch s {
	case string(PasteShortcutCtrlV):
		return PasteShortcutCtrlV
	case string(PasteShortcutCtrlShiftV):
		return PasteShortcutCtrlShiftV
	case string(PasteShortcutShiftInsert):
		return PasteShortcutShiftInsert
	default:
		return PasteShortcutCtrlV // Fallback to default
	}
}

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
func (w *Instance) Write(ctx context.Context, mode config.OutputMode, pasteShortcutRaw string, text string) error {
	if text == "" {
		return nil
	}

	shortcut := parsePasteShortcut(pasteShortcutRaw)

	switch mode {
	case config.OutputModeCopyOnly:
		return w.copyToClipboard(ctx, text)
	case config.OutputModeCopyPaste:
		return w.pasteWorkflow(ctx, shortcut, text, false)
	case config.OutputModeGhostPaste:
		return w.pasteWorkflow(ctx, shortcut, text, true)
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
func (w *Instance) pasteWorkflow(ctx context.Context, shortcut PasteShortcut, text string, restore bool) error {
	var originalContent string

	if restore {
		originalContent, _ = atclip.ReadAll()
	}

	if err := w.copyToClipboard(ctx, text); err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	if err := triggerPastePlatform(shortcut); err != nil {
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
