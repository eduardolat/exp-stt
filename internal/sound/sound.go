// Package sound provides audio feedback functionality for application events.
// It uses simple system commands to play audio cues for transcription start/end events.
package sound

import (
	"context"
	"os/exec"
	"runtime"

	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

// Instance handles audio feedback.
type Instance struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager
}

// New creates a new sound instance.
func New(logger logger.Logger, settingsManager *config.SettingsManager) *Instance {
	return &Instance{
		logger:          logger,
		settingsManager: settingsManager,
	}
}

// TranscriptionStarted plays a sound when transcription starts.
func (s *Instance) TranscriptionStarted(ctx context.Context) {
	settings := s.settingsManager.Get()
	if !settings.SoundFeedbackEnable {
		return
	}

	go s.playBeep(ctx, 440, 100) // A4 note, 100ms
}

// TranscriptionFinished plays a sound when transcription completes.
func (s *Instance) TranscriptionFinished(ctx context.Context) {
	settings := s.settingsManager.Get()
	if !settings.SoundFeedbackEnable {
		return
	}

	go s.playBeep(ctx, 880, 150) // A5 note, 150ms
}

// playBeep plays a beep sound using system tools.
func (s *Instance) playBeep(ctx context.Context, frequency, durationMs int) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		// Try paplay with a generated tone, fallback to speaker-test
		cmd = exec.CommandContext(ctx, "paplay", "--volume=32768", "/usr/share/sounds/freedesktop/stereo/message.oga")
		if err := cmd.Run(); err != nil {
			// Fallback: try using beep command if available
			_ = exec.CommandContext(ctx, "beep", "-f", itoa(frequency), "-l", itoa(durationMs)).Run()
		}
	case "darwin":
		// macOS: use afplay with system sound
		cmd = exec.CommandContext(ctx, "afplay", "/System/Library/Sounds/Pop.aiff")
		_ = cmd.Run()
	case "windows":
		// Windows: use PowerShell to play a beep
		cmd = exec.CommandContext(ctx, "powershell", "-c", "[console]::beep("+itoa(frequency)+","+itoa(durationMs)+")")
		_ = cmd.Run()
	default:
		s.logger.Debug(ctx, "sound playback not supported on this platform")
	}
}

// itoa converts an int to string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var digits []byte
	negative := n < 0
	if negative {
		n = -n
	}

	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}

// Shutdown is a no-op for this implementation.
func (s *Instance) Shutdown() {
	// Nothing to clean up
}
