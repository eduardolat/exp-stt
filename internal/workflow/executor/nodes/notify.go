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

	// Get message from config
	message, _ := input.Config["message"].(string)
	if message == "" {
		// Try to get from text field (common pattern)
		message, _ = input.Config["text"].(string)
	}

	if message != "" {
		notifier.TranscriptionFinished(ctx, message)
	}

	return EmptyOutput(), nil
}
