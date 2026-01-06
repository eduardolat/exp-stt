package postprocess

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

// MockSettingsProvider implements SettingsProvider for testing
type MockSettingsProvider struct {
	Settings config.Settings
}

func (m *MockSettingsProvider) Get() config.Settings {
	return m.Settings
}

// MockHTTPClient implements HTTPClient for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}

// MockLogger (minimal implementation or use logger.NewSlogLogger(false))
// Since we are not asserting log output, we can use a discard logger or the real one with false debug.
func newTestLogger() logger.Logger {
	return logger.NewSlogLogger(false)
}

func TestProcess(t *testing.T) {
	ctx := context.Background()

	defaultSettings := config.Settings{
		PostProcessEnabled:  true,
		PostProcessAPIKey:   "test-key",
		PostProcessBaseURL:  "https://api.example.com",
		PostProcessModel:    "test-model",
		PostProcessPromptID: "prompt-1",
		Prompts: []config.Prompt{
			{
				ID:   "prompt-1",
				Body: "Fix this: ${output}",
			},
		},
	}

	t.Run("Happy Path", func(t *testing.T) {
		settingsProvider := &MockSettingsProvider{Settings: defaultSettings}

		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request
				require.Equal(t, "POST", req.Method)
				require.Equal(t, "https://api.example.com/chat/completions", req.URL.String())
				require.Equal(t, "Bearer test-key", req.Header.Get("Authorization"))

				// Return success response
				respBody := `{
					"choices": [
						{
							"message": {
								"content": "Fixed text"
							}
						}
					]
				}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(respBody)),
				}, nil
			},
		}

		instance := New(newTestLogger(), settingsProvider)
		instance.SetClient(mockClient)

		result, err := instance.Process(ctx, "raw text")
		require.NoError(t, err)
		require.Equal(t, "Fixed text", result)
	})

	t.Run("Disabled via Settings", func(t *testing.T) {
		disabledSettings := defaultSettings
		disabledSettings.PostProcessEnabled = false
		settingsProvider := &MockSettingsProvider{Settings: disabledSettings}

		instance := New(newTestLogger(), settingsProvider)
		// No client needed, should not be called

		result, err := instance.Process(ctx, "raw text")
		require.NoError(t, err)
		require.Equal(t, "raw text", result)
	})

	t.Run("Missing API Key", func(t *testing.T) {
		noKeySettings := defaultSettings
		noKeySettings.PostProcessAPIKey = ""
		settingsProvider := &MockSettingsProvider{Settings: noKeySettings}

		instance := New(newTestLogger(), settingsProvider)

		result, err := instance.Process(ctx, "raw text")
		require.NoError(t, err)
		require.Equal(t, "raw text", result)
	})

	t.Run("Empty Input", func(t *testing.T) {
		settingsProvider := &MockSettingsProvider{Settings: defaultSettings}
		instance := New(newTestLogger(), settingsProvider)

		result, err := instance.Process(ctx, "   ")
		require.NoError(t, err)
		require.Equal(t, "   ", result)
	})

	t.Run("Empty Prompt Body", func(t *testing.T) {
		emptyPromptSettings := defaultSettings
		emptyPromptSettings.Prompts = []config.Prompt{
			{ID: "prompt-1", Body: ""}, // Empty body
		}
		settingsProvider := &MockSettingsProvider{Settings: emptyPromptSettings}
		instance := New(newTestLogger(), settingsProvider)

		result, err := instance.Process(ctx, "raw text")
		require.NoError(t, err)
		require.Equal(t, "raw text", result)
	})

	t.Run("Prompt ID Not Found", func(t *testing.T) {
		settingsProvider := &MockSettingsProvider{Settings: defaultSettings}
		settingsProvider.Settings.PostProcessPromptID = "non-existent"
		instance := New(newTestLogger(), settingsProvider)

		result, err := instance.Process(ctx, "raw text")
		require.NoError(t, err)
		require.Equal(t, "raw text", result)
	})

	t.Run("API Error (Network)", func(t *testing.T) {
		settingsProvider := &MockSettingsProvider{Settings: defaultSettings}
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}

		instance := New(newTestLogger(), settingsProvider)
		instance.SetClient(mockClient)

		result, err := instance.Process(ctx, "raw text")
		require.Error(t, err)
		require.Contains(t, err.Error(), "API request failed")
		require.Equal(t, "raw text", result) // Returns original text on error? Wait, implementation returns error.
		// Wait, look at implementation: return text, fmt.Errorf(...)
		// So it returns original text AND error.
	})

	t.Run("API Error (JSON Error Response)", func(t *testing.T) {
		settingsProvider := &MockSettingsProvider{Settings: defaultSettings}
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				respBody := `{"error": {"message": "Rate limit exceeded"}}`
				return &http.Response{
					StatusCode: 429,
					Body:       io.NopCloser(bytes.NewBufferString(respBody)),
				}, nil
			},
		}

		instance := New(newTestLogger(), settingsProvider)
		instance.SetClient(mockClient)

		result, err := instance.Process(ctx, "raw text")
		require.Error(t, err)
		require.Contains(t, err.Error(), "API error: Rate limit exceeded")
		require.Equal(t, "raw text", result)
	})

	t.Run("API Invalid JSON Response", func(t *testing.T) {
		settingsProvider := &MockSettingsProvider{Settings: defaultSettings}
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
				}, nil
			},
		}

		instance := New(newTestLogger(), settingsProvider)
		instance.SetClient(mockClient)

		result, err := instance.Process(ctx, "raw text")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse response")
		require.Equal(t, "raw text", result)
	})

	t.Run("API Empty Choices", func(t *testing.T) {
		settingsProvider := &MockSettingsProvider{Settings: defaultSettings}
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"choices": []}`)),
				}, nil
			},
		}

		instance := New(newTestLogger(), settingsProvider)
		instance.SetClient(mockClient)

		result, err := instance.Process(ctx, "raw text")
		require.Error(t, err)
		require.Contains(t, err.Error(), "no response from API")
		require.Equal(t, "raw text", result)
	})

	t.Run("Empty Content Response", func(t *testing.T) {
		// If API returns empty content string, it should return original text (nil error)
		settingsProvider := &MockSettingsProvider{Settings: defaultSettings}
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"choices": [{"message": {"content": "  "}}]}`)),
				}, nil
			},
		}

		instance := New(newTestLogger(), settingsProvider)
		instance.SetClient(mockClient)

		result, err := instance.Process(ctx, "raw text")
		require.NoError(t, err)
		require.Equal(t, "raw text", result)
	})
}
