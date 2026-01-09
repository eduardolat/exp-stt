package api

import (
	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

func (h *handlers) registerProcWorkflowGetActive() {
	h.uforpcServer.Procs.WorkflowGetActive.Handle(func(c *uforpc.WorkflowGetActiveHandlerContext[urpcProps]) (uforpc.WorkflowGetActiveOutput, error) {
		settings := h.settingsManager.Get()

		if settings.ActiveWorkflowID == "" {
			return uforpc.WorkflowGetActiveOutput{
				Workflow: uforpc.Optional[uforpc.Workflow]{Present: false},
			}, nil
		}

		wf, ok := h.workflowManager.GetByID(settings.ActiveWorkflowID)
		if !ok {
			return uforpc.WorkflowGetActiveOutput{
				Workflow: uforpc.Optional[uforpc.Workflow]{Present: false},
			}, nil
		}

		apiWf := workflowToAPI(wf)
		return uforpc.WorkflowGetActiveOutput{
			Workflow: uforpc.Optional[uforpc.Workflow]{Present: true, Value: apiWf},
		}, nil
	})
}
