// Package postprocess provides LLM-based text enhancement for transcriptions.
// It supports OpenAI-compatible APIs to improve grammar, punctuation, and formatting.
package postprocess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

// Instance handles LLM-based text post-processing.
type Instance struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager
	client          *http.Client
}

const defaultTimeout = 30 * time.Second

// New creates a new post-processor instance.
func New(logger logger.Logger, settingsManager *config.SettingsManager) *Instance {
	return &Instance{
		logger:          logger,
		settingsManager: settingsManager,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// IsEnabled returns whether post-processing is enabled.
func (p *Instance) IsEnabled() bool {
	settings := p.settingsManager.Get()
	return settings.PostProcessEnabled && settings.PostProcessAPIKey != ""
}

// Process enhances the transcription using the configured LLM.
func (p *Instance) Process(ctx context.Context, text string) (string, error) {
	if !p.IsEnabled() {
		return text, nil
	}

	if strings.TrimSpace(text) == "" {
		return text, nil
	}

	prompt := p.getSystemPrompt()
	if prompt == "" {
		return text, nil
	}

	// Replace placeholder with actual transcription text
	input := strings.ReplaceAll(prompt, "${output}", text)
	return p.callAPI(ctx, input)
}

// getSystemPrompt returns the prompt body for the configured prompt ID.
func (p *Instance) getSystemPrompt() string {
	settings := p.settingsManager.Get()

	for _, prompt := range settings.Prompts {
		if prompt.ID == settings.PostProcessPromptID {
			return prompt.Body
		}
	}

	return ""
}

// chatRequest represents the OpenAI chat completion request.
type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatResponse represents the OpenAI chat completion response.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// callAPI sends the text to the LLM API for enhancement.
func (p *Instance) callAPI(ctx context.Context, text string) (string, error) {
	settings := p.settingsManager.Get()

	reqBody := chatRequest{
		Model: settings.PostProcessModel,
		Messages: []message{
			{Role: "user", Content: text},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return text, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := strings.TrimSuffix(settings.PostProcessBaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return text, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+settings.PostProcessAPIKey)
	req.Header.Set("User-Agent", "Tribar/"+config.AppVersion)
	req.Header.Set("X-Title", config.AppName)
	req.Header.Set("HTTP-Referer", "https://github.com/varavel/tribar")

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Error(ctx, "LLM API request failed", "err", err)
		return text, fmt.Errorf("API request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return text, fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return text, fmt.Errorf("failed to parse response: %w", err)
	}

	if chatResp.Error != nil {
		return text, fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return text, fmt.Errorf("no response from API")
	}

	result := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	if result == "" {
		return text, nil
	}

	p.logger.Debug(ctx, "text post-processed",
		"original_length", len(text),
		"processed_length", len(result),
	)

	return result, nil
}
