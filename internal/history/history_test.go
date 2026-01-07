package history

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/varavelio/tribar/internal/config"
)

// MockSettingsProvider mocks the settings configuration for testing.
type MockSettingsProvider struct {
	Settings config.Settings
}

func (m *MockSettingsProvider) Get() config.Settings {
	return m.Settings
}

// MockLogger is a no-op logger for testing.
type MockLogger struct{}

func (l *MockLogger) SetDebug(enabled bool)                                               {}
func (l *MockLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {}
func (l *MockLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Warn(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {}

func TestHistoryManager(t *testing.T) {
	// Setup temporary directory for history files
	tempDir := t.TempDir()

	// Setup mock settings with a history limit of 2
	mockSettings := &MockSettingsProvider{
		Settings: config.Settings{
			HistoryLimit: 2,
		},
	}

	// Create manager
	manager := NewManager(&MockLogger{}, mockSettings, tempDir)
	ctx := context.Background()

	// 1. Test Write
	req1 := WriteRequest{
		StartedAt:           time.Now().Add(-2 * time.Minute),
		FinishedAt:          time.Now().Add(-1 * time.Minute),
		RecordingDurationMs: 60000,
		TranscriptionRaw:    "Test Transcription 1",
		TranscriptionFinal:  "Test Transcription 1 Final",
		PostProcessed:       true,
		AudioData:           []byte("audio data 1"),
	}

	entry1, err := manager.Write(ctx, req1)
	require.NoError(t, err)
	require.NotEmpty(t, entry1.ID)
	require.Equal(t, req1.TranscriptionRaw, entry1.TranscriptionRaw)

	// Verify files exist
	jsonPath1 := filepath.Join(tempDir, entry1.ID+".json")
	wavPath1 := filepath.Join(tempDir, entry1.ID+".wav")
	require.FileExists(t, jsonPath1)
	require.FileExists(t, wavPath1)

	// Verify content of JSON
	data, err := os.ReadFile(jsonPath1)
	require.NoError(t, err)
	var loadedEntry1 Entry
	err = json.Unmarshal(data, &loadedEntry1)
	require.NoError(t, err)
	require.Equal(t, entry1.ID, loadedEntry1.ID)

	// 2. Test GetByID
	gotEntry, found := manager.GetByID(entry1.ID)
	require.True(t, found)
	require.Equal(t, entry1.ID, gotEntry.ID)

	// 3. Test Pruning (add 2 more entries, limit is 2, so entry1 should be removed)
	req2 := WriteRequest{
		TranscriptionRaw: "Test Transcription 2",
		AudioData:        []byte("audio data 2"),
	}
	entry2, err := manager.Write(ctx, req2)
	require.NoError(t, err)

	// Ensure different timestamps for stable sorting if needed, but UUIDv7 is time-ordered
	time.Sleep(1 * time.Millisecond)

	req3 := WriteRequest{
		TranscriptionRaw: "Test Transcription 3",
		AudioData:        []byte("audio data 3"),
	}
	entry3, err := manager.Write(ctx, req3)
	require.NoError(t, err)

	// Now we should have entry2 and entry3. Entry1 should be pruned.
	require.Equal(t, 2, manager.Count())
	entries := manager.GetAll(ctx)
	require.Len(t, entries, 2)

	// entries are sorted newest first
	require.Equal(t, entry3.ID, entries[0].ID)
	require.Equal(t, entry2.ID, entries[1].ID)

	// Verify entry1 files are gone
	require.NoFileExists(t, jsonPath1)
	require.NoFileExists(t, wavPath1)

	// 4. Test Load (create new manager pointing to same dir)
	manager2 := NewManager(&MockLogger{}, mockSettings, tempDir)
	err = manager2.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, manager2.Count())

	loadedEntries := manager2.GetAll(ctx)
	// Sort order check: newest first
	require.Equal(t, entry3.ID, loadedEntries[0].ID)
	require.Equal(t, entry2.ID, loadedEntries[1].ID)

	// 5. Test Delete
	err = manager.Delete(ctx, entry2.ID)
	require.NoError(t, err)
	require.Equal(t, 1, manager.Count())
	require.NoFileExists(t, filepath.Join(tempDir, entry2.ID+".json"))

	// 6. Test Clear
	err = manager.Clear(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, manager.Count())
	require.NoFileExists(t, filepath.Join(tempDir, entry3.ID+".json"))

	// 7. Test GetAudioPath
	path := manager.GetAudioPath("some-id")
	require.Equal(t, filepath.Join(tempDir, "some-id.wav"), path)
}

func TestHistoryManager_Load_Corrupted(t *testing.T) {
	tempDir := t.TempDir()
	mockSettings := &MockSettingsProvider{Settings: config.Settings{HistoryLimit: 10}}
	manager := NewManager(&MockLogger{}, mockSettings, tempDir)
	ctx := context.Background()

	// Create a corrupted JSON file
	err := os.WriteFile(filepath.Join(tempDir, "bad.json"), []byte("{ invalid json"), 0644)
	require.NoError(t, err)

	// Create a valid entry manually
	validID := "01943d2c-8c1a-7b3b-8267-8f5b8c6d1e4f" // Example UUIDv7-like
	validEntry := Entry{
		Version: 1,
		ID:      validID,
	}
	validJson, _ := json.Marshal(validEntry)
	err = os.WriteFile(filepath.Join(tempDir, validID+".json"), validJson, 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, validID+".wav"), []byte("audio"), 0644)
	require.NoError(t, err)

	// Load should succeed but skip corrupted file
	err = manager.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, manager.Count())
	entries := manager.GetAll(ctx)
	require.Equal(t, validID, entries[0].ID)
}

func TestHistoryManager_Load_MissingWav(t *testing.T) {
	tempDir := t.TempDir()
	mockSettings := &MockSettingsProvider{Settings: config.Settings{HistoryLimit: 10}}
	manager := NewManager(&MockLogger{}, mockSettings, tempDir)
	ctx := context.Background()

	// Create a valid JSON but no WAV
	validID := "01943d2c-8c1a-7b3b-8267-8f5b8c6d1e4f"
	validEntry := Entry{Version: 1, ID: validID}
	validJson, _ := json.Marshal(validEntry)
	err := os.WriteFile(filepath.Join(tempDir, validID+".json"), validJson, 0644)
	require.NoError(t, err)

	// Load should succeed but skip entry because wav is missing
	err = manager.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, manager.Count())
}

func TestHistoryManager_Sort(t *testing.T) {
	tempDir := t.TempDir()
	mockSettings := &MockSettingsProvider{Settings: config.Settings{HistoryLimit: 10}}
	manager := NewManager(&MockLogger{}, mockSettings, tempDir)
	ctx := context.Background()

	// Create 3 entries with IDs that sort alphabetically (which correlates to time in UUIDv7)
	// We'll simulate UUIDv7 strings for sorting property
	ids := []string{
		"01943d2c-1111-7b3b-8267-8f5b8c6d1e4f", // oldest
		"01943d2c-2222-7b3b-8267-8f5b8c6d1e4f", // middle
		"01943d2c-3333-7b3b-8267-8f5b8c6d1e4f", // newest
	}

	for _, id := range ids {
		entry := Entry{Version: 1, ID: id}
		jsonBytes, _ := json.Marshal(entry)
		_ = os.WriteFile(filepath.Join(tempDir, id+".json"), jsonBytes, 0644)
		_ = os.WriteFile(filepath.Join(tempDir, id+".wav"), []byte("wav"), 0644)
	}

	err := manager.Load(ctx)
	require.NoError(t, err)

	entries := manager.GetAll(ctx)
	require.Len(t, entries, 3)
	// Should be sorted newest first (descending ID)
	require.Equal(t, ids[2], entries[0].ID)
	require.Equal(t, ids[1], entries[1].ID)
	require.Equal(t, ids[0], entries[2].ID)
}
