package workflow

import "errors"

// ErrCycleDetected is returned when a workflow contains a cycle.
var ErrCycleDetected = errors.New("workflow contains a cycle - this is not allowed")

// ErrEmptyWorkflow is returned when a workflow has no nodes.
var ErrEmptyWorkflow = errors.New("workflow has no nodes")

// ErrInvalidEdge is returned when an edge references a non-existent node.
var ErrInvalidEdge = errors.New("edge references non-existent node")

// ValidateDAG ensures the workflow graph is a valid Directed Acyclic Graph.
// It checks for cycles using DFS and validates that all edges reference existing nodes.
// Empty workflows (no nodes) are valid - they represent a blank canvas for the user.
func ValidateDAG(wf *Workflow) error {
	// Empty workflows are valid - users create them empty and add nodes later
	if len(wf.Nodes) == 0 {
		return nil
	}

	// Build node ID set for validation
	nodeSet := make(map[string]struct{}, len(wf.Nodes))
	for _, node := range wf.Nodes {
		nodeSet[node.ID] = struct{}{}
	}

	// Build adjacency list and validate edges
	adj := make(map[string][]string)
	for _, edge := range wf.Edges {
		if _, ok := nodeSet[edge.Source]; !ok {
			return ErrInvalidEdge
		}
		if _, ok := nodeSet[edge.Target]; !ok {
			return ErrInvalidEdge
		}
		adj[edge.Source] = append(adj[edge.Source], edge.Target)
	}

	// Track visit state: 0=unvisited, 1=visiting (in current path), 2=visited (fully processed)
	state := make(map[string]int)

	// DFS function to detect cycles
	var dfs func(nodeID string) bool
	dfs = func(nodeID string) bool {
		if state[nodeID] == 1 {
			return true // Cycle detected: we're revisiting a node in the current path
		}
		if state[nodeID] == 2 {
			return false // Already fully processed
		}

		state[nodeID] = 1 // Mark as visiting
		for _, neighbor := range adj[nodeID] {
			if dfs(neighbor) {
				return true
			}
		}
		state[nodeID] = 2 // Mark as fully visited
		return false
	}

	// Run DFS from each unvisited node
	for _, node := range wf.Nodes {
		if state[node.ID] == 0 {
			if dfs(node.ID) {
				return ErrCycleDetected
			}
		}
	}

	return nil
}

// ValidateWorkflow performs full validation of a workflow.
func ValidateWorkflow(wf *Workflow) error {
	if wf.Name == "" {
		return errors.New("workflow name is required")
	}
	return ValidateDAG(wf)
}
