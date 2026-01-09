package nodes

import (
	"context"
	"fmt"
)

// TranscribeNode converts audio data to text.
type TranscribeNode struct{}

// NewTranscribeNode creates a new transcribe node.
func NewTranscribeNode() *TranscribeNode {
	return &TranscribeNode{}
}

// Type returns the node type identifier.
func (n *TranscribeNode) Type() string {
	return "transcribe"
}

// Execute transcribes the audio data from the trigger.
func (n *TranscribeNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	// Get WAV data from trigger data or config
	var wavData []byte
	if data, ok := input.Config["audioData"].([]byte); ok {
		wavData = data
	} else if data, ok := input.Data["wavData"].([]byte); ok {
		wavData = data
	}

	if len(wavData) == 0 {
		return EmptyOutput(), fmt.Errorf("no audio data provided")
	}

	// Perform transcription
	text, err := services.GetTranscriber().TranscribeWAV(wavData)
	if err != nil {
		return EmptyOutput(), fmt.Errorf("transcription failed: %w", err)
	}

	return NewNodeOutput(map[string]interface{}{
		"rawText": text,
	}), nil
}
