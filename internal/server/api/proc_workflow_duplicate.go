package api

import (
	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

func (h *handlers) registerProcWorkflowDuplicate() {
	h.uforpcServer.Procs.WorkflowDuplicate.Handle(func(c *uforpc.WorkflowDuplicateHandlerContext[urpcProps]) (uforpc.WorkflowDuplicateOutput, error) {
		wf, err := h.workflowManager.Duplicate(c.Input.Id)
		if err != nil {
			return uforpc.WorkflowDuplicateOutput{}, err
		}
		return uforpc.WorkflowDuplicateOutput{Workflow: workflowToAPI(wf)}, nil
	})
}
