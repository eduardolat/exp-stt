package api

import (
	"github.com/varavelio/tribar/internal/server/api/uforpc"
	"github.com/varavelio/tribar/internal/workflow"
)

func (h *handlers) registerProcWorkflowGetAvailableNodes() {
	h.uforpcServer.Procs.WorkflowGetAvailableNodes.Handle(func(c *uforpc.WorkflowGetAvailableNodesHandlerContext[urpcProps]) (uforpc.WorkflowGetAvailableNodesOutput, error) {
		definitions := workflow.GetAvailableNodeDefinitions()
		return uforpc.WorkflowGetAvailableNodesOutput{
			Nodes: nodeDefinitionsToAPI(definitions),
		}, nil
	})
}
