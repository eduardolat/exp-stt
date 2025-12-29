package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcRecordingToggle() {
	h.uforpcServer.Procs.RecordingToggle.Handle(func(c *uforpc.RecordingToggleHandlerContext[urpcProps]) (uforpc.RecordingToggleOutput, error) {
		h.engine.ToggleRecording()
		return uforpc.RecordingToggleOutput{}, nil
	})
}
