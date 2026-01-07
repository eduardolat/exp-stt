package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varavelio/tribar/internal/eventbus"
)

func TestSettingsManager_Load_Defaults(t *testing.T) {
	tempDir := t.TempDir()
	bus := eventbus.New()

	// Initialize manager in a new empty directory
	sm, err := NewSettingsManager(bus, tempDir)
	require.NoError(t, err)

	// Should have default settings
	settings := sm.Get()
	require.Equal(t, defaultSettings.InputDevice, settings.InputDevice)
	require.Equal(t, defaultSettings.Version, settings.Version)
	require.NotEmpty(t, settings.Prompts)

	// File should have been created
	require.FileExists(t, filepath.Join(tempDir, settingsFileName))
}

func TestSettingsManager_Save_Load(t *testing.T) {
	tempDir := t.TempDir()
	bus := eventbus.New()
	sm, err := NewSettingsManager(bus, tempDir)
	require.NoError(t, err)

	// Modify settings
	newSettings := sm.Get()
	newSettings.InputDevice = "new-device-id"
	newSettings.SoundFeedbackVolume = 50

	// Save
	err = sm.Update(newSettings)
	require.NoError(t, err)

	// Create new manager from same dir
	sm2, err := NewSettingsManager(bus, tempDir)
	require.NoError(t, err)

	loadedSettings := sm2.Get()
	require.Equal(t, "new-device-id", loadedSettings.InputDevice)
	require.Equal(t, 50, loadedSettings.SoundFeedbackVolume)
}

func TestSettingsManager_Update(t *testing.T) {
	tempDir := t.TempDir()
	bus := eventbus.New()
	sm, err := NewSettingsManager(bus, tempDir)
	require.NoError(t, err)

	// Subscribe to changes
	changes := make(chan struct{}, 1)
	_ = bus.SubscribeSettingsChanged(func() {
		changes <- struct{}{}
	})

	// Update settings
	newSettings := sm.Get()
	newSettings.NotifyOnStart = true
	err = sm.Update(newSettings)
	require.NoError(t, err)

	// Verify update in memory
	require.True(t, sm.Get().NotifyOnStart)

	// Verify event fired
	select {
	case <-changes:
		// success
	default:
		t.Fatal("settings changed event not fired")
	}

	// Verify persistence
	data, err := os.ReadFile(filepath.Join(tempDir, settingsFileName))
	require.NoError(t, err)
	var saved Settings
	err = json.Unmarshal(data, &saved)
	require.NoError(t, err)
	require.True(t, saved.NotifyOnStart)
}

func TestSettingsManager_Load_Merge(t *testing.T) {
	tempDir := t.TempDir()
	bus := eventbus.New()

	// Create an existing settings file with ONLY one field (simulating old version or partial file)
	partialSettings := `{"input_device": "old-device"}`
	err := os.WriteFile(filepath.Join(tempDir, settingsFileName), []byte(partialSettings), 0644)
	require.NoError(t, err)

	// Initialize manager
	sm, err := NewSettingsManager(bus, tempDir)
	require.NoError(t, err)

	settings := sm.Get()
	// Should preserve existing value
	require.Equal(t, "old-device", settings.InputDevice)
	// Should populate missing fields with defaults
	require.Equal(t, defaultSettings.SoundFeedbackVolume, settings.SoundFeedbackVolume)
	require.Equal(t, defaultSettings.HistoryLimit, settings.HistoryLimit)

	// Should have re-saved the full file
	data, err := os.ReadFile(filepath.Join(tempDir, settingsFileName))
	require.NoError(t, err)
	var saved Settings
	err = json.Unmarshal(data, &saved)
	require.NoError(t, err)
	require.Equal(t, defaultSettings.HistoryLimit, saved.HistoryLimit)
}

func TestSettingsManager_Load_Corrupted(t *testing.T) {
	tempDir := t.TempDir()
	bus := eventbus.New()

	// Create a corrupted settings file
	err := os.WriteFile(filepath.Join(tempDir, settingsFileName), []byte("{ invalid json"), 0644)
	require.NoError(t, err)

	// Initialize manager should fail
	_, err = NewSettingsManager(bus, tempDir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse settings")
}
