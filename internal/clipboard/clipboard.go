// Package clipboard provides output functionality for transcription results.
// It supports three modes: copy only, copy and paste, and ghost paste.
package clipboard

import (
	"context"
	"fmt"
	"sync"
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
	mu     sync.Mutex
}

// New creates a new clipboard instance.
func New(logger logger.Logger) *Instance {
	return &Instance{
		logger: logger,
		mu:     sync.Mutex{},
	}
}

// Write outputs the transcription result based on the configured mode.
func (w *Instance) Write(ctx context.Context, mode config.OutputMode, pasteShortcutRaw string, text string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if text == "" {
		return nil
	}

	shortcut := parsePasteShortcut(pasteShortcutRaw)

	// workflowFunc is the function that will be called to write to the clipboard
	// it defaults and fallbacks to copy only workflow
	workflowFunc := w.copyOnlyWorkflow
	if mode == config.OutputModeCopyPaste {
		workflowFunc = w.copyPasteWorkflow
	}
	if mode == config.OutputModeGhostPaste {
		workflowFunc = w.ghostPasteWorkflow
	}

	if err := workflowFunc(ctx, shortcut, text); err != nil {
		w.logger.Error(ctx, "failed to write to clipboard", "mode", mode, "shortcut", shortcut, "text", text, "err", err)
		return fmt.Errorf("%s error: %w", mode, err)
	}

	return nil
}

// copyOnlyWorkflow copies text to the system clipboard.
func (w *Instance) copyOnlyWorkflow(ctx context.Context, _ PasteShortcut, text string) error {
	if err := atclip.WriteAll(text); err != nil {
		w.logger.Error(ctx, "failed to copy to clipboard", "err", err)
		return fmt.Errorf("clipboard error: %w", err)
	}
	return nil
}

// copyPasteWorkflow handles the copy-paste workflow.
func (w *Instance) copyPasteWorkflow(ctx context.Context, shortcut PasteShortcut, text string) error {
	if err := w.copyOnlyWorkflow(ctx, shortcut, text); err != nil {
		return err
	}

	// Wait for the OS to process the copy before pasting
	time.Sleep(50 * time.Millisecond)
	if err := triggerPastePlatform(shortcut); err != nil {
		w.logger.Warn(ctx, "paste trigger failed, text remains in clipboard", "err", err)
		return err
	}

	return nil
}

// ghostPasteWorkflow handles the ghost-paste workflow.
func (w *Instance) ghostPasteWorkflow(ctx context.Context, shortcut PasteShortcut, text string) error {
	originalContent, _ := atclip.ReadAll()

	if err := w.copyPasteWorkflow(ctx, shortcut, text); err != nil {
		return err
	}

	// Wait for the OS to process the paste before restoring
	time.Sleep(150 * time.Millisecond)
	_ = atclip.WriteAll(originalContent)

	return nil
}
