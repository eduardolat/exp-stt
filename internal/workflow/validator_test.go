package workflow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateDAG(t *testing.T) {
	t.Run("valid linear workflow", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
				{ID: "b", Type: NodeTypeClipboardPaste},
				{ID: "c", Type: NodeTypeSound},
			},
			Edges: []Edge{
				{ID: "e1", Source: "a", Target: "b"},
				{ID: "e2", Source: "b", Target: "c"},
			},
		}
		err := ValidateDAG(wf)
		require.NoError(t, err)
	})

	t.Run("valid branching workflow", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
				{ID: "b", Type: NodeTypeCondition},
				{ID: "c", Type: NodeTypeAIProcess},
				{ID: "d", Type: NodeTypeClipboardPaste},
			},
			Edges: []Edge{
				{ID: "e1", Source: "a", Target: "b"},
				{ID: "e2", Source: "b", Target: "c", SourceHandle: "true"},
				{ID: "e3", Source: "b", Target: "d", SourceHandle: "false"},
			},
		}
		err := ValidateDAG(wf)
		require.NoError(t, err)
	})

	t.Run("cycle detection - simple loop", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
				{ID: "b", Type: NodeTypeClipboardPaste},
			},
			Edges: []Edge{
				{ID: "e1", Source: "a", Target: "b"},
				{ID: "e2", Source: "b", Target: "a"},
			},
		}
		err := ValidateDAG(wf)
		require.ErrorIs(t, err, ErrCycleDetected)
	})

	t.Run("cycle detection - self reference", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
			},
			Edges: []Edge{
				{ID: "e1", Source: "a", Target: "a"},
			},
		}
		err := ValidateDAG(wf)
		require.ErrorIs(t, err, ErrCycleDetected)
	})

	t.Run("cycle detection - complex loop", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
				{ID: "b", Type: NodeTypeAIProcess},
				{ID: "c", Type: NodeTypeClipboardPaste},
			},
			Edges: []Edge{
				{ID: "e1", Source: "a", Target: "b"},
				{ID: "e2", Source: "b", Target: "c"},
				{ID: "e3", Source: "c", Target: "a"},
			},
		}
		err := ValidateDAG(wf)
		require.ErrorIs(t, err, ErrCycleDetected)
	})

	t.Run("empty workflow is valid", func(t *testing.T) {
		wf := &Workflow{
			Name:  "Test",
			Nodes: []Node{},
		}
		err := ValidateDAG(wf)
		require.NoError(t, err, "empty workflows should be valid - users create them empty")
	})

	t.Run("invalid edge - missing source", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
			},
			Edges: []Edge{
				{ID: "e1", Source: "nonexistent", Target: "a"},
			},
		}
		err := ValidateDAG(wf)
		require.ErrorIs(t, err, ErrInvalidEdge)
	})

	t.Run("invalid edge - missing target", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
			},
			Edges: []Edge{
				{ID: "e1", Source: "a", Target: "nonexistent"},
			},
		}
		err := ValidateDAG(wf)
		require.ErrorIs(t, err, ErrInvalidEdge)
	})

	t.Run("disconnected nodes are valid", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
				{ID: "b", Type: NodeTypeSound},
				{ID: "c", Type: NodeTypeNotify},
			},
			Edges: []Edge{}, // No connections
		}
		err := ValidateDAG(wf)
		require.NoError(t, err)
	})
}

func TestValidateWorkflow(t *testing.T) {
	t.Run("missing name", func(t *testing.T) {
		wf := &Workflow{
			Name: "",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
			},
		}
		err := ValidateWorkflow(wf)
		require.Error(t, err)
		require.Contains(t, err.Error(), "name is required")
	})

	t.Run("valid workflow", func(t *testing.T) {
		wf := &Workflow{
			Name: "Test Workflow",
			Nodes: []Node{
				{ID: "a", Type: NodeTypeTranscribe},
			},
		}
		err := ValidateWorkflow(wf)
		require.NoError(t, err)
	})
}
