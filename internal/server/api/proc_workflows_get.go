package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcWorkflowsGet() {
	h.uforpcServer.Procs.WorkflowsGet.Handle(func(c *uforpc.WorkflowsGetHandlerContext[urpcProps]) (uforpc.WorkflowsGetOutput, error) {
		workflows := h.workflowManager.GetAll()
		result := make([]uforpc.Workflow, len(workflows))
		for i, wf := range workflows {
			result[i] = workflowToAPI(wf)
		}
		return uforpc.WorkflowsGetOutput{Workflows: result}, nil
	})
}
