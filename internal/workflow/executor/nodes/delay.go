package nodes

import (
	"context"
	"fmt"
	"time"
)

// DelayNode waits for a specified duration.
type DelayNode struct{}

// NewDelayNode creates a new delay node.
func NewDelayNode() *DelayNode {
	return &DelayNode{}
}

// Type returns the node type identifier.
func (n *DelayNode) Type() string {
	return "delay"
}

// Execute waits for the configured duration and passes through input.
func (n *DelayNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	// Get delay duration in milliseconds
	delayMs := 0
	if v, ok := input.Config["delayMs"].(float64); ok {
		delayMs = int(v)
	} else if v, ok := input.Config["delayMs"].(int); ok {
		delayMs = v
	}

	if delayMs <= 0 {
		delayMs = 1000 // Default to 1 second
	}

	// Cap at 30 seconds to prevent abuse
	if delayMs > 30000 {
		delayMs = 30000
	}

	// Wait with context awareness
	select {
	case <-ctx.Done():
		return EmptyOutput(), ctx.Err()
	case <-time.After(time.Duration(delayMs) * time.Millisecond):
	}

	// Pass through the input
	passthrough := input.Config["input"]
	if passthrough == nil {
		passthrough = ""
	}

	return NewNodeOutput(map[string]interface{}{
		"output": fmt.Sprintf("%v", passthrough),
	}), nil
}
