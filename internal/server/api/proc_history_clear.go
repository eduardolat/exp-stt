package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcHistoryClear() {
	h.uforpcServer.Procs.HistoryClear.Handle(func(c *uforpc.HistoryClearHandlerContext[urpcProps]) (uforpc.HistoryClearOutput, error) {
		if err := h.appState.ClearHistory(c.Context); err != nil {
			return uforpc.HistoryClearOutput{}, uforpc.Error{
				Category: "InternalError",
				Code:     "HISTORY_CLEAR_FAILED",
				Message:  err.Error(),
			}
		}
		return uforpc.HistoryClearOutput{}, nil
	})
}
