package api

import (
	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

func (h *handlers) registerProcWorkflowUpdate() {
	h.uforpcServer.Procs.WorkflowUpdate.Handle(func(c *uforpc.WorkflowUpdateHandlerContext[urpcProps]) (uforpc.WorkflowUpdateOutput, error) {
		input := workflowInputFromAPI(c.Input.Workflow)
		wf, err := h.workflowManager.Update(c.Input.Id, input)
		if err != nil {
			return uforpc.WorkflowUpdateOutput{}, err
		}
		return uforpc.WorkflowUpdateOutput{Workflow: workflowToAPI(wf)}, nil
	})
}
