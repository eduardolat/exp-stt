package api

import (
	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

func (h *handlers) registerProcWorkflowCreate() {
	h.uforpcServer.Procs.WorkflowCreate.Handle(func(c *uforpc.WorkflowCreateHandlerContext[urpcProps]) (uforpc.WorkflowCreateOutput, error) {
		input := workflowInputFromAPI(c.Input.Workflow)
		wf, err := h.workflowManager.Create(input)
		if err != nil {
			return uforpc.WorkflowCreateOutput{}, err
		}
		return uforpc.WorkflowCreateOutput{Workflow: workflowToAPI(wf)}, nil
	})
}
