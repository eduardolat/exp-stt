package toggle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/varavelio/tribar/internal/instance"
)

// Execute sends a toggle recording request to the running Tribar instance.
// Must be called after config.EnsureDirectories.
func Execute() error {
	port, err := instance.ReadServerPort()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://127.0.0.1:%d/api/v1/urpc/RecordingToggle", port)
	body := bytes.NewBufferString(`{}`)

	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return fmt.Errorf("failed to connect to Tribar: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			if errMsg, ok := result["error"].(string); ok {
				return fmt.Errorf("toggle failed: %s", errMsg)
			}
		}
		return fmt.Errorf("toggle failed with status: %d", resp.StatusCode)
	}

	return nil
}
