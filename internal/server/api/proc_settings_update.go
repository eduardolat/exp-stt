package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcSettingsUpdate() {
	h.uforpcServer.Procs.SettingsUpdate.Handle(func(c *uforpc.SettingsUpdateHandlerContext[urpcProps]) (uforpc.SettingsUpdateOutput, error) {
		settings := settingsFromAPI(c.Input.Settings)
		if err := h.settingsManager.Update(settings); err != nil {
			return uforpc.SettingsUpdateOutput{}, uforpc.Error{
				Category: "InternalError",
				Code:     "SETTINGS_UPDATE_FAILED",
				Message:  err.Error(),
			}
		}
		return uforpc.SettingsUpdateOutput{}, nil
	})
}
