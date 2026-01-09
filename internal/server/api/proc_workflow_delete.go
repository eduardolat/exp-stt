package api

import (
	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

func (h *handlers) registerProcWorkflowDelete() {
	h.uforpcServer.Procs.WorkflowDelete.Handle(func(c *uforpc.WorkflowDeleteHandlerContext[urpcProps]) (uforpc.WorkflowDeleteOutput, error) {
		if err := h.workflowManager.Delete(c.Input.Id); err != nil {
			return uforpc.WorkflowDeleteOutput{}, err
		}
		return uforpc.WorkflowDeleteOutput{}, nil
	})
}
