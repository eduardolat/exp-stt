package nodes

import (
	"context"

	"github.com/varavelio/tribar/internal/clipboard"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/history"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/notify"
	"github.com/varavelio/tribar/internal/postprocess"
	"github.com/varavelio/tribar/internal/sound"
	"github.com/varavelio/tribar/internal/transcribe"
)

// NodeInput contains the data passed to a node during execution.
type NodeInput struct {
	// Config is the resolved node configuration.
	Config map[string]interface{}
	// Data holds additional node-specific data.
	Data map[string]interface{}
}

// NodeExecutor is the interface all workflow nodes must implement.
type NodeExecutor interface {
	// Execute runs the node with the given input and services.
	// Nodes receive services via the ServiceProvider - they never instantiate services.
	Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error)

	// Type returns the node type identifier (e.g., "transcribe", "ai_process").
	Type() string
}

// ServiceProvider provides access to shared services for node execution.
type ServiceProvider interface {
	GetLogger() logger.Logger
	GetSettingsManager() *config.SettingsManager
	GetTranscriber() *transcribe.Instance
	GetPostProcessor() *postprocess.Instance
	GetNotifier() *notify.Instance
	GetSound() *sound.Instance
	GetClipboard() *clipboard.Instance
	GetHistoryManager() *history.Manager
}
