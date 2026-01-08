package state

import (
	"context"
	"errors"
	"testing"

	"github.com/gen2brain/malgo"
	"github.com/stretchr/testify/require"
	"github.com/varavelio/tribar/internal/eventbus"
	"github.com/varavelio/tribar/internal/history"
)

// MockAudioContext implements AudioContext for testing
type MockAudioContext struct {
	DevicesFunc func(kind malgo.DeviceType) ([]malgo.DeviceInfo, error)
}

func (m *MockAudioContext) Devices(kind malgo.DeviceType) ([]malgo.DeviceInfo, error) {
	if m.DevicesFunc != nil {
		return m.DevicesFunc(kind)
	}
	return nil, nil
}

// MockHistoryManager implements HistoryManager for testing
type MockHistoryManager struct {
	GetAllFunc       func(ctx context.Context) []history.Entry
	GetByIDFunc      func(id string) (history.Entry, bool)
	GetAudioPathFunc func(id string) string
	DeleteFunc       func(ctx context.Context, id string) error
	ClearFunc        func(ctx context.Context) error
	CountFunc        func() int
}

func (m *MockHistoryManager) GetAll(ctx context.Context) []history.Entry {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil
}

func (m *MockHistoryManager) GetByID(id string) (history.Entry, bool) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return history.Entry{}, false
}

func (m *MockHistoryManager) GetAudioPath(id string) string {
	if m.GetAudioPathFunc != nil {
		return m.GetAudioPathFunc(id)
	}
	return ""
}

func (m *MockHistoryManager) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockHistoryManager) Clear(ctx context.Context) error {
	if m.ClearFunc != nil {
		return m.ClearFunc(ctx)
	}
	return nil
}

func (m *MockHistoryManager) Count() int {
	if m.CountFunc != nil {
		return m.CountFunc()
	}
	return 0
}

func TestStateStatus(t *testing.T) {
	eb := eventbus.New()
	st := NewWithDependencies(eb, &MockHistoryManager{}, &MockAudioContext{})

	// Test Initial Status
	current, previous := st.GetStatus()
	require.Equal(t, StatusUnloaded, current)
	require.Equal(t, StatusUnknown, previous)

	// Subscribe to event bus
	stateChanged := false
	eb.SubscribeStateChanged(func() {
		stateChanged = true
	})

	// Test SetStatus
	st.SetStatus(StatusLoading)

	// Use assert for concurrent safety if needed, though here we are single threaded mostly
	// But memory says: "Tests for event-driven logic... require synchronization"
	// However, here the callback is executed synchronously by the event bus usually?
	// Let's check eventbus implementation briefly. Assuming standard synchronous dispatch or we need to wait.
	// We'll assume sync dispatch for now, but if it fails we add sync.

	require.True(t, stateChanged, "StateChanged event should have been fired")

	current, previous = st.GetStatus()
	require.Equal(t, StatusLoading, current)
	require.Equal(t, StatusUnloaded, previous)
}

func TestDownloadProgress(t *testing.T) {
	eb := eventbus.New()
	st := NewWithDependencies(eb, &MockHistoryManager{}, &MockAudioContext{})

	eventCount := 0
	eb.SubscribeStateChanged(func() {
		eventCount++
	})

	// 0% -> Should fire
	st.SetDownloadProgress("test.bin", 0, 100, 0)
	require.Equal(t, 1, eventCount)

	// 2% -> Should NOT fire (only every 5%)
	st.SetDownloadProgress("test.bin", 2, 100, 2)
	require.Equal(t, 1, eventCount)

	// 5% -> Should fire
	st.SetDownloadProgress("test.bin", 5, 100, 5)
	require.Equal(t, 2, eventCount)

	// Check Values
	prog := st.GetDownloadProgress()
	require.Equal(t, "test.bin", prog.FileName)
	require.Equal(t, float64(5), prog.Percent)

	// Clear
	st.ClearDownloadProgress()
	prog = st.GetDownloadProgress()
	require.Equal(t, "", prog.FileName)
	require.Equal(t, 3, eventCount) // Clear fires event
}

func TestGetAvailableDevices(t *testing.T) {
	mockAudio := &MockAudioContext{
		DevicesFunc: func(kind malgo.DeviceType) ([]malgo.DeviceInfo, error) {
			if kind == malgo.Capture {
				return []malgo.DeviceInfo{
					{}, // Empty device 1
					{}, // Empty device 2
				}, nil
			}
			return []malgo.DeviceInfo{
				{}, // Output device 1
			}, nil
		},
	}

	st := NewWithDependencies(eventbus.New(), &MockHistoryManager{}, mockAudio)

	devices := st.GetAvailableDevices()

	require.Len(t, devices.InputDevices, 2)
	require.Len(t, devices.OutputDevices, 1)

	// Verify IsDefault logic
	require.True(t, devices.InputDevices[0].IsDefault)
	require.False(t, devices.InputDevices[1].IsDefault)

	// Error handling
	mockAudio.DevicesFunc = func(kind malgo.DeviceType) ([]malgo.DeviceInfo, error) {
		return nil, errors.New("fail")
	}
	devices = st.GetAvailableDevices()
	require.Empty(t, devices.InputDevices)
	require.Empty(t, devices.OutputDevices)
}

func TestHistoryOperations(t *testing.T) {
	mockHist := &MockHistoryManager{}
	eb := eventbus.New()
	st := NewWithDependencies(eb, mockHist, &MockAudioContext{})

	// Test GetHistory
	mockHist.GetAllFunc = func(ctx context.Context) []history.Entry {
		return []history.Entry{{ID: "1"}}
	}
	entries := st.GetHistory(context.Background())
	require.Len(t, entries, 1)
	require.Equal(t, "1", entries[0].ID)

	// Test GetHistoryEntry
	mockHist.GetByIDFunc = func(id string) (history.Entry, bool) {
		if id == "1" {
			return history.Entry{ID: "1"}, true
		}
		return history.Entry{}, false
	}
	entry, found := st.GetHistoryEntry("1")
	require.True(t, found)
	require.Equal(t, "1", entry.ID)

	// Test DeleteHistoryEntry
	eventFired := false
	eb.SubscribeStateChanged(func() { eventFired = true })

	mockHist.DeleteFunc = func(ctx context.Context, id string) error {
		return nil
	}
	err := st.DeleteHistoryEntry(context.Background(), "1")
	require.NoError(t, err)
	require.True(t, eventFired)

	// Test Delete failure
	eventFired = false
	mockHist.DeleteFunc = func(ctx context.Context, id string) error {
		return errors.New("oops")
	}
	err = st.DeleteHistoryEntry(context.Background(), "1")
	require.Error(t, err)
	require.False(t, eventFired) // Should not fire on error

	// Test ClearHistory
	eventFired = false
	mockHist.ClearFunc = func(ctx context.Context) error {
		return nil
	}
	err = st.ClearHistory(context.Background())
	require.NoError(t, err)
	require.True(t, eventFired)

	// Test Count
	mockHist.CountFunc = func() int { return 42 }
	require.Equal(t, 42, st.HistoryCount())

	// Test GetAudioPath
	mockHist.GetAudioPathFunc = func(id string) string { return "/tmp/" + id + ".wav" }
	require.Equal(t, "/tmp/123.wav", st.GetHistoryAudioPath("123"))
}
