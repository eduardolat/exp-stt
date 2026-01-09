package nodes

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// TerminalNode executes shell commands.
type TerminalNode struct{}

// NewTerminalNode creates a new terminal node.
func NewTerminalNode() *TerminalNode {
	return &TerminalNode{}
}

// Type returns the node type identifier.
func (n *TerminalNode) Type() string {
	return "terminal"
}

// Execute runs a shell command and returns its output.
func (n *TerminalNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	command, _ := input.Config["command"].(string)
	if command == "" {
		return EmptyOutput(), fmt.Errorf("no command specified")
	}

	// Get timeout (default 30 seconds)
	timeoutMs := 30000
	if v, ok := input.Config["timeoutMs"].(float64); ok {
		timeoutMs = int(v)
	}

	// Create context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	// Get arguments
	var args []string
	if argsRaw, ok := input.Config["args"].([]interface{}); ok {
		for _, arg := range argsRaw {
			if s, ok := arg.(string); ok {
				args = append(args, s)
			}
		}
	} else if argsStr, ok := input.Config["args"].(string); ok && argsStr != "" {
		args = strings.Fields(argsStr)
	}

	// Create command
	cmd := exec.CommandContext(cmdCtx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return EmptyOutput(), fmt.Errorf("command execution failed: %w", err)
		}
	}

	return NewNodeOutput(map[string]interface{}{
		"stdout":   stdout.String(),
		"stderr":   stderr.String(),
		"exitCode": exitCode,
	}), nil
}
