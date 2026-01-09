package nodes

import (
	"context"
)

// NotifyNode shows a desktop notification.
type NotifyNode struct{}

// NewNotifyNode creates a new notification node.
func NewNotifyNode() *NotifyNode {
	return &NotifyNode{}
}

// Type returns the node type identifier.
func (n *NotifyNode) Type() string {
	return "notify"
}

// Execute shows a desktop notification.
func (n *NotifyNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	notifier := services.GetNotifier()

	// Get message from config first, fallback to input data (from previous node)
	message, _ := input.Config["message"].(string)
	if message == "" {
		message, _ = input.Config["text"].(string)
	}
	// Fallback to input data from previous node
	if message == "" {
		message, _ = input.Data["rawText"].(string)
	}
	if message == "" {
		message, _ = input.Data["text"].(string)
	}
	if message == "" {
		message, _ = input.Data["processedText"].(string)
	}

	if message != "" {
		notifier.TranscriptionFinished(ctx, message)
	}

	return NewNodeOutput(map[string]interface{}{
		"text": message,
	}), nil
}
