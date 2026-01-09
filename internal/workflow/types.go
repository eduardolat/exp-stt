// Package workflow provides the core types and management for Tribar's workflow system.
package workflow

import (
	"encoding/json"
	"time"
)

// Workflow represents a complete automation workflow with nodes and edges.
type Workflow struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Trigger     Trigger   `json:"trigger"`
	Nodes       []Node    `json:"nodes"`
	Edges       []Edge    `json:"edges"`
}

// Trigger defines what activates a workflow.
type Trigger struct {
	Type   TriggerType            `json:"type"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// TriggerType represents the type of trigger that activates a workflow.
type TriggerType string

const (
	TriggerTypeVoice TriggerType = "voice"
)

// Node represents a single step in a workflow.
type Node struct {
	ID       string                 `json:"id"`
	Type     NodeType               `json:"type"`
	Position Position               `json:"position"`
	Config   map[string]interface{} `json:"config,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// NodeType represents the type of a workflow node.
type NodeType string

const (
	NodeTypeTranscribe     NodeType = "transcribe"
	NodeTypeAIProcess      NodeType = "ai_process"
	NodeTypeClipboardCopy  NodeType = "clipboard_copy"
	NodeTypeClipboardPaste NodeType = "clipboard_paste"
	NodeTypeJavaScript     NodeType = "javascript"
	NodeTypeTerminal       NodeType = "terminal"
	NodeTypeNotify         NodeType = "notify"
	NodeTypeSound          NodeType = "sound"
	NodeTypeCondition      NodeType = "condition"
	NodeTypeDelay          NodeType = "delay"
	NodeTypeHTTP           NodeType = "http"
)

// Position represents the visual position of a node in the editor.
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Edge represents a connection between two nodes.
type Edge struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	Target       string `json:"target"`
	SourceHandle string `json:"sourceHandle,omitempty"`
	TargetHandle string `json:"targetHandle,omitempty"`
}

// NodeDefinition describes a node type for the frontend palette.
type NodeDefinition struct {
	Type        NodeType               `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Inputs      []string               `json:"inputs"`
	Outputs     map[string]OutputField `json:"outputs"`
}

// OutputField describes an output field of a node.
type OutputField struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// WorkflowInput is used for creating/updating workflows (without readonly fields).
type WorkflowInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Enabled     bool    `json:"enabled"`
	Trigger     Trigger `json:"trigger"`
	Nodes       []Node  `json:"nodes"`
	Edges       []Edge  `json:"edges"`
}

// Clone creates a deep copy of the workflow.
func (w *Workflow) Clone() *Workflow {
	data, _ := json.Marshal(w)
	var clone Workflow
	_ = json.Unmarshal(data, &clone)
	return &clone
}

// GetAvailableNodeDefinitions returns all available node types for the frontend.
func GetAvailableNodeDefinitions() []NodeDefinition {
	return []NodeDefinition{
		{
			Type:        NodeTypeTranscribe,
			Name:        "Transcribe Audio",
			Description: "Convert audio to text using Parakeet/ONNX",
			Category:    "core",
			Inputs:      []string{"audioData"},
			Outputs: map[string]OutputField{
				"rawText": {Type: "string", Description: "Raw transcription text"},
			},
		},
		{
			Type:        NodeTypeAIProcess,
			Name:        "AI Processing",
			Description: "Process text with LLM for enhancement",
			Category:    "ai",
			Inputs:      []string{"text", "prompt"},
			Outputs: map[string]OutputField{
				"processedText": {Type: "string", Description: "AI-enhanced text"},
				"tokensUsed":    {Type: "number", Description: "API tokens consumed"},
			},
		},
		{
			Type:        NodeTypeClipboardCopy,
			Name:        "Clipboard Copy",
			Description: "Copy text to clipboard",
			Category:    "output",
			Inputs:      []string{"text"},
			Outputs: map[string]OutputField{
				"text": {Type: "string", Description: "Text passed through"},
			},
		},
		{
			Type:        NodeTypeClipboardPaste,
			Name:        "Clipboard Paste",
			Description: "Paste text using keyboard shortcut",
			Category:    "output",
			Inputs:      []string{"text"},
			Outputs: map[string]OutputField{
				"text": {Type: "string", Description: "Text passed through"},
			},
		},
		{
			Type:        NodeTypeJavaScript,
			Name:        "JavaScript",
			Description: "Execute custom JavaScript code",
			Category:    "transform",
			Inputs:      []string{"input"},
			Outputs: map[string]OutputField{
				"output": {Type: "any", Description: "Script output"},
			},
		},
		{
			Type:        NodeTypeTerminal,
			Name:        "Terminal Command",
			Description: "Execute a shell command",
			Category:    "system",
			Inputs:      []string{"command", "args"},
			Outputs: map[string]OutputField{
				"stdout":   {Type: "string", Description: "Standard output"},
				"stderr":   {Type: "string", Description: "Standard error"},
				"exitCode": {Type: "number", Description: "Exit code"},
			},
		},
		{
			Type:        NodeTypeNotify,
			Name:        "Notification",
			Description: "Show a desktop notification",
			Category:    "feedback",
			Inputs:      []string{"title", "message"},
			Outputs:     map[string]OutputField{},
		},
		{
			Type:        NodeTypeSound,
			Name:        "Play Sound",
			Description: "Play an audio feedback sound",
			Category:    "feedback",
			Inputs:      []string{"soundId"},
			Outputs:     map[string]OutputField{},
		},
		{
			Type:        NodeTypeCondition,
			Name:        "Condition",
			Description: "Branch workflow based on a condition",
			Category:    "flow",
			Inputs:      []string{"value", "condition"},
			Outputs: map[string]OutputField{
				"true":  {Type: "any", Description: "Output when condition is true"},
				"false": {Type: "any", Description: "Output when condition is false"},
			},
		},
		{
			Type:        NodeTypeDelay,
			Name:        "Delay",
			Description: "Wait for a specified duration",
			Category:    "flow",
			Inputs:      []string{"input", "delayMs"},
			Outputs: map[string]OutputField{
				"output": {Type: "any", Description: "Input passed through after delay"},
			},
		},
		{
			Type:        NodeTypeHTTP,
			Name:        "HTTP Request",
			Description: "Make an HTTP request",
			Category:    "network",
			Inputs:      []string{"url", "method", "headers", "body"},
			Outputs: map[string]OutputField{
				"response":   {Type: "string", Description: "Response body"},
				"statusCode": {Type: "number", Description: "HTTP status code"},
				"headers":    {Type: "object", Description: "Response headers"},
			},
		},
	}
}
