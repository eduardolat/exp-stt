package nodes

import (
	"context"
)

// SoundNode plays audio feedback.
type SoundNode struct{}

// NewSoundNode creates a new sound node.
func NewSoundNode() *SoundNode {
	return &SoundNode{}
}

// Type returns the node type identifier.
func (n *SoundNode) Type() string {
	return "sound"
}

// Execute plays the configured sound.
func (n *SoundNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	snd := services.GetSound()
	soundType, _ := input.Config["soundType"].(string)

	switch soundType {
	case "success":
		snd.TranscriptionSuccess(ctx)
	case "error":
		snd.TranscriptionError(ctx)
	case "start":
		snd.RecordingStarted(ctx)
	case "stop":
		snd.RecordingStopped(ctx)
	default:
		snd.TranscriptionSuccess(ctx)
	}

	return EmptyOutput(), nil
}
