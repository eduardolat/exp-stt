package nodes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPNode makes HTTP requests.
type HTTPNode struct{}

// NewHTTPNode creates a new HTTP request node.
func NewHTTPNode() *HTTPNode {
	return &HTTPNode{}
}

// Type returns the node type identifier.
func (n *HTTPNode) Type() string {
	return "http"
}

// Execute makes an HTTP request and returns the response.
func (n *HTTPNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	url, _ := input.Config["url"].(string)
	if url == "" {
		return EmptyOutput(), fmt.Errorf("no URL specified")
	}

	method, _ := input.Config["method"].(string)
	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	// Get timeout (default 30 seconds)
	timeoutMs := 30000
	if v, ok := input.Config["timeoutMs"].(float64); ok {
		timeoutMs = int(v)
	}

	// Get body
	var body io.Reader
	if bodyStr, ok := input.Config["body"].(string); ok && bodyStr != "" {
		body = strings.NewReader(bodyStr)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return EmptyOutput(), fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if headers, ok := input.Config["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strVal, ok := value.(string); ok {
				req.Header.Set(key, strVal)
			}
		}
	}

	// Set default content-type for POST/PUT
	if (method == "POST" || method == "PUT") && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make request
	client := &http.Client{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}

	resp, err := client.Do(req)
	if err != nil {
		return EmptyOutput(), fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return EmptyOutput(), fmt.Errorf("failed to read response: %w", err)
	}

	// Convert response headers to map
	responseHeaders := make(map[string]interface{})
	for key, values := range resp.Header {
		if len(values) == 1 {
			responseHeaders[key] = values[0]
		} else {
			responseHeaders[key] = values
		}
	}

	return NewNodeOutput(map[string]interface{}{
		"response":   string(respBody),
		"statusCode": resp.StatusCode,
		"headers":    responseHeaders,
	}), nil
}
