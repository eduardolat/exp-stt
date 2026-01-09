package executor

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/varavelio/tribar/internal/workflow"
	"github.com/varavelio/tribar/internal/workflow/executor/nodes"
)

// Executor runs workflows by traversing the node graph in topological order.
type Executor struct {
	services  *ServiceContainer
	registry  *NodeRegistry
	eventChan chan ExecutionEvent
	isRunning atomic.Bool
}

// New creates a new workflow executor.
func New(services *ServiceContainer, registry *NodeRegistry) *Executor {
	return &Executor{
		services:  services,
		registry:  registry,
		eventChan: make(chan ExecutionEvent, 100),
	}
}

// EventChan returns the channel for execution events (for SSE streaming).
func (e *Executor) EventChan() <-chan ExecutionEvent {
	return e.eventChan
}

// IsRunning returns true if a workflow is currently executing.
func (e *Executor) IsRunning() bool {
	return e.isRunning.Load()
}

// Execute runs the given workflow with the provided execution context.
// Returns an error if the workflow fails or if another workflow is already running.
func (e *Executor) Execute(ctx context.Context, wf *workflow.Workflow, execCtx *ExecutionContext) error {
	// Single execution policy: prevent concurrent workflow execution
	if !e.isRunning.CompareAndSwap(false, true) {
		return fmt.Errorf("workflow already in progress")
	}
	defer e.isRunning.Store(false)

	// Build execution order using topological sort
	order, err := e.topologicalSort(wf)
	if err != nil {
		return fmt.Errorf("failed to determine execution order: %w", err)
	}

	// Execute nodes in order
	for _, node := range order {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := e.executeNode(ctx, node, wf, execCtx); err != nil {
			return fmt.Errorf("node %s (%s) failed: %w", node.ID, node.Type, err)
		}
	}

	return nil
}

// executeNode runs a single node and stores its output.
func (e *Executor) executeNode(ctx context.Context, node *workflow.Node, wf *workflow.Workflow, execCtx *ExecutionContext) error {
	// Emit "running" event
	e.emitEvent(NewExecutionEvent(execCtx.ExecutionID, node.ID, StatusRunning))

	// Get the executor for this node type
	executor, ok := e.registry.Get(node.Type)
	if !ok {
		err := fmt.Errorf("unknown node type: %s", node.Type)
		e.emitEvent(NewExecutionEvent(execCtx.ExecutionID, node.ID, StatusError).WithError(err))
		return err
	}

	// Check if this is a condition node - handle branching
	if node.Type == workflow.NodeTypeCondition {
		return e.executeConditionNode(ctx, node, wf, execCtx, executor)
	}

	// Resolve variables in config
	resolvedConfig := execCtx.ResolveConfig(node.Config)

	// Build input data for the node
	// Priority: 1) Trigger data for transcribe, 2) Parent node output, 3) Node's own data
	var nodeData map[string]interface{}
	if node.Type == workflow.NodeTypeTranscribe {
		// Entry point receives trigger data
		nodeData = execCtx.TriggerData
	} else {
		// Find parent node (node that has an edge targeting this node)
		nodeData = e.getParentNodeOutput(node.ID, wf, execCtx)
	}

	input := nodes.NodeInput{
		Config: resolvedConfig,
		Data:   nodeData,
	}

	// Execute the node
	output, err := executor.Execute(ctx, input, e.services)
	if err != nil {
		e.emitEvent(NewExecutionEvent(execCtx.ExecutionID, node.ID, StatusError).WithError(err))
		return err
	}

	// Store output for downstream nodes
	execCtx.SetOutput(node.ID, output)

	// Emit "completed" event
	e.emitEvent(NewExecutionEvent(execCtx.ExecutionID, node.ID, StatusCompleted).WithOutputPreview(output.Preview()))

	return nil
}

// getParentNodeOutput finds the output of the parent node (the one connected via incoming edge).
func (e *Executor) getParentNodeOutput(nodeID string, wf *workflow.Workflow, execCtx *ExecutionContext) map[string]interface{} {
	// Find edges that target this node
	for _, edge := range wf.Edges {
		if edge.Target == nodeID {
			// Get the output of the source node
			if output, ok := execCtx.GetOutput(edge.Source); ok {
				return output.Fields()
			}
		}
	}
	return nil
}

// executeConditionNode handles condition nodes with branching logic.
func (e *Executor) executeConditionNode(ctx context.Context, node *workflow.Node, wf *workflow.Workflow, execCtx *ExecutionContext, executor nodes.NodeExecutor) error {
	// Resolve variables in config
	resolvedConfig := execCtx.ResolveConfig(node.Config)

	input := nodes.NodeInput{
		Config: resolvedConfig,
		Data:   node.Data,
	}

	output, err := executor.Execute(ctx, input, e.services)
	if err != nil {
		e.emitEvent(NewExecutionEvent(execCtx.ExecutionID, node.ID, StatusError).WithError(err))
		return err
	}

	// Store condition result
	execCtx.SetOutput(node.ID, output)

	// The condition node will output which branch to take
	e.emitEvent(NewExecutionEvent(execCtx.ExecutionID, node.ID, StatusCompleted).WithOutputPreview(output.Preview()))

	return nil
}

// topologicalSort returns nodes in execution order using Kahn's algorithm.
func (e *Executor) topologicalSort(wf *workflow.Workflow) ([]*workflow.Node, error) {
	// Build node map for quick lookup
	nodeMap := make(map[string]*workflow.Node, len(wf.Nodes))
	for i := range wf.Nodes {
		nodeMap[wf.Nodes[i].ID] = &wf.Nodes[i]
	}

	// Build adjacency list and in-degree count
	adj := make(map[string][]string)
	inDegree := make(map[string]int)

	for _, node := range wf.Nodes {
		inDegree[node.ID] = 0
		adj[node.ID] = nil
	}

	for _, edge := range wf.Edges {
		adj[edge.Source] = append(adj[edge.Source], edge.Target)
		inDegree[edge.Target]++
	}

	// Start with nodes that have no incoming edges
	queue := make([]string, 0)
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	// Process nodes
	var result []*workflow.Node
	for len(queue) > 0 {
		nodeID := queue[0]
		queue = queue[1:]

		node := nodeMap[nodeID]
		result = append(result, node)

		for _, neighbor := range adj[nodeID] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check if all nodes were processed (no cycles)
	if len(result) != len(wf.Nodes) {
		return nil, fmt.Errorf("workflow contains a cycle")
	}

	return result, nil
}

// emitEvent sends an execution event to the event channel.
func (e *Executor) emitEvent(event ExecutionEvent) {
	select {
	case e.eventChan <- event:
	default:
		// Channel full, drop event to prevent blocking
	}
}
