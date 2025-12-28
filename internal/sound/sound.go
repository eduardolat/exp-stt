// Package sound provides audio feedback functionality for application events.
// It plays embedded WAV files for recording start/stop, success, and error events.
package sound

import (
	"bytes"
	"context"
	"math"
	"sync"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
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
	settings := s.settingsManager.Get()

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

	// Apply volume setting (0-100 scale converted to dB-like behavior)
	volume := min(max(settings.SoundFeedbackVolume, 0), 100)

	if volume < 100 {
		// Convert 0-100 to a logarithmic volume scale
		// 100 = 0 dB (no change), 0 = -infinity (silence)
		// We use a range of about -40 dB to 0 dB for practical volume control
		var volumeDB float64
		if volume == 0 {
			volumeDB = -100 // Effectively silent
		} else {
			// Map 1-100 to -40dB to 0dB using logarithmic scale
			volumeDB = 20 * math.Log10(float64(volume)/100)
		}
		toPlay = &effects.Volume{
			Streamer: toPlay,
			Base:     2,
			Volume:   volumeDB / 10, // Convert to beep's volume scale (base 2)
			Silent:   volume == 0,
		}
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
