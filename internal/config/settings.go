package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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

// Settings holds all user-configurable preferences.
type Settings struct {
	Version int `json:"version"`

	// Notification settings
	NotifyOnError  bool `json:"notify_on_error"`
	NotifyOnStart  bool `json:"notify_on_start"`
	NotifyOnFinish bool `json:"notify_on_finish"`

	// Sound settings
	SoundOnStart  bool `json:"sound_on_start"`
	SoundOnFinish bool `json:"sound_on_finish"`

	// Output settings
	OutputMode OutputMode `json:"output_mode"`

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
}

// defaultPrompts returns the predefined prompts for post-processing.
var defaultPrompts = []Prompt{
	{
		ID:   "bc3eb08b-be67-4055-9e3f-40a43a6cc142",
		Name: "Cleanup",
		Body: `You are a transcription cleanup assistant. Your task is to clean up speech-to-text output.

Instructions:
- Fix obvious transcription errors and typos
- Add proper punctuation and capitalization
- Remove filler words (um, uh, like, you know) unless they add meaning
- Preserve the original meaning and tone exactly
- Do not add, remove, or change any substantive content
- Do not add explanations or commentary
- Return only the cleaned text

Output the cleaned transcription directly without any prefix or explanation.

Raw transcription: ${output}`,
	},
	{
		ID:   "55f21e4e-b314-4a49-9a61-8d22ddc8713b",
		Name: "Formal",
		Body: `You are a professional writing assistant. Your task is to transform speech-to-text output into formal written text.

Instructions:
- Convert spoken language to formal written style
- Fix grammar, punctuation, and capitalization
- Remove filler words and verbal pauses
- Use professional vocabulary where appropriate
- Maintain the original meaning and intent
- Structure sentences for clarity
- Do not add new information or change the meaning
- Return only the formatted text

Output the formal text directly without any prefix or explanation.

Raw transcription: ${output}`,
	},
	{
		ID:   "3eeed235-49e3-487c-84a2-85028e15ca93",
		Name: "Technical",
		Body: `You are a technical documentation assistant. Your task is to format speech-to-text output for technical contexts.

Instructions:
- Fix transcription errors, especially technical terms
- Use proper technical terminology and formatting
- Add appropriate punctuation for code-related content
- Format lists and steps clearly if present
- Preserve technical accuracy
- Use clear, concise language
- Do not add explanations or change meaning
- Return only the formatted text

Output the technical text directly without any prefix or explanation.

Raw transcription: ${output}`,
	},
	{
		ID:   "1dee2faf-4c65-4349-a48b-16768791efe9",
		Name: "Creative",
		Body: `You are a creative writing assistant. Your task is to polish speech-to-text output while preserving creative voice.

Instructions:
- Fix transcription errors and add punctuation
- Preserve the speaker's unique voice and style
- Keep creative expressions and metaphors
- Maintain emotional tone and emphasis
- Remove only meaningless filler words
- Do not change the creative intent
- Return only the polished text

Output the polished text directly without any prefix or explanation.

Raw transcription: ${output}`,
	},
}

// defaultSettings returns the default application settings.
var defaultSettings = Settings{
	Version: 1,

	NotifyOnError:  true,
	NotifyOnStart:  false,
	NotifyOnFinish: false,

	SoundOnStart:  true,
	SoundOnFinish: true,

	OutputMode: OutputModeCopyPaste,

	PostProcessEnabled:  false,
	PostProcessBaseURL:  "https://api.openai.com/v1",
	PostProcessAPIKey:   "",
	PostProcessModel:    "gpt-4o-mini",
	PostProcessPromptID: defaultPrompts[0].ID,

	Prompts: defaultPrompts,

	HistoryLimit: 10,
}

// SettingsManager handles loading and saving of user settings.
type SettingsManager struct {
	mu       sync.RWMutex
	settings Settings
	filePath string
}

// NewSettingsManager creates a new settings manager and loads existing settings.
func NewSettingsManager() (*SettingsManager, error) {
	sm := &SettingsManager{
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
	return sm.saveUnsafe()
}

// Load reads settings from the config file.
func (sm *SettingsManager) Load() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse settings: %w", err)
	}

	sm.settings = settings
	return nil
}

// Save writes the current settings to the config file.
func (sm *SettingsManager) Save() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.saveUnsafe()
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
