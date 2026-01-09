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
	"github.com/varavelio/tribar/internal/workflow"
	"github.com/varavelio/tribar/internal/workflow/executor"
)

const unloadTickerInterval = 10 * time.Second

// Dependencies contains all required dependencies for the engine.
type Dependencies struct {
	Logger           logger.Logger
	SettingsManager  *config.SettingsManager
	HistoryManager   *history.Manager
	State            *state.Instance
	Recorder         *record.Recorder
	Transcriber      *transcribe.Instance
	PostProcess      *postprocess.Instance
	Writer           *clipboard.Instance
	Notifier         *notify.Instance
	Sound            *sound.Instance
	WorkflowManager  *workflow.Manager
	WorkflowExecutor *executor.Executor
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

	// Workflow execution (Advanced Mode)
	workflowManager  *workflow.Manager
	workflowExecutor *executor.Executor

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
		logger:           deps.Logger,
		settingsManager:  deps.SettingsManager,
		historyManager:   deps.HistoryManager,
		state:            deps.State,
		recorder:         deps.Recorder,
		transcriber:      deps.Transcriber,
		postprocess:      deps.PostProcess,
		writer:           deps.Writer,
		notifier:         deps.Notifier,
		sound:            deps.Sound,
		workflowManager:  deps.WorkflowManager,
		workflowExecutor: deps.WorkflowExecutor,
		ctx:              ctx,
		cancel:           cancel,
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
// All recordings are processed through the workflow executor - either using:
// - The active workflow (when in advanced mode with a selected workflow)
// - The default workflow (standard transcribe → process → paste flow)
func (e *Engine) processRecording() {
	settings := e.settingsManager.Get()

	// Get the raw audio data from the recorder
	audioData := e.recorder.GetData()
	wavData := e.recorder.BuildWAV()

	// Determine which workflow to execute
	var wf *workflow.Workflow
	var workflowID string

	if settings.AdvancedMode && settings.ActiveWorkflowID != "" {
		// Use the user-selected workflow
		workflowID = settings.ActiveWorkflowID
		var ok bool
		wf, ok = e.workflowManager.GetByID(workflowID)
		if !ok {
			e.logger.Warn(e.ctx, "active workflow not found, using default", "workflowID", workflowID)
			wf = workflow.DefaultWorkflow()
			workflowID = "default"
		}
	} else {
		// Use the default workflow (standard pipeline)
		wf = workflow.DefaultWorkflow()
		workflowID = "default"
	}

	e.logger.Info(e.ctx, "executing workflow", "workflowID", workflowID, "workflowName", wf.Name)

	// Ensure models are loaded before transcription
	if !e.ensureModelsLoaded() {
		e.sound.TranscriptionError(e.ctx)
		e.notifier.Error(e.ctx, config.AppName, "Failed to load models for transcription")
		e.state.SetStatus(state.StatusUnloaded)
		return
	}

	e.state.SetStatus(state.StatusTranscribing)

	// Create execution context with trigger data
	triggerData := map[string]interface{}{
		"audioData": wavData,
		"rawAudio":  audioData,
	}
	execCtx := executor.NewExecutionContext(triggerData)

	// Execute the workflow
	if err := e.workflowExecutor.Execute(e.ctx, wf, execCtx); err != nil {
		e.handleError("workflow execution failed", err)
		return
	}

	e.state.SetStatus(state.StatusLoaded)
	e.recordActivity()
	e.logger.Info(e.ctx, "workflow execution complete", "workflowID", workflowID)
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
}
