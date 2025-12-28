// Package sound provides audio feedback functionality for application events.
// It plays embedded WAV files for recording start/stop, success, and error events.
package sound

import (
	"bytes"
	"context"
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

// RecordingStarted plays a sound when recording starts.
func (s *Instance) RecordingStarted(ctx context.Context) {
	settings := s.settingsManager.Get()
	if !settings.SoundFeedbackEnable {
		return
	}

	sound, ok := sounds.InOutSoundsMap[settings.SoundFeedbackRecordID]
	if !ok {
		sound = sounds.InOutSounds[0]
	}

	go s.playSound(ctx, sound.Input)
}

// RecordingStopped plays a sound when recording stops.
func (s *Instance) RecordingStopped(ctx context.Context) {
	settings := s.settingsManager.Get()
	if !settings.SoundFeedbackEnable {
		return
	}

	sound, ok := sounds.InOutSoundsMap[settings.SoundFeedbackRecordID]
	if !ok {
		sound = sounds.InOutSounds[0]
	}

	go s.playSound(ctx, sound.Output)
}

// TranscriptionSuccess plays a sound when transcription completes successfully.
func (s *Instance) TranscriptionSuccess(ctx context.Context) {
	settings := s.settingsManager.Get()
	if !settings.SoundFeedbackEnable {
		return
	}

	sound, ok := sounds.SuccessSoundsMap[settings.SoundFeedbackSuccessID]
	if !ok {
		sound = sounds.SuccessSounds[0]
	}

	go s.playSound(ctx, sound.Sound)
}

// TranscriptionError plays a sound when transcription fails.
func (s *Instance) TranscriptionError(ctx context.Context) {
	settings := s.settingsManager.Get()
	if !settings.SoundFeedbackEnable {
		return
	}

	sound, ok := sounds.ErrorSoundsMap[settings.SoundFeedbackErrorID]
	if !ok {
		sound = sounds.ErrorSounds[0]
	}

	go s.playSound(ctx, sound.Sound)
}

// playSound plays the provided WAV data through the speaker.
func (s *Instance) playSound(ctx context.Context, data []byte) {
	streamer, format, err := wav.Decode(bytes.NewReader(data))
	if err != nil {
		s.logger.Error(ctx, "failed to decode WAV", "error", err)
		return
	}
	defer func() { _ = streamer.Close() }()

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

// Shutdown is a no-op for this implementation.
func (s *Instance) Shutdown() {
	// Nothing to clean up
}
