package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/varavelio/tribar/internal/eventbus"
)

const settingsFileName = "settings.json"

// OutputMode defines how transcription results are delivered to the user.
type OutputMode string

const (
	OutputModeCopyOnly   OutputMode = "copy_only"
	OutputModeCopyPaste  OutputMode = "copy_paste"
	OutputModeGhostPaste OutputMode = "ghost_paste"
)

// Prompt represents a user-configurable prompt for post-processing.
type Prompt struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}

// Shortcut represents a global keyboard shortcut configuration.
type Shortcut struct {
	Modifiers []string `json:"modifiers"` // "ctrl", "alt", "shift", "meta"
	Key       string   `json:"key"`       // "space", "a", etc.
}

// Settings holds all user-configurable preferences.
type Settings struct {
	Version int `json:"version"`

	// Notification settings
	NotifyOnError  bool `json:"notify_on_error"`
	NotifyOnStart  bool `json:"notify_on_start"`
	NotifyOnFinish bool `json:"notify_on_finish"`

	// Sound feedback settings
	SoundFeedbackEnable    bool   `json:"sound_feedback_enable"`
	SoundFeedbackRecordID  string `json:"sound_feedback_record_id"`
	SoundFeedbackSuccessID string `json:"sound_feedback_success_id"`
	SoundFeedbackErrorID   string `json:"sound_feedback_error_id"`
	SoundFeedbackVolume    int    `json:"sound_feedback_volume"`

	// Audio device settings
	InputDevice string `json:"input_device"`

	// Output settings
	OutputMode          OutputMode `json:"output_mode"`
	OutputTrailingSpace bool       `json:"output_trailing_space"`

	// Post-processing settings
	PostProcessEnabled  bool   `json:"postprocess_enabled"`
	PostProcessBaseURL  string `json:"postprocess_base_url"`
	PostProcessAPIKey   string `json:"postprocess_api_key"`
	PostProcessModel    string `json:"postprocess_model"`
	PostProcessPromptID string `json:"postprocess_prompt_id"`

	// Prompts for post-processing
	Prompts []Prompt `json:"prompts"`

	// History settings
	HistoryLimit int `json:"history_limit"`

	// Model auto-unload settings
	ModelUnloadEnable  bool `json:"model_unload_enable"`
	ModelUnloadSeconds int  `json:"model_unload_seconds"`

	// Security settings
	// Specific origins allowed to access the API via CORS.
	// If any element is "*", all origins are allowed.
	AllowedCORSOrigins []string `json:"allowed_cors_origins"`

	// Global shortcut settings
	ShortcutToggle Shortcut `json:"shortcut_toggle"`

	// Paste shortcut settings
	PasteShortcut string `json:"paste_shortcut"`
}

// defaultPrompts returns the predefined prompts for post-processing.
var defaultPrompts = []Prompt{
	{
		ID:   "bc3eb08b-be67-4055-9e3f-40a43a6cc142",
		Name: "Cleanup transcription",

		Body: strings.Join([]string{
			"You are an expert editor for raw speech-to-text transcriptions. Your task is to format the input into clear, readable written text.",
			"",
			"Strictly follow these rules:",
			"1. Punctuation & Formatting: Add proper punctuation (periods, commas, question marks) and capitalization.",
			"2. Language: Output MUST be in the same language as the input. Do not translate.",
			"3. Cleanup: Remove filler words (like \"um\", \"uh\", \"like\", \"you know\") and stuttering, but keep the core meaning intact.",
			"4. Output: specificy ONLY the processed text. Do not add quotes, prefixes (like \"Here is the text:\"), or explanations. Respond literally with the processed text.",
			"",
			"Raw transcription:",
			"${output}",
		}, "\n"),
	},
}

// defaultSettings returns the default application settings.
var defaultSettings = Settings{
	Version: 1,

	NotifyOnError:  true,
	NotifyOnStart:  false,
	NotifyOnFinish: false,

	SoundFeedbackEnable:    true,
	SoundFeedbackRecordID:  "1",
	SoundFeedbackSuccessID: "1",
	SoundFeedbackErrorID:   "1",
	SoundFeedbackVolume:    100,

	InputDevice: "default",

	OutputMode:          OutputModeCopyPaste,
	OutputTrailingSpace: false,

	PostProcessEnabled:  false,
	PostProcessBaseURL:  "https://api.openai.com/v1",
	PostProcessAPIKey:   "",
	PostProcessModel:    "gpt-4o-mini",
	PostProcessPromptID: defaultPrompts[0].ID,

	Prompts: defaultPrompts,

	HistoryLimit: 10,

	ModelUnloadEnable:  true,
	ModelUnloadSeconds: 300,

	AllowedCORSOrigins: []string{},

	ShortcutToggle: Shortcut{
		Modifiers: []string{"ctrl"},
		Key:       "space",
	},

	PasteShortcut: "ctrl+v",
}

// SettingsManager handles loading and saving of user settings.
type SettingsManager struct {
	eventBus *eventbus.EventBus
	mu       sync.RWMutex
	settings Settings
	filePath string
}

// NewSettingsManager creates a new settings manager and loads existing settings.
func NewSettingsManager(eventBus *eventbus.EventBus) (*SettingsManager, error) {
	sm := &SettingsManager{
		eventBus: eventBus,
		settings: defaultSettings,
		filePath: filepath.Join(DirectoryConfig, settingsFileName),
	}

	if err := sm.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load settings: %w", err)
		}
		// File doesn't exist, save defaults
		if err := sm.Save(); err != nil {
			return nil, fmt.Errorf("failed to save default settings: %w", err)
		}
	}

	return sm, nil
}

// Get returns a copy of the current settings.
func (sm *SettingsManager) Get() Settings {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.settings
}

// Update updates the settings and saves them to disk.
func (sm *SettingsManager) Update(settings Settings) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.settings = settings
	if err := sm.saveUnsafe(); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	sm.eventBus.PublishSettingsChanged()
	return nil
}

// Load reads settings from the config file, merging with defaults.
// This ensures new settings get their default values when upgrading.
func (sm *SettingsManager) Load() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return err
	}

	// Start with defaults, then unmarshal file contents on top
	// This preserves default values for any new fields not in the file
	merged := defaultSettings
	if err := json.Unmarshal(data, &merged); err != nil {
		return fmt.Errorf("failed to parse settings: %w", err)
	}

	sm.settings = merged

	// Re-save to persist any new default fields to the file
	return sm.saveUnsafe()
}

// Save writes the current settings to the config file.
func (sm *SettingsManager) Save() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if err := sm.saveUnsafe(); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}
	return nil
}

// saveUnsafe writes settings to disk without acquiring the lock.
func (sm *SettingsManager) saveUnsafe() error {
	data, err := json.MarshalIndent(sm.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(sm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}
