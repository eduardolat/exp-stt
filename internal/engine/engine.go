// Package engine is the central orchestrator that connects all components and manages
// the transcription workflow. It is the only package allowed to modify the application state.
package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/varavelio/tribar/internal/clipboard"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/notify"
	"github.com/varavelio/tribar/internal/postprocess"
	"github.com/varavelio/tribar/internal/record"
	"github.com/varavelio/tribar/internal/sound"
	"github.com/varavelio/tribar/internal/state"
	"github.com/varavelio/tribar/internal/transcribe"
)

// Dependencies contains all required dependencies for the engine.
type Dependencies struct {
	Logger          logger.Logger
	SettingsManager *config.SettingsManager
	State           *state.Instance
	Recorder        *record.Recorder
	Transcriber     *transcribe.Instance
	PostProcess     *postprocess.Instance
	Writer          *clipboard.Instance
	Notifier        *notify.Instance
	Sound           *sound.Instance
}

// Engine orchestrates the transcription workflow.
type Engine struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager
	state           *state.Instance
	recorder        *record.Recorder
	transcriber     *transcribe.Instance
	postprocess     *postprocess.Instance
	writer          *clipboard.Instance
	notifier        *notify.Instance
	sound           *sound.Instance

	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new Engine instance with all dependencies.
func New(deps Dependencies) *Engine {
	ctx, cancel := context.WithCancel(context.Background())

	return &Engine{
		logger:          deps.Logger,
		settingsManager: deps.SettingsManager,
		state:           deps.State,
		recorder:        deps.Recorder,
		transcriber:     deps.Transcriber,
		postprocess:     deps.PostProcess,
		writer:          deps.Writer,
		notifier:        deps.Notifier,
		sound:           deps.Sound,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// LoadModels loads the transcription models with progress reporting.
func (e *Engine) LoadModels(progressCallback transcribe.DownloadProgressCallback) error {
	e.state.SetStatus(state.StatusLoading)

	allExist, _ := e.transcriber.CheckModels()
	if !allExist {
		e.logger.Info(e.ctx, "downloading missing models...")
		if err := e.transcriber.DownloadModels(progressCallback); err != nil {
			e.state.SetStatus(state.StatusUnloaded)
			e.notifier.Error(e.ctx, "Model Download Failed", err.Error())
			return fmt.Errorf("failed to download models: %w", err)
		}
	}

	if err := e.transcriber.LoadModels(); err != nil {
		e.state.SetStatus(state.StatusUnloaded)
		e.notifier.Error(e.ctx, "Model Load Failed", err.Error())
		return fmt.Errorf("failed to load models: %w", err)
	}

	e.state.SetStatus(state.StatusLoaded)
	e.logger.Info(e.ctx, "models loaded successfully")
	return nil
}

// ToggleRecording starts or stops the recording based on current state.
func (e *Engine) ToggleRecording() {
	status, _ := e.state.GetStatus()

	switch status {
	case state.StatusListening:
		e.stopRecording()
	case state.StatusLoaded:
		e.startRecording()
	case state.StatusUnloaded:
		e.logger.Warn(e.ctx, "cannot start recording, models not loaded")
	}
}

// StartRecording begins audio capture.
func (e *Engine) startRecording() {
	if err := e.recorder.Start(); err != nil {
		e.logger.Error(e.ctx, "failed to start recording", "err", err)
		e.notifier.Error(e.ctx, "Recording Failed", err.Error())
		return
	}

	e.state.SetStatus(state.StatusListening)
	e.sound.TranscriptionStarted(e.ctx)
	e.notifier.TranscriptionStarted(e.ctx)
	e.logger.Info(e.ctx, "recording started")
}

// stopRecording stops audio capture and processes the recording.
func (e *Engine) stopRecording() {
	e.recorder.Stop()
	e.logger.Info(e.ctx, "recording stopped, processing...")

	go e.processRecording()
}

// processRecording handles the transcription pipeline in a goroutine.
func (e *Engine) processRecording() {
	settings := e.settingsManager.Get()
	e.state.SetStatus(state.StatusTranscribing)

	audioPath := e.generateAudioPath()
	if err := e.recorder.SaveWAV(audioPath); err != nil {
		e.handleError("failed to save audio", err)
		return
	}

	wavData, err := os.ReadFile(audioPath)
	if err != nil {
		e.handleError("failed to read audio file", err)
		return
	}

	text, err := e.transcriber.TranscribeWAV(wavData)
	if err != nil {
		e.handleError("transcription failed", err)
		return
	}

	e.logger.Debug(e.ctx, "transcription complete", "text", text)

	if e.postprocess.IsEnabled() {
		e.state.SetStatus(state.StatusPostProcessing)
		processed, err := e.postprocess.Process(e.ctx, text)
		if err != nil {
			e.logger.Warn(e.ctx, "post-processing failed, using raw transcription", "err", err)
		} else {
			text = processed
		}
	}

	if err := e.writer.Write(e.ctx, settings.OutputMode, text); err != nil {
		e.logger.Error(e.ctx, "failed to write output", "err", err)
	}

	e.state.AddHistoryEntry(text, audioPath)
	e.sound.TranscriptionFinished(e.ctx)
	e.notifier.TranscriptionFinished(e.ctx, text)
	e.state.SetStatus(state.StatusLoaded)

	e.logger.Info(e.ctx, "transcription complete", "length", len(text))
}

// handleError logs the error, notifies the user, and resets state.
func (e *Engine) handleError(message string, err error) {
	e.logger.Error(e.ctx, message, "err", err)
	e.notifier.Error(e.ctx, config.AppName, fmt.Sprintf("%s: %v", message, err))
	e.state.SetStatus(state.StatusLoaded)
}

// generateAudioPath creates a unique path for the audio file.
func (e *Engine) generateAudioPath() string {
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("recording-%s.wav", timestamp)
	return filepath.Join(config.DirectoryRecordings, filename)
}

// GetState returns the current application state (read-only access for UI).
func (e *Engine) GetState() *state.Instance {
	return e.state
}

// Shutdown gracefully stops the engine and releases resources.
func (e *Engine) Shutdown() {
	e.cancel()

	status, _ := e.state.GetStatus()
	if status == state.StatusListening {
		e.recorder.Stop()
	}

	e.logger.Info(e.ctx, "engine shutdown complete")
}
