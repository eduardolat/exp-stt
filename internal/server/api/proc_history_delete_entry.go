package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcHistoryDeleteEntry() {
	h.uforpcServer.Procs.HistoryDeleteEntry.Handle(func(c *uforpc.HistoryDeleteEntryHandlerContext[urpcProps]) (uforpc.HistoryDeleteEntryOutput, error) {
		if err := h.appState.DeleteHistoryEntry(c.Context, c.Input.Id); err != nil {
			return uforpc.HistoryDeleteEntryOutput{}, uforpc.Error{
				Category: "NotFoundError",
				Code:     "ENTRY_NOT_FOUND",
				Message:  err.Error(),
			}
		}
		return uforpc.HistoryDeleteEntryOutput{}, nil
	})
}
