// Package sound provides audio feedback functionality for application events.
// It plays embedded WAV files when transcription starts and finishes.
package sound

import (
	"bytes"
	"context"
	"strconv"
	"sync"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
	"github.com/varavelio/tribar/assets/sounds"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

// Instance handles audio feedback playback.
type Instance struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager

	speakerOnce sync.Once
	speakerErr  error
	sampleRate  beep.SampleRate
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

	go s.playSound(ctx, true)
}

// TranscriptionFinished plays a sound when transcription completes.
func (s *Instance) TranscriptionFinished(ctx context.Context) {
	settings := s.settingsManager.Get()
	if !settings.SoundFeedbackEnable {
		return
	}

	go s.playSound(ctx, false)
}

// playSound plays the appropriate sound based on the current settings.
func (s *Instance) playSound(ctx context.Context, isInput bool) {
	settings := s.settingsManager.Get()

	idx := s.parseSoundID(settings.SoundFeedbackID)
	sound := sounds.Sounds[idx]

	var data []byte
	if isInput {
		data = sound.Input
	} else {
		data = sound.Output
	}

	streamer, format, err := wav.Decode(bytes.NewReader(data))
	if err != nil {
		s.logger.Error(ctx, "failed to decode WAV", "error", err)
		return
	}
	defer streamer.Close()

	s.speakerOnce.Do(func() {
		s.sampleRate = format.SampleRate
		s.speakerErr = speaker.Init(format.SampleRate, format.SampleRate.N(format.SampleRate.D(512)))
	})

	if s.speakerErr != nil {
		s.logger.Error(ctx, "failed to initialize speaker", "error", s.speakerErr)
		return
	}

	// Resample if the audio sample rate differs from speaker sample rate
	var toPlay beep.Streamer = streamer
	if format.SampleRate != s.sampleRate {
		toPlay = beep.Resample(4, format.SampleRate, s.sampleRate, streamer)
	}

	done := make(chan struct{})
	speaker.Play(beep.Seq(toPlay, beep.Callback(func() {
		close(done)
	})))

	select {
	case <-done:
	case <-ctx.Done():
	}
}

// parseSoundID converts the sound ID string to a valid array index (0-8).
func (s *Instance) parseSoundID(id string) int {
	n, err := strconv.Atoi(id)
	if err != nil || n < 1 || n > 9 {
		return 0 // Default to first sound
	}
	return n - 1
}

// Shutdown is a no-op for this implementation.
func (s *Instance) Shutdown() {
	// Nothing to clean up
}
