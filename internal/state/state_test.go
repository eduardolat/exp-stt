package state

import (
	"errors"
	"sync"
	"testing"

	"github.com/gen2brain/malgo"
	"github.com/stretchr/testify/require"
	"github.com/varavelio/tribar/internal/eventbus"
)

// MockDeviceInfo implements DeviceInfo interface
type MockDeviceInfo struct {
	ID   string
	Name string
}

func (m MockDeviceInfo) IDString() string {
	return m.ID
}

func (m MockDeviceInfo) NameString() string {
	return m.Name
}

// MockAudioContext implements AudioContext interface
type MockAudioContext struct {
	CaptureDevices  []DeviceInfo
	PlaybackDevices []DeviceInfo
	Err             error
}

func (m *MockAudioContext) Devices(deviceType malgo.DeviceType) ([]DeviceInfo, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if deviceType == malgo.Capture {
		return m.CaptureDevices, nil
	}
	if deviceType == malgo.Playback {
		return m.PlaybackDevices, nil
	}
	return nil, nil
}

func TestGetAvailableDevices(t *testing.T) {
	eb := eventbus.New()

	t.Run("Happy Path", func(t *testing.T) {
		mockCtx := &MockAudioContext{
			CaptureDevices: []DeviceInfo{
				MockDeviceInfo{ID: "mic1", Name: "Microphone 1"},
				MockDeviceInfo{ID: "mic2", Name: "Microphone 2"},
			},
			PlaybackDevices: []DeviceInfo{
				MockDeviceInfo{ID: "spk1", Name: "Speaker 1"},
			},
		}

		inst := NewWithContext(eb, nil, mockCtx)
		devices := inst.GetAvailableDevices()

		require.Len(t, devices.InputDevices, 2)
		require.Len(t, devices.OutputDevices, 1)

		// Check Input Devices
		require.Equal(t, "mic1", devices.InputDevices[0].ID)
		require.Equal(t, "Microphone 1", devices.InputDevices[0].Name)
		require.True(t, devices.InputDevices[0].IsDefault, "First device should be default")

		require.Equal(t, "mic2", devices.InputDevices[1].ID)
		require.Equal(t, "Microphone 2", devices.InputDevices[1].Name)
		require.False(t, devices.InputDevices[1].IsDefault)

		// Check Output Devices
		require.Equal(t, "spk1", devices.OutputDevices[0].ID)
		require.True(t, devices.OutputDevices[0].IsDefault)
	})

	t.Run("Empty Devices", func(t *testing.T) {
		mockCtx := &MockAudioContext{
			CaptureDevices:  []DeviceInfo{},
			PlaybackDevices: []DeviceInfo{},
		}

		inst := NewWithContext(eb, nil, mockCtx)
		devices := inst.GetAvailableDevices()

		require.Empty(t, devices.InputDevices)
		require.Empty(t, devices.OutputDevices)
	})

	t.Run("Error fetching devices", func(t *testing.T) {
		mockCtx := &MockAudioContext{
			Err: errors.New("malgo error"),
		}

		inst := NewWithContext(eb, nil, mockCtx)
		devices := inst.GetAvailableDevices()

		require.Empty(t, devices.InputDevices)
		require.Empty(t, devices.OutputDevices)
	})

	t.Run("Nil Context (Defensive)", func(t *testing.T) {
		// NewWithContext allows passing nil, though our constructor doesn't typically allow it.
		// If audioCtx is nil, GetAvailableDevices should return empty.
		inst := NewWithContext(eb, nil, nil)
		devices := inst.GetAvailableDevices()

		require.Empty(t, devices.InputDevices)
		require.Empty(t, devices.OutputDevices)
	})
}

func TestStatusTransitions(t *testing.T) {
	eb := eventbus.New()
	inst := NewWithContext(eb, nil, &MockAudioContext{})

	// Initial State
	current, prev := inst.GetStatus()
	require.Equal(t, StatusUnloaded, current)
	require.Equal(t, StatusUnknown, prev)

	// Subscribe to verify event emission
	var wg sync.WaitGroup
	wg.Add(1)
	err := eb.SubscribeStateChanged(func() {
		wg.Done()
	})
	require.NoError(t, err)

	// Transition
	inst.SetStatus(StatusLoading)

	wg.Wait() // Wait for event

	current, prev = inst.GetStatus()
	require.Equal(t, StatusLoading, current)
	require.Equal(t, StatusUnloaded, prev)
}

func TestDownloadProgress(t *testing.T) {
	eb := eventbus.New()
	inst := NewWithContext(eb, nil, &MockAudioContext{})

	var callCount int
	var mu sync.Mutex

	err := eb.SubscribeStateChanged(func() {
		mu.Lock()
		callCount++
		mu.Unlock()
	})
	require.NoError(t, err)

	// Set progress to 4% (Should NOT trigger event)
	inst.SetDownloadProgress("model.bin", 4, 100, 4.0)
	mu.Lock()
	require.Equal(t, 0, callCount)
	mu.Unlock()

	// Set progress to 5% (Should trigger event)
	inst.SetDownloadProgress("model.bin", 5, 100, 5.0)
	mu.Lock()
	require.Equal(t, 1, callCount)
	mu.Unlock()

	// Set progress to 6% (Should NOT trigger event)
	inst.SetDownloadProgress("model.bin", 6, 100, 6.0)
	mu.Lock()
	require.Equal(t, 1, callCount)
	mu.Unlock()

	// Set progress to 10% (Should trigger event)
	inst.SetDownloadProgress("model.bin", 10, 100, 10.0)
	mu.Lock()
	require.Equal(t, 2, callCount)
	mu.Unlock()
}

func TestClearDownloadProgress(t *testing.T) {
	eb := eventbus.New()
	inst := NewWithContext(eb, nil, &MockAudioContext{})

	inst.SetDownloadProgress("file", 50, 100, 50.0)
	progress := inst.GetDownloadProgress()
	require.Equal(t, "file", progress.FileName)

	inst.ClearDownloadProgress()
	progress = inst.GetDownloadProgress()
	require.Empty(t, progress.FileName)
	require.Equal(t, int64(0), progress.Downloaded)
}
