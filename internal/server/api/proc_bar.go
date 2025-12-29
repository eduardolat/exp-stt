package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcBar() {
	h.uforpcServer.Procs.BarProc.Handle(func(c *uforpc.BarProcHandlerContext[urpcProps]) (uforpc.BarProcOutput, error) {
		return uforpc.BarProcOutput{}, nil
	})
}
