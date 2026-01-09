package api

import (
	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

func (h *handlers) registerProcWorkflowSetActive() {
	h.uforpcServer.Procs.WorkflowSetActive.Handle(func(c *uforpc.WorkflowSetActiveHandlerContext[urpcProps]) (uforpc.WorkflowSetActiveOutput, error) {
		// Validate workflow exists (optional - silently ignore if not found)
		if _, ok := h.workflowManager.GetByID(c.Input.Id); !ok {
			return uforpc.WorkflowSetActiveOutput{}, nil
		}

		// Update settings
		settings := h.settingsManager.Get()
		settings.ActiveWorkflowID = c.Input.Id
		if err := h.settingsManager.Update(settings); err != nil {
			return uforpc.WorkflowSetActiveOutput{}, err
		}
		return uforpc.WorkflowSetActiveOutput{}, nil
	})
}
