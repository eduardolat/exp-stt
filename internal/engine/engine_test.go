package engine

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/history"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/state"
	"github.com/varavelio/tribar/internal/transcribe"
)

// Mock implementations
type MockSettingsManager struct {
	mock.Mock
}

func (m *MockSettingsManager) Get() config.Settings {
	args := m.Called()
	return args.Get(0).(config.Settings)
}

type MockHistoryManager struct {
	mock.Mock
}

func (m *MockHistoryManager) Write(ctx context.Context, req history.WriteRequest) (history.Entry, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(history.Entry), args.Error(1)
}

type MockStateManager struct {
	mock.Mock
}

func (m *MockStateManager) SetStatus(status state.Status) {
	m.Called(status)
}

func (m *MockStateManager) GetStatus() (state.Status, state.Status) {
	args := m.Called()
	return args.Get(0).(state.Status), args.Get(1).(state.Status)
}

func (m *MockStateManager) SetDownloadProgress(fileName string, downloaded, total int64, percent float64) {
	m.Called(fileName, downloaded, total, percent)
}

func (m *MockStateManager) ClearDownloadProgress() {
	m.Called()
}

type MockRecorder struct {
	mock.Mock
}

func (m *MockRecorder) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRecorder) Stop() {
	m.Called()
}

func (m *MockRecorder) GetData() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *MockRecorder) BuildWAV() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

type MockTranscriber struct {
	mock.Mock
}

func (m *MockTranscriber) CheckModels() (bool, []transcribe.ModelFile) {
	args := m.Called()
	return args.Bool(0), args.Get(1).([]transcribe.ModelFile)
}

func (m *MockTranscriber) DownloadModels(progress transcribe.DownloadProgressCallback) error {
	args := m.Called(progress)
	// Execute the callback if provided in args/stub setup, here we just return
	return args.Error(0)
}

func (m *MockTranscriber) LoadModels() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTranscriber) UnloadModels() {
	m.Called()
}

func (m *MockTranscriber) TranscribeWAV(wavData []byte) (string, error) {
	args := m.Called(wavData)
	return args.String(0), args.Error(1)
}

type MockPostProcess struct {
	mock.Mock
}

func (m *MockPostProcess) IsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockPostProcess) Process(ctx context.Context, text string) (string, error) {
	args := m.Called(ctx, text)
	return args.String(0), args.Error(1)
}

type MockClipboardWriter struct {
	mock.Mock
}

func (m *MockClipboardWriter) Write(ctx context.Context, mode config.OutputMode, pasteShortcut string, text string) error {
	args := m.Called(ctx, mode, pasteShortcut, text)
	return args.Error(0)
}

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) Error(ctx context.Context, title string, message string) {
	m.Called(ctx, title, message)
}

func (m *MockNotifier) TranscriptionStarted(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockNotifier) TranscriptionFinished(ctx context.Context, text string) {
	m.Called(ctx, text)
}

type MockSoundManager struct {
	mock.Mock
}

func (m *MockSoundManager) RecordingStarted(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockSoundManager) RecordingStopped(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockSoundManager) TranscriptionError(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockSoundManager) TranscriptionSuccess(ctx context.Context) {
	m.Called(ctx)
}

// Helper to create engine with mocks
func createTestEngine() (*Engine, *Dependencies) {
	deps := &Dependencies{
		Logger:          &logger.NoOpLogger{},
		SettingsManager: new(MockSettingsManager),
		HistoryManager:  new(MockHistoryManager),
		State:           new(MockStateManager),
		Recorder:        new(MockRecorder),
		Transcriber:     new(MockTranscriber),
		PostProcess:     new(MockPostProcess),
		Writer:          new(MockClipboardWriter),
		Notifier:        new(MockNotifier),
		Sound:           new(MockSoundManager),
	}

	return New(*deps), deps
}

func TestToggleRecording(t *testing.T) {
	t.Run("Starts recording when status is StatusLoaded", func(t *testing.T) {
		e, deps := createTestEngine()
		mockState := deps.State.(*MockStateManager)
		mockRecorder := deps.Recorder.(*MockRecorder)
		mockSound := deps.Sound.(*MockSoundManager)
		mockNotifier := deps.Notifier.(*MockNotifier)

		// Setup expectations
		mockState.On("GetStatus").Return(state.StatusLoaded, state.StatusUnloaded)
		mockRecorder.On("Start").Return(nil)
		mockState.On("SetStatus", state.StatusListening).Return()
		mockSound.On("RecordingStarted", mock.Anything).Return()
		mockNotifier.On("TranscriptionStarted", mock.Anything).Return()

		e.ToggleRecording()

		mockState.AssertExpectations(t)
		mockRecorder.AssertExpectations(t)
		mockSound.AssertExpectations(t)
		mockNotifier.AssertExpectations(t)
	})

	t.Run("Stops recording when status is StatusListening", func(t *testing.T) {
		e, deps := createTestEngine()
		mockState := deps.State.(*MockStateManager)
		mockRecorder := deps.Recorder.(*MockRecorder)
		mockSound := deps.Sound.(*MockSoundManager)
		mockTranscriber := deps.Transcriber.(*MockTranscriber)
		mockSettings := deps.SettingsManager.(*MockSettingsManager)
		mockWriter := deps.Writer.(*MockClipboardWriter)
		mockHistory := deps.HistoryManager.(*MockHistoryManager)
		mockNotifier := deps.Notifier.(*MockNotifier)
		mockPostProcess := deps.PostProcess.(*MockPostProcess)

		// Setup expectations for ToggleRecording -> stopRecording
		mockState.On("GetStatus").Return(state.StatusListening, state.StatusLoaded)
		mockRecorder.On("Stop").Return()
		mockSound.On("RecordingStopped", mock.Anything).Return()

		// Expectations for processRecording (async)
		mockSettings.On("Get").Return(config.Settings{
			OutputMode:          config.OutputModeCopyOnly,
			PasteShortcut:       "ctrl+v",
			OutputTrailingSpace: false,
		})
		mockRecorder.On("GetData").Return([]byte{1, 2, 3, 4})
		mockRecorder.On("BuildWAV").Return([]byte{1, 2, 3, 4})

		// EnsureModelsLoaded logic
		mockState.On("SetStatus", state.StatusLoading).Return() // if models not loaded
		mockTranscriber.On("LoadModels").Return(nil)

		mockState.On("SetStatus", state.StatusTranscribing).Return()
		mockTranscriber.On("TranscribeWAV", []byte{1, 2, 3, 4}).Return("Hello World", nil)

		mockPostProcess.On("IsEnabled").Return(false)

		mockWriter.On("Write", mock.Anything, config.OutputModeCopyOnly, "ctrl+v", "Hello World").Return(nil)

		mockHistory.On("Write", mock.Anything, mock.MatchedBy(func(req history.WriteRequest) bool {
			return req.TranscriptionFinal == "Hello World"
		})).Return(history.Entry{}, nil)

		mockSound.On("TranscriptionSuccess", mock.Anything).Return()
		mockNotifier.On("TranscriptionFinished", mock.Anything, "Hello World").Return()
		mockState.On("SetStatus", state.StatusLoaded).Return()

		e.ToggleRecording()

		// Wait for goroutine to finish
		time.Sleep(100 * time.Millisecond)

		mockState.AssertExpectations(t)
		mockRecorder.AssertExpectations(t)
		mockSound.AssertExpectations(t)
		mockTranscriber.AssertExpectations(t)
		mockWriter.AssertExpectations(t)
		mockHistory.AssertExpectations(t)
	})

	t.Run("Fails to start recording if StatusDownloading", func(t *testing.T) {
		e, deps := createTestEngine()
		mockState := deps.State.(*MockStateManager)
		mockNotifier := deps.Notifier.(*MockNotifier)

		mockState.On("GetStatus").Return(state.StatusDownloading, state.StatusUnloaded)
		mockNotifier.On("Error", mock.Anything, config.AppName, "Please wait, models are being downloaded").Return()

		e.ToggleRecording()

		mockState.AssertExpectations(t)
		mockNotifier.AssertExpectations(t)
	})
}

func TestEnsureModelsDownloaded(t *testing.T) {
	t.Run("Models already downloaded", func(t *testing.T) {
		e, deps := createTestEngine()
		mockTranscriber := deps.Transcriber.(*MockTranscriber)

		mockTranscriber.On("CheckModels").Return(true, []transcribe.ModelFile{})

		err := e.EnsureModelsDownloaded()

		assert.NoError(t, err)
		mockTranscriber.AssertExpectations(t)
	})

	t.Run("Downloads missing models", func(t *testing.T) {
		e, deps := createTestEngine()
		mockTranscriber := deps.Transcriber.(*MockTranscriber)
		mockState := deps.State.(*MockStateManager)

		mockTranscriber.On("CheckModels").Return(false, []transcribe.ModelFile{{Name: "model1"}})
		mockState.On("SetStatus", state.StatusDownloading).Return()
		mockTranscriber.On("DownloadModels", mock.Anything).Return(nil)
		mockState.On("ClearDownloadProgress").Return()
		mockState.On("SetStatus", state.StatusUnloaded).Return()

		err := e.EnsureModelsDownloaded()

		assert.NoError(t, err)
		mockTranscriber.AssertExpectations(t)
		mockState.AssertExpectations(t)
	})

	t.Run("Download fails", func(t *testing.T) {
		e, deps := createTestEngine()
		mockTranscriber := deps.Transcriber.(*MockTranscriber)
		mockState := deps.State.(*MockStateManager)
		mockNotifier := deps.Notifier.(*MockNotifier)

		expectedErr := errors.New("network error")
		mockTranscriber.On("CheckModels").Return(false, []transcribe.ModelFile{{Name: "model1"}})
		mockState.On("SetStatus", state.StatusDownloading).Return()
		mockTranscriber.On("DownloadModels", mock.Anything).Return(expectedErr)
		mockState.On("SetStatus", state.StatusUnloaded).Return()
		mockState.On("ClearDownloadProgress").Return()
		mockNotifier.On("Error", mock.Anything, "Model Download Failed", expectedErr.Error()).Return()

		err := e.EnsureModelsDownloaded()

		assert.Error(t, err)
		mockTranscriber.AssertExpectations(t)
		mockState.AssertExpectations(t)
		mockNotifier.AssertExpectations(t)
	})
}

func TestCheckAndUnload(t *testing.T) {
	t.Run("Unloads models when inactive", func(t *testing.T) {
		e, deps := createTestEngine()
		mockSettings := deps.SettingsManager.(*MockSettingsManager)
		mockState := deps.State.(*MockStateManager)
		mockTranscriber := deps.Transcriber.(*MockTranscriber)

		// Set manually internal state
		e.modelsLoaded = true
		// Hack to set lastActivityTime way back
		// We can't set unexported field easily without reflection or exposing it.
		// But verify logic: if we don't touch recordActivity, it's time.Time{} (zero), so elapsed is huge.

		mockSettings.On("Get").Return(config.Settings{
			ModelUnloadEnable:  true,
			ModelUnloadSeconds: 10,
		})

		mockState.On("GetStatus").Return(state.StatusLoaded, state.StatusLoaded)
		mockTranscriber.On("UnloadModels").Return()
		mockState.On("SetStatus", state.StatusUnloaded).Return()

		e.checkAndUnload()

		mockTranscriber.AssertExpectations(t)
		mockState.AssertExpectations(t)
	})

	t.Run("Does not unload if disabled", func(t *testing.T) {
		e, deps := createTestEngine()
		mockSettings := deps.SettingsManager.(*MockSettingsManager)

		mockSettings.On("Get").Return(config.Settings{
			ModelUnloadEnable: false,
		})

		e.checkAndUnload()

		mockSettings.AssertExpectations(t)
	})
}
