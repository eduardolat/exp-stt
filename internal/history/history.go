// Package history manages transcription history using a flat-file database approach.
// Each transcription is stored as a pair of files: [UUIDv7].wav and [UUIDv7].json.
// The history is loaded dynamically at startup and kept in memory for fast access.
package history

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

const (
	currentSchemaVersion = 1
	jsonExtension        = ".json"
	wavExtension         = ".wav"
)

// Entry represents a single transcription record stored on disk.
type Entry struct {
	Version             int       `json:"version"`
	ID                  string    `json:"id"`
	StartedAt           time.Time `json:"started_at"`
	FinishedAt          time.Time `json:"finished_at"`
	RecordingDurationMs int64     `json:"recording_duration_ms"`
	TranscriptionRaw    string    `json:"transcription_raw"`
	TranscriptionFinal  string    `json:"transcription_final"`
	PostProcessed       bool      `json:"post_processed"`
}

// WriteRequest contains the data needed to create a new history entry.
type WriteRequest struct {
	StartedAt           time.Time
	FinishedAt          time.Time
	RecordingDurationMs int64
	TranscriptionRaw    string
	TranscriptionFinal  string
	PostProcessed       bool
	AudioData           []byte
}

// Manager handles all history operations including loading, saving, and pruning.
type Manager struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager
	directory       string

	mu      sync.RWMutex
	entries []Entry
}

// NewManager creates a new history manager.
func NewManager(logger logger.Logger, settingsManager *config.SettingsManager) *Manager {
	return &Manager{
		logger:          logger,
		settingsManager: settingsManager,
		directory:       config.DirectoryRecordings,
		entries:         make([]Entry, 0),
	}
}

// LoadAsync loads the history from disk in a separate goroutine.
// This should be called during application startup to avoid blocking.
func (m *Manager) LoadAsync(ctx context.Context) {
	go func() {
		if err := m.Load(ctx); err != nil {
			m.logger.Error(ctx, "failed to load history", "err", err)
		}
	}()
}

// Load reads all history entries from disk and populates the in-memory cache.
func (m *Manager) Load(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.loadLocked(ctx)
}

// loadLocked performs the actual loading. Caller must hold m.mu lock.
func (m *Manager) loadLocked(ctx context.Context) error {
	files, err := os.ReadDir(m.directory)
	if err != nil {
		if os.IsNotExist(err) {
			m.entries = make([]Entry, 0)
			return nil
		}
		return fmt.Errorf("failed to read history directory: %w", err)
	}

	entries := make([]Entry, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), jsonExtension) {
			continue
		}

		entry, err := m.loadEntry(ctx, file.Name())
		if err != nil {
			m.logger.Warn(ctx, "skipping corrupted history entry",
				"file", file.Name(),
				"err", err,
			)
			continue
		}

		entries = append(entries, entry)
	}

	// Sort by ID descending (UUIDv7 is K-sortable, so alphabetical = chronological)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID > entries[j].ID
	})

	m.entries = entries
	return nil
}

// loadEntry reads and validates a single history entry from disk.
func (m *Manager) loadEntry(_ context.Context, jsonFileName string) (Entry, error) {
	jsonPath := filepath.Join(m.directory, jsonFileName)

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to read json file: %w", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, fmt.Errorf("failed to parse json: %w", err)
	}

	if entry.Version > currentSchemaVersion {
		return Entry{}, fmt.Errorf("unsupported schema version: %d", entry.Version)
	}

	// Apply migrations for older versions if needed
	entry = m.migrateEntry(entry)

	// Verify that the corresponding audio file exists
	wavPath := filepath.Join(m.directory, entry.ID+wavExtension)
	if _, err := os.Stat(wavPath); os.IsNotExist(err) {
		return Entry{}, fmt.Errorf("audio file not found: %s", wavPath)
	}

	return entry, nil
}

// migrateEntry applies necessary transformations for entries with older schema versions.
func (m *Manager) migrateEntry(entry Entry) Entry {
	// Currently at version 1, no migrations needed yet.
	// Future migrations would be added here:
	// if entry.Version < 2 { ... }
	return entry
}

// Write creates a new history entry with atomic write semantics.
// It generates a UUIDv7 ID, saves the audio file first, then the JSON metadata.
// After successful write, it performs pruning if necessary.
func (m *Manager) Write(ctx context.Context, req WriteRequest) (Entry, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Entry{}, fmt.Errorf("failed to generate UUIDv7: %w", err)
	}

	entry := Entry{
		Version:             currentSchemaVersion,
		ID:                  id.String(),
		StartedAt:           req.StartedAt,
		FinishedAt:          req.FinishedAt,
		RecordingDurationMs: req.RecordingDurationMs,
		TranscriptionRaw:    req.TranscriptionRaw,
		TranscriptionFinal:  req.TranscriptionFinal,
		PostProcessed:       req.PostProcessed,
	}

	// Step 1: Save the audio file
	wavPath := filepath.Join(m.directory, entry.ID+wavExtension)
	if err := os.WriteFile(wavPath, req.AudioData, 0644); err != nil {
		return Entry{}, fmt.Errorf("failed to write audio file: %w", err)
	}

	// Step 2: Save the JSON metadata
	jsonPath := filepath.Join(m.directory, entry.ID+jsonExtension)
	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		// Cleanup audio file on failure
		_ = os.Remove(wavPath)
		return Entry{}, fmt.Errorf("failed to marshal entry: %w", err)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		// Cleanup audio file on failure
		_ = os.Remove(wavPath)
		return Entry{}, fmt.Errorf("failed to write json file: %w", err)
	}

	// Step 3: Update in-memory cache
	m.mu.Lock()
	m.entries = append([]Entry{entry}, m.entries...)
	m.mu.Unlock()

	// Step 4: Prune old entries if over limit
	m.prune(ctx)

	m.logger.Debug(ctx, "history entry written", "id", entry.ID)

	return entry, nil
}

// prune removes the oldest entries if the history exceeds the configured limit.
func (m *Manager) prune(ctx context.Context) {
	limit := m.settingsManager.Get().HistoryLimit
	if limit <= 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.entries) <= limit {
		return
	}

	// Entries are sorted newest first, so entries beyond the limit are the oldest
	toRemove := m.entries[limit:]
	m.entries = m.entries[:limit]

	// Delete files for removed entries
	for _, entry := range toRemove {
		if err := m.deleteFilesLocked(entry.ID); err != nil {
			m.logger.Warn(ctx, "failed to delete pruned entry files",
				"id", entry.ID,
				"err", err,
			)
		} else {
			m.logger.Debug(ctx, "pruned history entry", "id", entry.ID)
		}
	}
}

// Delete removes a history entry by ID from both disk and memory.
func (m *Manager) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find and remove from memory
	found := false
	for i, entry := range m.entries {
		if entry.ID == id {
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("entry not found: %s", id)
	}

	// Delete files from disk
	if err := m.deleteFilesLocked(id); err != nil {
		return fmt.Errorf("failed to delete entry files: %w", err)
	}

	m.logger.Debug(ctx, "history entry deleted", "id", id)

	return nil
}

// deleteFilesLocked removes both JSON and WAV files for an entry.
// Caller must hold m.mu lock.
func (m *Manager) deleteFilesLocked(id string) error {
	jsonPath := filepath.Join(m.directory, id+jsonExtension)
	wavPath := filepath.Join(m.directory, id+wavExtension)

	var errs []error

	if err := os.Remove(jsonPath); err != nil && !os.IsNotExist(err) {
		errs = append(errs, fmt.Errorf("json: %w", err))
	}

	if err := os.Remove(wavPath); err != nil && !os.IsNotExist(err) {
		errs = append(errs, fmt.Errorf("wav: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("delete errors: %v", errs)
	}

	return nil
}

// GetAll returns a copy of all history entries.
// This method also refreshes the cache from disk to ensure fresh data.
func (m *Manager) GetAll(ctx context.Context) []Entry {
	m.mu.Lock()
	_ = m.loadLocked(ctx)
	m.mu.Unlock()

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Entry, len(m.entries))
	copy(result, m.entries)
	return result
}

// GetByID returns a specific history entry by ID.
func (m *Manager) GetByID(id string) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, entry := range m.entries {
		if entry.ID == id {
			return entry, true
		}
	}

	return Entry{}, false
}

// GetAudioPath returns the full path to the audio file for an entry.
func (m *Manager) GetAudioPath(id string) string {
	return filepath.Join(m.directory, id+wavExtension)
}

// Count returns the number of entries in the history.
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}

// Clear removes all history entries from disk and memory.
func (m *Manager) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, entry := range m.entries {
		if err := m.deleteFilesLocked(entry.ID); err != nil {
			m.logger.Warn(ctx, "failed to delete entry during clear",
				"id", entry.ID,
				"err", err,
			)
		}
	}

	m.entries = make([]Entry, 0)
	m.logger.Info(ctx, "history cleared")

	return nil
}
