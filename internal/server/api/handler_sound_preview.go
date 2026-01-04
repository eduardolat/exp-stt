package api

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/varavelio/tribar/assets/sounds"
)

// handleSoundPreview serves embedded WAV sound files for browser preview.
func (h *handlers) handleSoundPreview(c echo.Context) error {
	soundType := c.Param("type")
	soundID := c.Param("id")

	var audioData []byte

	switch soundType {
	case "record":
		sound, ok := sounds.InOutSoundsMap[soundID]
		if !ok {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "sound not found",
			})
		}
		audioData = sound.Input

	case "success":
		sound, ok := sounds.SuccessSoundsMap[soundID]
		if !ok {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "sound not found",
			})
		}
		audioData = sound.Sound

	case "error":
		sound, ok := sounds.ErrorSoundsMap[soundID]
		if !ok {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "sound not found",
			})
		}
		audioData = sound.Sound

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid sound type, must be one of: record, success, error",
		})
	}

	c.Response().Header().Set("Content-Type", "audio/wav")
	c.Response().Header().Set("Cache-Control", "public, max-age=31536000")

	return c.Stream(http.StatusOK, "audio/wav", bytes.NewReader(audioData))
}
