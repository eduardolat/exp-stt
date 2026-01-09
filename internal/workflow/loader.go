package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Loader handles reading and writing workflows to the filesystem.
type Loader struct {
	workflowsDir string
}

// NewLoader creates a new workflow loader.
func NewLoader(workflowsDir string) *Loader {
	return &Loader{workflowsDir: workflowsDir}
}

// EnsureDir creates the workflows directory if it doesn't exist.
func (l *Loader) EnsureDir() error {
	return os.MkdirAll(l.workflowsDir, 0755)
}

// workflowPath returns the file path for a workflow by ID.
func (l *Loader) workflowPath(id string) string {
	return filepath.Join(l.workflowsDir, id+".json")
}

// Load reads a workflow from disk by ID.
func (l *Loader) Load(id string) (*Workflow, error) {
	data, err := os.ReadFile(l.workflowPath(id))
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var wf Workflow
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}

	return &wf, nil
}

// Save writes a workflow to disk.
func (l *Loader) Save(wf *Workflow) error {
	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	if err := os.WriteFile(l.workflowPath(wf.ID), data, 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	return nil
}

// Delete removes a workflow file from disk.
func (l *Loader) Delete(id string) error {
	path := l.workflowPath(id)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete workflow file: %w", err)
	}
	return nil
}

// LoadAll reads all workflows from the workflows directory.
func (l *Loader) LoadAll() ([]*Workflow, error) {
	entries, err := os.ReadDir(l.workflowsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read workflows directory: %w", err)
	}

	var workflows []*Workflow
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := entry.Name()[:len(entry.Name())-5] // Remove .json extension
		wf, err := l.Load(id)
		if err != nil {
			continue // Skip invalid files
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil
}

// Exists checks if a workflow file exists.
func (l *Loader) Exists(id string) bool {
	_, err := os.Stat(l.workflowPath(id))
	return err == nil
}
