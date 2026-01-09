package api

import (
	"encoding/json"

	"github.com/varavelio/tribar/internal/server/api/uforpc"
	"github.com/varavelio/tribar/internal/workflow"
)

// workflowToAPI converts a workflow.Workflow to uforpc.Workflow
func workflowToAPI(wf *workflow.Workflow) uforpc.Workflow {
	nodes := make([]uforpc.WorkflowNode, len(wf.Nodes))
	for i, n := range wf.Nodes {
		configJSON, _ := json.Marshal(n.Config)
		dataJSON, _ := json.Marshal(n.Data)
		nodes[i] = uforpc.WorkflowNode{
			Id:       n.ID,
			NodeType: string(n.Type),
			Position: uforpc.Position{X: n.Position.X, Y: n.Position.Y},
			Config:   string(configJSON),
			Data:     string(dataJSON),
		}
	}

	edges := make([]uforpc.WorkflowEdge, len(wf.Edges))
	for i, e := range wf.Edges {
		edges[i] = uforpc.WorkflowEdge{
			Id:           e.ID,
			Source:       e.Source,
			Target:       e.Target,
			SourceHandle: uforpc.Optional[string]{Present: e.SourceHandle != "", Value: e.SourceHandle},
			TargetHandle: uforpc.Optional[string]{Present: e.TargetHandle != "", Value: e.TargetHandle},
		}
	}

	triggerConfigJSON, _ := json.Marshal(wf.Trigger.Config)

	return uforpc.Workflow{
		Id:          wf.ID,
		Name:        wf.Name,
		Description: uforpc.Optional[string]{Present: wf.Description != "", Value: wf.Description},
		Enabled:     wf.Enabled,
		CreatedAt:   wf.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   wf.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Trigger: uforpc.WorkflowTrigger{
			TriggerType: string(wf.Trigger.Type),
			Config:      string(triggerConfigJSON),
		},
		Nodes: nodes,
		Edges: edges,
	}
}

// workflowInputFromAPI converts uforpc.WorkflowInput to workflow.WorkflowInput
func workflowInputFromAPI(input uforpc.WorkflowInput) workflow.WorkflowInput {
	nodes := make([]workflow.Node, len(input.Nodes))
	for i, n := range input.Nodes {
		var config, data map[string]interface{}
		_ = json.Unmarshal([]byte(n.Config), &config)
		_ = json.Unmarshal([]byte(n.Data), &data)

		nodes[i] = workflow.Node{
			ID:       n.Id,
			Type:     workflow.NodeType(n.NodeType),
			Position: workflow.Position{X: n.Position.X, Y: n.Position.Y},
			Config:   config,
			Data:     data,
		}
	}

	edges := make([]workflow.Edge, len(input.Edges))
	for i, e := range input.Edges {
		edge := workflow.Edge{
			ID:     e.Id,
			Source: e.Source,
			Target: e.Target,
		}
		if e.SourceHandle.Present {
			edge.SourceHandle = e.SourceHandle.Value
		}
		if e.TargetHandle.Present {
			edge.TargetHandle = e.TargetHandle.Value
		}
		edges[i] = edge
	}

	var triggerConfig map[string]interface{}
	_ = json.Unmarshal([]byte(input.Trigger.Config), &triggerConfig)

	description := ""
	if input.Description.Present {
		description = input.Description.Value
	}

	return workflow.WorkflowInput{
		Name:        input.Name,
		Description: description,
		Enabled:     input.Enabled,
		Trigger: workflow.Trigger{
			Type:   workflow.TriggerType(input.Trigger.TriggerType),
			Config: triggerConfig,
		},
		Nodes: nodes,
		Edges: edges,
	}
}

// nodeDefinitionsToAPI converts workflow node definitions to API format
func nodeDefinitionsToAPI(defs []workflow.NodeDefinition) []uforpc.NodeDefinition {
	result := make([]uforpc.NodeDefinition, len(defs))
	for i, d := range defs {
		outputsJSON, _ := json.Marshal(d.Outputs)
		result[i] = uforpc.NodeDefinition{
			NodeType:    string(d.Type),
			Name:        d.Name,
			Description: d.Description,
			Category:    d.Category,
			Inputs:      d.Inputs,
			Outputs:     string(outputsJSON),
		}
	}
	return result
}
