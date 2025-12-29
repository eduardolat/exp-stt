package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcStateGet() {
	h.uforpcServer.Procs.StateGet.Handle(func(c *uforpc.StateGetHandlerContext[urpcProps]) (uforpc.StateGetOutput, error) {
		entries := h.appState.GetHistory(c.Context)
		state := h.buildState(entries)
		return uforpc.StateGetOutput{State: state}, nil
	})
}
