package postprocess

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/varavelio/tribar/internal/config"
)

// MockSettingsProvider implements SettingsProvider for testing.
type MockSettingsProvider struct {
	settings config.Settings
}

func (m *MockSettingsProvider) Get() config.Settings {
	return m.settings
}

// MockRoundTripper implements http.RoundTripper for testing.
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

// MockLogger implements logger.Logger for testing (no-op).
type MockLogger struct{}

func (m *MockLogger) SetDebug(enabled bool)                             {}
func (m *MockLogger) Debug(ctx context.Context, msg string, args ...any) {}
func (m *MockLogger) Info(ctx context.Context, msg string, args ...any)  {}
func (m *MockLogger) Warn(ctx context.Context, msg string, args ...any)  {}
func (m *MockLogger) Error(ctx context.Context, msg string, args ...any) {}

func TestProcess(t *testing.T) {
	// Setup default test settings
	defaultPrompt := config.Prompt{
		ID:   "test-prompt-id",
		Body: "Fix this: ${output}",
	}

	baseSettings := config.Settings{
		PostProcessEnabled:  true,
		PostProcessAPIKey:   "sk-test-key",
		PostProcessBaseURL:  "https://api.test.com/v1",
		PostProcessModel:    "gpt-test-model",
		PostProcessPromptID: "test-prompt-id",
		Prompts:             []config.Prompt{defaultPrompt},
	}

	tests := []struct {
		name          string
		settings      config.Settings
		inputText     string
		mockResponse  *http.Response
		mockError     error
		expectedText  string
		expectedError string
	}{
		{
			name: "Happy Path",
			settings: baseSettings,
			inputText: "hello world",
			mockResponse: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"choices": [{
						"message": {
							"content": "Hello, World!"
						}
					}]
				}`)),
			},
			expectedText: "Hello, World!",
		},
		{
			name: "Disabled via Settings",
			settings: func() config.Settings {
				s := baseSettings
				s.PostProcessEnabled = false
				return s
			}(),
			inputText:    "hello",
			expectedText: "hello",
		},
		{
			name: "Missing API Key",
			settings: func() config.Settings {
				s := baseSettings
				s.PostProcessAPIKey = ""
				return s
			}(),
			inputText:    "hello",
			expectedText: "hello",
		},
		{
			name:         "Empty Input",
			settings:     baseSettings,
			inputText:    "   ",
			expectedText: "   ",
		},
		{
			name: "Missing Prompt Configuration",
			settings: func() config.Settings {
				s := baseSettings
				s.PostProcessPromptID = "non-existent-id"
				return s
			}(),
			inputText:    "hello",
			expectedText: "hello",
		},
		{
			name:      "API Network Error",
			settings:  baseSettings,
			inputText: "hello",
			mockError: errors.New("connection refused"),
			expectedText: "hello", // Returns original text on error
			expectedError: "API request failed",
		},
		{
			name:      "API JSON Error Response",
			settings:  baseSettings,
			inputText: "hello",
			mockResponse: &http.Response{
				StatusCode: 400,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"error": {
						"message": "Invalid model"
					}
				}`)),
			},
			expectedText: "hello",
			expectedError: "API error: Invalid model",
		},
		{
			name:      "API Empty Response Choices",
			settings:  baseSettings,
			inputText: "hello",
			mockResponse: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"choices": []
				}`)),
			},
			expectedText: "hello",
			expectedError: "no response from API",
		},
		{
			name:      "API Malformed JSON",
			settings:  baseSettings,
			inputText: "hello",
			mockResponse: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{ invalid json }`)),
			},
			expectedText: "hello",
			expectedError: "failed to parse response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settingsProvider := &MockSettingsProvider{settings: tt.settings}
			mockLogger := &MockLogger{}
			instance := New(mockLogger, settingsProvider)

			if tt.mockResponse != nil || tt.mockError != nil {
				mockTransport := &MockRoundTripper{
					RoundTripFunc: func(req *http.Request) (*http.Response, error) {
						// Verify request headers
						require.Equal(t, "application/json", req.Header.Get("Content-Type"))
						require.Equal(t, "Bearer "+tt.settings.PostProcessAPIKey, req.Header.Get("Authorization"))

						// Verify request body
						body, _ := io.ReadAll(req.Body)
						var reqBody chatRequest
						_ = json.Unmarshal(body, &reqBody)

						require.Equal(t, tt.settings.PostProcessModel, reqBody.Model)
						require.NotEmpty(t, reqBody.Messages)

						// Verify content (simple check as it includes prompt)
						require.Contains(t, reqBody.Messages[0].Content, tt.inputText)

						if tt.mockError != nil {
							return nil, tt.mockError
						}
						return tt.mockResponse, nil
					},
				}
				instance.SetHTTPClient(&http.Client{Transport: mockTransport})
			}

			result, err := instance.Process(context.Background(), tt.inputText)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedText, result)
		})
	}
}
