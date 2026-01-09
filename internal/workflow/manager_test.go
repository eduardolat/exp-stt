package workflow

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	t.Run("create and retrieve workflow", func(t *testing.T) {
		tempDir := t.TempDir()
		m, err := NewManager(tempDir)
		require.NoError(t, err)

		input := WorkflowInput{
			Name:        "Test Workflow",
			Description: "A test workflow",
			Enabled:     true,
			Trigger:     Trigger{Type: TriggerTypeVoice},
			Nodes: []Node{
				{ID: "n1", Type: NodeTypeTranscribe, Position: Position{X: 100, Y: 100}},
			},
			Edges: []Edge{},
		}

		created, err := m.Create(input)
		require.NoError(t, err)
		require.NotEmpty(t, created.ID)
		require.Equal(t, "Test Workflow", created.Name)

		// Retrieve by ID
		retrieved, ok := m.GetByID(created.ID)
		require.True(t, ok)
		require.Equal(t, created.ID, retrieved.ID)
		require.Equal(t, created.Name, retrieved.Name)

		// Verify persistence
		require.FileExists(t, filepath.Join(tempDir, created.ID+".json"))
	})

	t.Run("get all workflows", func(t *testing.T) {
		tempDir := t.TempDir()
		m, err := NewManager(tempDir)
		require.NoError(t, err)

		// Create two workflows
		for i := 0; i < 2; i++ {
			_, err := m.Create(WorkflowInput{
				Name:    "Workflow",
				Enabled: true,
				Trigger: Trigger{Type: TriggerTypeVoice},
				Nodes:   []Node{{ID: "n1", Type: NodeTypeTranscribe}},
			})
			require.NoError(t, err)
		}

		all := m.GetAll()
		require.Len(t, all, 2)
	})

	t.Run("update workflow", func(t *testing.T) {
		tempDir := t.TempDir()
		m, err := NewManager(tempDir)
		require.NoError(t, err)

		created, err := m.Create(WorkflowInput{
			Name:    "Original Name",
			Enabled: true,
			Trigger: Trigger{Type: TriggerTypeVoice},
			Nodes:   []Node{{ID: "n1", Type: NodeTypeTranscribe}},
		})
		require.NoError(t, err)

		updated, err := m.Update(created.ID, WorkflowInput{
			Name:    "Updated Name",
			Enabled: false,
			Trigger: Trigger{Type: TriggerTypeVoice},
			Nodes:   []Node{{ID: "n1", Type: NodeTypeTranscribe}},
		})
		require.NoError(t, err)
		require.Equal(t, "Updated Name", updated.Name)
		require.False(t, updated.Enabled)
		require.Equal(t, created.CreatedAt, updated.CreatedAt)
		require.True(t, updated.UpdatedAt.After(created.UpdatedAt))
	})

	t.Run("delete workflow", func(t *testing.T) {
		tempDir := t.TempDir()
		m, err := NewManager(tempDir)
		require.NoError(t, err)

		created, err := m.Create(WorkflowInput{
			Name:    "To Delete",
			Enabled: true,
			Trigger: Trigger{Type: TriggerTypeVoice},
			Nodes:   []Node{{ID: "n1", Type: NodeTypeTranscribe}},
		})
		require.NoError(t, err)

		err = m.Delete(created.ID)
		require.NoError(t, err)

		_, ok := m.GetByID(created.ID)
		require.False(t, ok)

		// Verify file removed
		_, err = os.Stat(filepath.Join(tempDir, created.ID+".json"))
		require.True(t, os.IsNotExist(err))
	})

	t.Run("duplicate workflow", func(t *testing.T) {
		tempDir := t.TempDir()
		m, err := NewManager(tempDir)
		require.NoError(t, err)

		created, err := m.Create(WorkflowInput{
			Name:        "Original",
			Description: "Original description",
			Enabled:     true,
			Trigger:     Trigger{Type: TriggerTypeVoice},
			Nodes:       []Node{{ID: "n1", Type: NodeTypeTranscribe}},
		})
		require.NoError(t, err)

		duplicated, err := m.Duplicate(created.ID)
		require.NoError(t, err)
		require.NotEqual(t, created.ID, duplicated.ID)
		require.Equal(t, "Original (Copy)", duplicated.Name)
		require.Equal(t, "Original description", duplicated.Description)
	})

	t.Run("get default workflow", func(t *testing.T) {
		tempDir := t.TempDir()
		m, err := NewManager(tempDir)
		require.NoError(t, err)

		// Default workflow should always be available
		wf, ok := m.GetByID("default")
		require.True(t, ok)
		require.Equal(t, "default", wf.ID)
		require.Equal(t, "Default Workflow", wf.Name)
	})

	t.Run("cannot modify default workflow", func(t *testing.T) {
		tempDir := t.TempDir()
		m, err := NewManager(tempDir)
		require.NoError(t, err)

		_, err = m.Update("default", WorkflowInput{Name: "Modified"})
		require.Error(t, err)

		err = m.Delete("default")
		require.Error(t, err)
	})

	t.Run("validation on create", func(t *testing.T) {
		tempDir := t.TempDir()
		m, err := NewManager(tempDir)
		require.NoError(t, err)

		// Missing name
		_, err = m.Create(WorkflowInput{
			Name:    "",
			Enabled: true,
			Trigger: Trigger{Type: TriggerTypeVoice},
			Nodes:   []Node{{ID: "n1", Type: NodeTypeTranscribe}},
		})
		require.Error(t, err)

		// Cycle in edges
		_, err = m.Create(WorkflowInput{
			Name:    "Cyclic",
			Enabled: true,
			Trigger: Trigger{Type: TriggerTypeVoice},
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
				{ID: "b", Type: NodeTypeClipboardPaste},
			},
			Edges: []Edge{
				{ID: "e1", Source: "a", Target: "b"},
				{ID: "e2", Source: "b", Target: "a"},
			},
		})
		require.ErrorIs(t, err, ErrCycleDetected)
	})

	t.Run("persistence across restarts", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create manager and add workflow
		m1, err := NewManager(tempDir)
		require.NoError(t, err)

		created, err := m1.Create(WorkflowInput{
			Name:    "Persistent",
			Enabled: true,
			Trigger: Trigger{Type: TriggerTypeVoice},
			Nodes:   []Node{{ID: "n1", Type: NodeTypeTranscribe}},
		})
		require.NoError(t, err)

		// Create new manager (simulating restart)
		m2, err := NewManager(tempDir)
		require.NoError(t, err)

		// Should load the workflow
		retrieved, ok := m2.GetByID(created.ID)
		require.True(t, ok)
		require.Equal(t, "Persistent", retrieved.Name)
	})
}
