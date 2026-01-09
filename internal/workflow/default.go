package workflow

import (
	_ "embed"
	"encoding/json"
	"sync"
)

//go:embed default.json
var defaultWorkflowJSON []byte

var (
	defaultWorkflow     *Workflow
	defaultWorkflowOnce sync.Once
)

// DefaultWorkflow returns the hardcoded default workflow.
// This is used when AdvancedMode is disabled or no workflow is selected.
// The workflow is parsed once and cached for subsequent calls.
func DefaultWorkflow() *Workflow {
	defaultWorkflowOnce.Do(func() {
		var wf Workflow
		if err := json.Unmarshal(defaultWorkflowJSON, &wf); err != nil {
			panic("failed to parse embedded default workflow: " + err.Error())
		}
		defaultWorkflow = &wf
	})
	return defaultWorkflow.Clone()
}

// IsDefaultWorkflow returns true if the given workflow ID is the default workflow.
func IsDefaultWorkflow(id string) bool {
	return id == "default" || id == ""
}
