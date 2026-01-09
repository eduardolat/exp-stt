package workflow

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// transcribeNodeID is the fixed ID for the Transcribe entry node.
const transcribeNodeID = "trigger-transcribe"

// ErrWorkflowNotFound is returned when a workflow is not found.
var ErrWorkflowNotFound = errors.New("workflow not found")

// Manager handles CRUD operations for workflows.
type Manager struct {
	loader    *Loader
	mu        sync.RWMutex
	workflows map[string]*Workflow
}

// NewManager creates a new workflow manager.
func NewManager(workflowsDir string) (*Manager, error) {
	loader := NewLoader(workflowsDir)
	if err := loader.EnsureDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure workflows directory: %w", err)
	}

	m := &Manager{
		loader:    loader,
		workflows: make(map[string]*Workflow),
	}

	// Load existing workflows
	workflows, err := loader.LoadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load workflows: %w", err)
	}

	for _, wf := range workflows {
		m.workflows[wf.ID] = wf
	}

	return m, nil
}

// GetAll returns all workflows sorted by creation date (newest first).
func (m *Manager) GetAll() []*Workflow {
	m.mu.RLock()
	defer m.mu.RUnlock()

	workflows := make([]*Workflow, 0, len(m.workflows))
	for _, wf := range m.workflows {
		workflows = append(workflows, wf.Clone())
	}

	sort.Slice(workflows, func(i, j int) bool {
		return workflows[i].CreatedAt.After(workflows[j].CreatedAt)
	})

	return workflows
}

// GetByID returns a workflow by its ID.
func (m *Manager) GetByID(id string) (*Workflow, bool) {
	if IsDefaultWorkflow(id) {
		return DefaultWorkflow(), true
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	wf, ok := m.workflows[id]
	if !ok {
		return nil, false
	}
	return wf.Clone(), true
}

// Create creates a new workflow from the input.
func (m *Manager) Create(input WorkflowInput) (*Workflow, error) {
	now := time.Now()

	// Ensure the Transcribe entry node always exists
	nodes := ensureTranscribeNode(input.Nodes)

	wf := &Workflow{
		ID:          uuid.Must(uuid.NewV7()).String(),
		Name:        input.Name,
		Description: input.Description,
		Enabled:     input.Enabled,
		CreatedAt:   now,
		UpdatedAt:   now,
		Trigger:     input.Trigger,
		Nodes:       nodes,
		Edges:       input.Edges,
	}

	if err := ValidateWorkflow(wf); err != nil {
		return nil, fmt.Errorf("invalid workflow: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.loader.Save(wf); err != nil {
		return nil, err
	}

	m.workflows[wf.ID] = wf
	return wf.Clone(), nil
}

// Update updates an existing workflow.
func (m *Manager) Update(id string, input WorkflowInput) (*Workflow, error) {
	if IsDefaultWorkflow(id) {
		return nil, errors.New("cannot modify default workflow")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	existing, ok := m.workflows[id]
	if !ok {
		return nil, ErrWorkflowNotFound
	}

	// Ensure the Transcribe entry node always exists
	nodes := ensureTranscribeNode(input.Nodes)

	wf := &Workflow{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Enabled:     input.Enabled,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   time.Now(),
		Trigger:     input.Trigger,
		Nodes:       nodes,
		Edges:       input.Edges,
	}

	if err := ValidateWorkflow(wf); err != nil {
		return nil, fmt.Errorf("invalid workflow: %w", err)
	}

	if err := m.loader.Save(wf); err != nil {
		return nil, err
	}

	m.workflows[id] = wf
	return wf.Clone(), nil
}

// Delete removes a workflow.
func (m *Manager) Delete(id string) error {
	if IsDefaultWorkflow(id) {
		return errors.New("cannot delete default workflow")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.workflows[id]; !ok {
		return ErrWorkflowNotFound
	}

	if err := m.loader.Delete(id); err != nil {
		return err
	}

	delete(m.workflows, id)
	return nil
}

// Duplicate creates a copy of an existing workflow.
func (m *Manager) Duplicate(id string) (*Workflow, error) {
	wf, ok := m.GetByID(id)
	if !ok {
		return nil, ErrWorkflowNotFound
	}

	input := WorkflowInput{
		Name:        wf.Name + " (Copy)",
		Description: wf.Description,
		Enabled:     wf.Enabled,
		Trigger:     wf.Trigger,
		Nodes:       wf.Nodes,
		Edges:       wf.Edges,
	}

	return m.Create(input)
}

// Count returns the number of workflows (excluding default).
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.workflows)
}

// ensureTranscribeNode ensures that a Transcribe node always exists in the workflow.
// If not present, it adds one with a fixed ID at position (100, 100).
func ensureTranscribeNode(nodes []Node) []Node {
	// Check if a Transcribe node already exists
	for _, node := range nodes {
		if node.Type == NodeTypeTranscribe {
			return nodes
		}
	}

	// Add the entry Transcribe node
	transcribeNode := Node{
		ID:   transcribeNodeID,
		Type: NodeTypeTranscribe,
		Position: Position{
			X: 100,
			Y: 100,
		},
		Config: nil,
		Data:   nil,
	}

	return append([]Node{transcribeNode}, nodes...)
}
