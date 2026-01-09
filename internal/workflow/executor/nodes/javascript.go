package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dop251/goja"
)

const jsExecutionTimeout = 5 * time.Second

// JavaScriptNode executes JavaScript code using Goja.
type JavaScriptNode struct{}

// NewJavaScriptNode creates a new JavaScript execution node.
func NewJavaScriptNode() *JavaScriptNode {
	return &JavaScriptNode{}
}

// Type returns the node type identifier.
func (n *JavaScriptNode) Type() string {
	return "javascript"
}

// Execute runs the JavaScript code with sandboxed environment.
func (n *JavaScriptNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	script, _ := input.Config["script"].(string)
	if script == "" {
		return EmptyOutput(), fmt.Errorf("no script provided")
	}

	vm := goja.New()

	// Set up timeout interrupt
	timer := time.AfterFunc(jsExecutionTimeout, func() {
		vm.Interrupt("execution timeout exceeded (5 seconds)")
	})
	defer timer.Stop()

	// Expose input data to the script
	inputData := input.Config["input"]
	if inputData == nil {
		inputData = ""
	}
	if err := vm.Set("input", inputData); err != nil {
		return EmptyOutput(), fmt.Errorf("failed to set input: %w", err)
	}

	// Expose a simple log function
	logs := make([]string, 0)
	logFunc := func(call goja.FunctionCall) goja.Value {
		for _, arg := range call.Arguments {
			logs = append(logs, fmt.Sprintf("%v", arg.Export()))
		}
		return goja.Undefined()
	}
	if err := vm.Set("log", logFunc); err != nil {
		return EmptyOutput(), fmt.Errorf("failed to set log function: %w", err)
	}

	// Expose console.log
	console := vm.NewObject()
	if err := console.Set("log", logFunc); err != nil {
		return EmptyOutput(), fmt.Errorf("failed to set console.log: %w", err)
	}
	if err := vm.Set("console", console); err != nil {
		return EmptyOutput(), fmt.Errorf("failed to set console: %w", err)
	}

	// Expose JSON utilities
	if err := vm.Set("JSON", map[string]interface{}{
		"parse": func(s string) (interface{}, error) {
			var result interface{}
			if err := json.Unmarshal([]byte(s), &result); err != nil {
				return nil, err
			}
			return result, nil
		},
		"stringify": func(v interface{}) (string, error) {
			data, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(data), nil
		},
	}); err != nil {
		return EmptyOutput(), fmt.Errorf("failed to set JSON: %w", err)
	}

	// Run the script
	result, err := vm.RunString(script)
	if err != nil {
		return EmptyOutput(), fmt.Errorf("javascript execution failed: %w", err)
	}

	// Export the result
	var output interface{}
	if result != nil && !goja.IsUndefined(result) && !goja.IsNull(result) {
		output = result.Export()
	}

	// Convert output to string if needed
	outputStr := ""
	if output != nil {
		if s, ok := output.(string); ok {
			outputStr = s
		} else {
			data, _ := json.Marshal(output)
			outputStr = string(data)
		}
	}

	return NewNodeOutput(map[string]interface{}{
		"output": outputStr,
		"logs":   logs,
	}), nil
}
