package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcSettingsGet() {
	h.uforpcServer.Procs.SettingsGet.Handle(func(c *uforpc.SettingsGetHandlerContext[urpcProps]) (uforpc.SettingsGetOutput, error) {
		settings := h.settingsManager.Get()
		return uforpc.SettingsGetOutput{Settings: buildSettings(settings)}, nil
	})
}
