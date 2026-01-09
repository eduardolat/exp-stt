package executor

import (
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/varavelio/tribar/internal/workflow/executor/nodes"
)

// ExecutionContext holds the state during workflow execution.
// It stores node outputs and provides variable interpolation.
type ExecutionContext struct {
	ExecutionID string
	TriggerData map[string]interface{}
	outputs     map[string]nodes.NodeOutput
	mu          sync.RWMutex
}

// NewExecutionContext creates a new execution context with trigger data.
func NewExecutionContext(triggerData map[string]interface{}) *ExecutionContext {
	return &ExecutionContext{
		ExecutionID: uuid.Must(uuid.NewV7()).String(),
		TriggerData: triggerData,
		outputs:     make(map[string]nodes.NodeOutput),
	}
}

// SetOutput stores a node's output for later reference by other nodes.
func (c *ExecutionContext) SetOutput(nodeID string, output nodes.NodeOutput) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.outputs[nodeID] = output
}

// GetOutput retrieves a specific node's output.
func (c *ExecutionContext) GetOutput(nodeID string) (nodes.NodeOutput, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out, ok := c.outputs[nodeID]
	return out, ok
}

// variableRegex matches {{nodes.nodeId.path.to.field}} patterns.
// The second group captures everything after the nodeId, supporting nested paths.
var variableRegex = regexp.MustCompile(`\{\{nodes\.([a-zA-Z0-9_-]+)\.(.+?)\}\}`)

// triggerRegex matches {{trigger.fieldName}} patterns.
var triggerRegex = regexp.MustCompile(`\{\{trigger\.([a-zA-Z0-9_]+)\}\}`)

// ResolveVariables interpolates variable references in a string.
// Supported patterns:
// - {{nodes.nodeId.field}} or {{nodes.nodeId.nested.field}} - Reference to another node's output
// - {{trigger.fieldName}} - Reference to trigger data
func (c *ExecutionContext) ResolveVariables(input string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Resolve node references: {{nodes.nodeId.field.path}}
	result := variableRegex.ReplaceAllStringFunc(input, func(match string) string {
		parts := variableRegex.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		nodeID, fieldPath := parts[1], parts[2]
		output, ok := c.outputs[nodeID]
		if !ok {
			return match // Leave unresolved if node not found
		}

		// Handle nested paths like "headers.Content-Type"
		return output.GetFieldPath(strings.Split(fieldPath, "."))
	})

	// Resolve trigger references: {{trigger.fieldName}}
	result = triggerRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := triggerRegex.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}

		field := parts[1]
		if val, ok := c.TriggerData[field]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return match
	})

	return result
}

// ResolveConfig resolves variables in a node config map.
func (c *ExecutionContext) ResolveConfig(config map[string]interface{}) map[string]interface{} {
	resolved := make(map[string]interface{}, len(config))
	for key, value := range config {
		if str, ok := value.(string); ok {
			resolved[key] = c.ResolveVariables(str)
		} else {
			resolved[key] = value
		}
	}
	return resolved
}
