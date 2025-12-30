package api

import "github.com/varavelio/tribar/internal/server/api/uforpc"

func (h *handlers) registerProcShortcutToggleUpdate() {
	h.uforpcServer.Procs.ShortcutToggleUpdate.Handle(func(c *uforpc.ShortcutToggleUpdateHandlerContext[urpcProps]) (uforpc.ShortcutToggleUpdateOutput, error) {
		if h.shortcutManager == nil {
			return uforpc.ShortcutToggleUpdateOutput{}, uforpc.Error{
				Code:    "SHORTCUT_MANAGER_NOT_INITIALIZED",
				Message: "shortcut manager not initialized",
			}
		}

		shortcut := shortcutFromAPI(c.Input.Shortcut)
		if err := h.shortcutManager.Update(shortcut); err != nil {
			return uforpc.ShortcutToggleUpdateOutput{}, uforpc.Error{
				Code:    "SHORTCUT_UPDATE_FAILED",
				Message: "shortcut update failed",
				Details: map[string]any{
					"error": err.Error(),
				},
			}
		}

		return uforpc.ShortcutToggleUpdateOutput{}, nil
	})
}
