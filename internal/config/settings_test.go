package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/eventbus"
)

func TestNewSettingsManager(t *testing.T) {
	tempDir := t.TempDir()
	eb := eventbus.New()

	sm, err := config.NewSettingsManager(eb, tempDir)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Check if settings file was created
	settingsPath := filepath.Join(tempDir, "settings.json")
	require.FileExists(t, settingsPath)

	// Verify default settings
	settings := sm.Get()
	require.Equal(t, "default", settings.InputDevice)
	require.Equal(t, 100, settings.SoundFeedbackVolume)
}

func TestSettingsManager_Load(t *testing.T) {
	tempDir := t.TempDir()
	eb := eventbus.New()
	settingsPath := filepath.Join(tempDir, "settings.json")

	// Create a pre-existing settings file with custom values
	initialSettings := config.Settings{
		Version:             1,
		InputDevice:         "custom-mic",
		SoundFeedbackVolume: 50,
	}
	data, err := json.Marshal(initialSettings)
	require.NoError(t, err)
	err = os.WriteFile(settingsPath, data, 0644)
	require.NoError(t, err)

	sm, err := config.NewSettingsManager(eb, tempDir)
	require.NoError(t, err)

	// Verify loaded settings
	settings := sm.Get()
	require.Equal(t, "custom-mic", settings.InputDevice)
	require.Equal(t, 50, settings.SoundFeedbackVolume)
}

func TestSettingsManager_Update(t *testing.T) {
	tempDir := t.TempDir()
	eb := eventbus.New()

	sm, err := config.NewSettingsManager(eb, tempDir)
	require.NoError(t, err)

	// Subscribe to settings changed event
	var wg sync.WaitGroup
	wg.Add(1)
	err = eb.SubscribeSettingsChanged(func() {
		wg.Done()
	})
	require.NoError(t, err)

	// Update settings
	newSettings := sm.Get()
	newSettings.InputDevice = "new-mic"
	err = sm.Update(newSettings)
	require.NoError(t, err)

	// Verify update in memory
	require.Equal(t, "new-mic", sm.Get().InputDevice)

	// Verify persistence
	data, err := os.ReadFile(filepath.Join(tempDir, "settings.json"))
	require.NoError(t, err)
	var savedSettings config.Settings
	err = json.Unmarshal(data, &savedSettings)
	require.NoError(t, err)
	require.Equal(t, "new-mic", savedSettings.InputDevice)

	// Wait for event
	wg.Wait()
}

func TestSettingsManager_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	eb := eventbus.New()

	sm, err := config.NewSettingsManager(eb, tempDir)
	require.NoError(t, err)

	var wg sync.WaitGroup
	iterations := 100

	// Concurrent Reads
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = sm.Get()
		}
	}()

	// Concurrent Updates
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			s := sm.Get()
			s.SoundFeedbackVolume = i
			_ = sm.Update(s)
		}
	}()

	wg.Wait()
}
