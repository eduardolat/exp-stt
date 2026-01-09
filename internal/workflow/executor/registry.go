package executor

import (
	"sync"

	"github.com/varavelio/tribar/internal/workflow"
	"github.com/varavelio/tribar/internal/workflow/executor/nodes"
)

// NodeRegistry manages the mapping of node types to their executors.
type NodeRegistry struct {
	mu        sync.RWMutex
	executors map[workflow.NodeType]nodes.NodeExecutor
}

// NewNodeRegistry creates a new node registry.
func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		executors: make(map[workflow.NodeType]nodes.NodeExecutor),
	}
}

// Register adds a node executor to the registry.
func (r *NodeRegistry) Register(executor nodes.NodeExecutor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.executors[workflow.NodeType(executor.Type())] = executor
}

// Get retrieves a node executor by type.
func (r *NodeRegistry) Get(nodeType workflow.NodeType) (nodes.NodeExecutor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	exec, ok := r.executors[nodeType]
	return exec, ok
}

// RegisteredTypes returns all registered node types.
func (r *NodeRegistry) RegisteredTypes() []workflow.NodeType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	types := make([]workflow.NodeType, 0, len(r.executors))
	for t := range r.executors {
		types = append(types, t)
	}
	return types
}

// DefaultRegistry creates a registry with all built-in node executors.
func DefaultRegistry() *NodeRegistry {
	r := NewNodeRegistry()

	// Register all built-in nodes
	r.Register(nodes.NewTranscribeNode())
	r.Register(nodes.NewAINode())
	r.Register(nodes.NewClipboardCopyNode())
	r.Register(nodes.NewClipboardPasteNode())
	r.Register(nodes.NewJavaScriptNode())
	r.Register(nodes.NewTerminalNode())
	r.Register(nodes.NewNotifyNode())
	r.Register(nodes.NewSoundNode())
	r.Register(nodes.NewConditionNode())
	r.Register(nodes.NewDelayNode())
	r.Register(nodes.NewHTTPNode())

	return r
}
