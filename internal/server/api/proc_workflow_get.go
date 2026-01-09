package api

import (
	"fmt"

	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

func (h *handlers) registerProcWorkflowGet() {
	h.uforpcServer.Procs.WorkflowGet.Handle(func(c *uforpc.WorkflowGetHandlerContext[urpcProps]) (uforpc.WorkflowGetOutput, error) {
		wf, ok := h.workflowManager.GetByID(c.Input.Id)
		if !ok {
			return uforpc.WorkflowGetOutput{}, fmt.Errorf("workflow not found")
		}
		return uforpc.WorkflowGetOutput{Workflow: workflowToAPI(wf)}, nil
	})
}
