// Package engine is the central orchestrator that connects all components and manages
// the transcription workflow. It is the only package allowed to modify the application state.
package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/varavelio/tribar/internal/clipboard"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/history"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/notify"
	"github.com/varavelio/tribar/internal/postprocess"
	"github.com/varavelio/tribar/internal/record"
	"github.com/varavelio/tribar/internal/sound"
	"github.com/varavelio/tribar/internal/state"
	"github.com/varavelio/tribar/internal/transcribe"
)

const unloadTickerInterval = 10 * time.Second

// Dependencies contains all required dependencies for the engine.
type Dependencies struct {
	Logger          logger.Logger
	SettingsManager *config.SettingsManager
	HistoryManager  *history.Manager
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
	historyManager  *history.Manager
	state           *state.Instance
	recorder        *record.Recorder
	transcriber     *transcribe.Instance
	postprocess     *postprocess.Instance
	writer          *clipboard.Instance
	notifier        *notify.Instance
	sound           *sound.Instance

	ctx    context.Context
	cancel context.CancelFunc

	lastActivityMu   sync.RWMutex
	lastActivityTime time.Time

	unloadTickerStarted bool
	modelsLoaded        bool
	modelsLoadedMu      sync.RWMutex
}

// New creates a new Engine instance with all dependencies.
func New(deps Dependencies) *Engine {
	ctx, cancel := context.WithCancel(context.Background())

	return &Engine{
		logger:          deps.Logger,
		settingsManager: deps.SettingsManager,
		historyManager:  deps.HistoryManager,
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

// EnsureModelsDownloaded downloads models if they don't exist.
// This should be called at startup - it only downloads, does not load into memory.
func (e *Engine) EnsureModelsDownloaded() error {
	allExist, _ := e.transcriber.CheckModels()
	if allExist {
		e.logger.Info(e.ctx, "models already downloaded")
		return nil
	}

	e.state.SetStatus(state.StatusDownloading)
	e.logger.Info(e.ctx, "downloading missing models...")

	progressCallback := func(filename string, downloaded, total int64, percent float64) {
		e.state.SetDownloadProgress(filename, downloaded, total, percent)
	}

	if err := e.transcriber.DownloadModels(progressCallback); err != nil {
		e.state.SetStatus(state.StatusUnloaded)
		e.state.ClearDownloadProgress()
		e.notifier.Error(e.ctx, "Model Download Failed", err.Error())
		return fmt.Errorf("failed to download models: %w", err)
	}

	e.state.ClearDownloadProgress()
	e.state.SetStatus(state.StatusUnloaded)
	e.logger.Info(e.ctx, "models downloaded successfully")
	return nil
}

// ensureModelsLoaded loads the transcription models into memory if not already loaded.
// Returns true if models are ready, false if there was an error.
func (e *Engine) ensureModelsLoaded() bool {
	e.modelsLoadedMu.RLock()
	if e.modelsLoaded {
		e.modelsLoadedMu.RUnlock()
		return true
	}
	e.modelsLoadedMu.RUnlock()

	e.modelsLoadedMu.Lock()
	defer e.modelsLoadedMu.Unlock()

	// Double-check after acquiring write lock
	if e.modelsLoaded {
		return true
	}

	e.state.SetStatus(state.StatusLoading)

	if err := e.transcriber.LoadModels(); err != nil {
		e.notifier.Error(e.ctx, "Model Load Failed", err.Error())
		e.logger.Error(e.ctx, "failed to load models", "err", err)
		return false
	}

	e.modelsLoaded = true
	e.logger.Info(e.ctx, "models loaded into memory")

	// Start unload ticker if not already started
	if !e.unloadTickerStarted {
		e.unloadTickerStarted = true
		go e.unloadTicker()
	}

	return true
}

// UnloadModels unloads the transcription models to free resources.
func (e *Engine) UnloadModels() {
	e.modelsLoadedMu.Lock()
	defer e.modelsLoadedMu.Unlock()

	if !e.modelsLoaded {
		return
	}

	e.transcriber.UnloadModels()
	e.modelsLoaded = false
	e.state.SetStatus(state.StatusUnloaded)
	e.logger.Info(e.ctx, "models unloaded due to inactivity")
}

// recordActivity updates the last activity timestamp.
func (e *Engine) recordActivity() {
	e.lastActivityMu.Lock()
	defer e.lastActivityMu.Unlock()
	e.lastActivityTime = time.Now()
}

// unloadTicker runs a background ticker that checks for inactivity and unloads models.
func (e *Engine) unloadTicker() {
	ticker := time.NewTicker(unloadTickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			e.checkAndUnload()
		}
	}
}

// checkAndUnload checks if the model should be unloaded based on settings and last activity.
func (e *Engine) checkAndUnload() {
	settings := e.settingsManager.Get()
	if !settings.ModelUnloadEnable {
		return
	}

	status, _ := e.state.GetStatus()
	if status != state.StatusUnloaded {
		// Only unload when idle (not listening, transcribing, etc.)
		if status != state.StatusLoaded {
			return
		}
	}

	e.modelsLoadedMu.RLock()
	loaded := e.modelsLoaded
	e.modelsLoadedMu.RUnlock()

	if !loaded {
		return
	}

	e.lastActivityMu.RLock()
	lastActivity := e.lastActivityTime
	e.lastActivityMu.RUnlock()

	elapsed := time.Since(lastActivity)
	timeout := time.Duration(settings.ModelUnloadSeconds) * time.Second

	if elapsed >= timeout {
		e.UnloadModels()
	}
}

// ToggleRecording starts or stops the recording based on current state.
func (e *Engine) ToggleRecording() {
	status, _ := e.state.GetStatus()

	switch status {
	case state.StatusListening:
		e.stopRecording()
	case state.StatusLoaded, state.StatusUnloaded:
		e.startRecording()
	case state.StatusDownloading:
		e.logger.Warn(e.ctx, "cannot start recording, models are being downloaded")
		e.notifier.Error(e.ctx, config.AppName, "Please wait, models are being downloaded")
	case state.StatusLoading:
		e.logger.Warn(e.ctx, "cannot start recording, models are being loaded")
		e.notifier.Error(e.ctx, config.AppName, "Please wait, models are being loaded")
	}
}

// startRecording begins audio capture immediately.
func (e *Engine) startRecording() {
	e.recordActivity()

	if err := e.recorder.Start(); err != nil {
		e.logger.Error(e.ctx, "failed to start recording", "err", err)
		e.notifier.Error(e.ctx, "Recording Failed", err.Error())
		return
	}

	e.state.SetStatus(state.StatusListening)
	e.sound.RecordingStarted(e.ctx)
	e.notifier.TranscriptionStarted(e.ctx)
	e.logger.Info(e.ctx, "recording started")
}

// stopRecording stops audio capture and processes the recording.
func (e *Engine) stopRecording() {
	e.recorder.Stop()
	e.sound.RecordingStopped(e.ctx)
	e.logger.Info(e.ctx, "recording stopped, processing...")

	go e.processRecording()
}

// processRecording handles the transcription pipeline in a goroutine.
func (e *Engine) processRecording() {
	settings := e.settingsManager.Get()
	startedAt := time.Now()

	// Get the raw audio data from the recorder
	audioData := e.recorder.GetData()
	recordingDurationMs := calculateDurationMs(len(audioData))
	wavData := e.recorder.BuildWAV()

	// Ensure models are loaded before transcription
	if !e.ensureModelsLoaded() {
		e.sound.TranscriptionError(e.ctx)
		e.notifier.Error(e.ctx, config.AppName, "Failed to load models for transcription")
		e.state.SetStatus(state.StatusUnloaded)
		return
	}

	e.state.SetStatus(state.StatusTranscribing)

	text, err := e.transcriber.TranscribeWAV(wavData)
	if err != nil {
		e.handleError("transcription failed", err)
		return
	}

	e.logger.Debug(e.ctx, "transcription complete", "text", text)

	rawText := text
	postProcessed := false

	if e.postprocess.IsEnabled() {
		e.state.SetStatus(state.StatusPostProcessing)
		processed, err := e.postprocess.Process(e.ctx, text)
		if err != nil {
			e.logger.Warn(e.ctx, "post-processing failed, using raw transcription", "err", err)
		} else {
			text = processed
			postProcessed = true
		}
	}

	// Apply trailing space if enabled
	outputText := text
	if settings.OutputTrailingSpace && outputText != "" {
		outputText = outputText + " "
	}

	if err := e.writer.Write(e.ctx, settings.OutputMode, settings.PasteShortcut, outputText); err != nil {
		e.logger.Error(e.ctx, "failed to write output", "err", err)
	}

	finishedAt := time.Now()

	// Write to history
	_, err = e.historyManager.Write(e.ctx, history.WriteRequest{
		StartedAt:           startedAt,
		FinishedAt:          finishedAt,
		RecordingDurationMs: recordingDurationMs,
		TranscriptionRaw:    rawText,
		TranscriptionFinal:  text,
		PostProcessed:       postProcessed,
		AudioData:           wavData,
	})
	if err != nil {
		e.logger.Error(e.ctx, "failed to save history entry", "err", err)
	}

	e.sound.TranscriptionSuccess(e.ctx)
	e.notifier.TranscriptionFinished(e.ctx, text)
	e.state.SetStatus(state.StatusLoaded)

	// Record activity after processing completes
	e.recordActivity()

	e.logger.Info(e.ctx, "transcription complete", "length", len(text))
}

// calculateDurationMs calculates audio duration in milliseconds from raw PCM data.
// Assumes 16-bit mono audio at 16kHz sample rate.
func calculateDurationMs(dataSize int) int64 {
	// 16-bit = 2 bytes per sample, mono = 1 channel, 16kHz sample rate
	bytesPerSecond := 16000 * 2 * 1
	if bytesPerSecond == 0 {
		return 0
	}
	return int64(dataSize) * 1000 / int64(bytesPerSecond)
}

// handleError logs the error, notifies the user, plays error sound, and resets state.
func (e *Engine) handleError(message string, err error) {
	e.logger.Error(e.ctx, message, "err", err)
	e.sound.TranscriptionError(e.ctx)
	e.notifier.Error(e.ctx, config.AppName, fmt.Sprintf("%s: %v", message, err))
	e.state.SetStatus(state.StatusLoaded)
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
