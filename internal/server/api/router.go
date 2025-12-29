package api

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/engine"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/server/api/uforpc"
	"github.com/varavelio/tribar/internal/state"
)

//go:embed all:playground/**
var playgroundFS embed.FS

// urpcProps contains the properties passed to all UFO RPC handlers and can be
// modified by middlewares
type urpcProps struct{}

type handlers struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager
	appState        *state.Instance
	engine          *engine.Engine
	uforpcServer    *uforpc.Server[urpcProps]
}

func MountRouter(
	parent *echo.Group,
	logger logger.Logger,
	settingsManager *config.SettingsManager,
	appState *state.Instance,
	eng *engine.Engine,
) {
	uforpcServer := uforpc.NewServer[urpcProps]()
	handlers := &handlers{
		logger:          logger,
		settingsManager: settingsManager,
		appState:        appState,
		engine:          eng,
		uforpcServer:    uforpcServer,
	}

	handlers.registerURPC()

	subFS, _ := fs.Sub(playgroundFS, "playground")
	parent.Group("/playground").StaticFS("", subFS)
	parent.POST("/urpc/:operationName", handlers.handleURPC)

	parent.GET("/audio/:id", handlers.handleAudioStream)
}

// registerURPC registers all UFO RPC handlers
func (h *handlers) registerURPC() {
	// Procedures
	h.registerProcStateGet()
	h.registerProcRecordingToggle()
	h.registerProcSettingsGet()
	h.registerProcSettingsUpdate()
	h.registerProcHistoryDeleteEntry()
	h.registerProcHistoryClear()

	// Streams
	h.registerStreamListenForEvents()
}

func (h *handlers) handleURPC(c echo.Context) error {
	ctx := c.Request().Context()
	props := urpcProps{}

	operationName := c.Param("operationName")
	httpAdapter := uforpc.NewNetHTTPAdapter(c.Response(), c.Request())

	return h.uforpcServer.HandleRequest(ctx, props, operationName, httpAdapter)
}

// handleAudioStream serves WAV audio files from history entries.
func (h *handlers) handleAudioStream(c echo.Context) error {
	id := c.Param("id")

	// Validate that the entry exists
	if _, ok := h.appState.GetHistoryEntry(id); !ok {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "entry not found",
		})
	}

	audioPath := h.appState.GetHistoryAudioPath(id)

	// Check if file exists
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "audio file not found",
		})
	}

	// Set appropriate headers for audio streaming
	c.Response().Header().Set("Content-Type", "audio/wav")
	c.Response().Header().Set("Accept-Ranges", "bytes")

	// Use filename for download
	filename := id + ".wav"
	c.Response().Header().Set("Content-Disposition", "inline; filename=\""+filename+"\"")

	// Enable range requests for seeking
	return c.File(audioPath)
}
