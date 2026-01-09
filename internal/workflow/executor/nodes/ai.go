package nodes

import (
	"context"
	"fmt"
)

// AINode processes text using LLM.
type AINode struct{}

// NewAINode creates a new AI processing node.
func NewAINode() *AINode {
	return &AINode{}
}

// Type returns the node type identifier.
func (n *AINode) Type() string {
	return "ai_process"
}

// Execute processes text with the configured LLM.
func (n *AINode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	// Get text to process
	text, _ := input.Config["text"].(string)
	if text == "" {
		return EmptyOutput(), fmt.Errorf("no text provided for AI processing")
	}

	postProcessor := services.GetPostProcessor()

	// Check if AI processing is enabled
	if !postProcessor.IsEnabled() {
		// Pass through if AI is disabled
		return NewNodeOutput(map[string]interface{}{
			"processedText": text,
			"tokensUsed":    0,
		}), nil
	}

	// Process with LLM
	processed, err := postProcessor.Process(ctx, text)
	if err != nil {
		return EmptyOutput(), fmt.Errorf("AI processing failed: %w", err)
	}

	return NewNodeOutput(map[string]interface{}{
		"processedText": processed,
		"tokensUsed":    0,
	}), nil
}
