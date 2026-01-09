package nodes

import (
	"context"
	"fmt"

	"github.com/varavelio/tribar/internal/config"
)

// ClipboardCopyNode copies text to the clipboard.
type ClipboardCopyNode struct{}

// NewClipboardCopyNode creates a new clipboard copy node.
func NewClipboardCopyNode() *ClipboardCopyNode {
	return &ClipboardCopyNode{}
}

// Type returns the node type identifier.
func (n *ClipboardCopyNode) Type() string {
	return "clipboard_copy"
}

// Execute copies text to the clipboard.
func (n *ClipboardCopyNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	text, _ := input.Config["text"].(string)
	if text == "" {
		return NewNodeOutput(map[string]interface{}{
			"text": "",
		}), nil
	}

	// Copy only, don't paste
	settings := services.GetSettingsManager().Get()
	if err := services.GetClipboard().Write(ctx, config.OutputModeCopyOnly, settings.PasteShortcut, text); err != nil {
		return EmptyOutput(), fmt.Errorf("clipboard copy failed: %w", err)
	}

	return NewNodeOutput(map[string]interface{}{
		"text": text,
	}), nil
}

// ClipboardPasteNode pastes text using keyboard shortcut.
type ClipboardPasteNode struct{}

// NewClipboardPasteNode creates a new clipboard paste node.
func NewClipboardPasteNode() *ClipboardPasteNode {
	return &ClipboardPasteNode{}
}

// Type returns the node type identifier.
func (n *ClipboardPasteNode) Type() string {
	return "clipboard_paste"
}

// Execute copies and pastes text.
func (n *ClipboardPasteNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	text, _ := input.Config["text"].(string)
	if text == "" {
		return NewNodeOutput(map[string]interface{}{
			"text": "",
		}), nil
	}

	// Apply trailing space if configured
	settings := services.GetSettingsManager().Get()
	outputText := text
	if settings.OutputTrailingSpace && outputText != "" {
		outputText = outputText + " "
	}

	// Use the configured output mode
	if err := services.GetClipboard().Write(ctx, settings.OutputMode, settings.PasteShortcut, outputText); err != nil {
		return EmptyOutput(), fmt.Errorf("clipboard paste failed: %w", err)
	}

	return NewNodeOutput(map[string]interface{}{
		"text": text,
	}), nil
}
