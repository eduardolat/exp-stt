package api

import (
	"embed"
	"io/fs"

	"github.com/labstack/echo/v4"
	"github.com/varavelio/tribar/internal/config"
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
	uforpcServer    *uforpc.Server[urpcProps]
}

func MountRouter(

	parent *echo.Group,
	logger logger.Logger,
	settingsManager *config.SettingsManager,
	appState *state.Instance,
) {
	uforpcServer := uforpc.NewServer[urpcProps]()
	handlers := &handlers{
		logger:          logger,
		settingsManager: settingsManager,
		appState:        appState,
		uforpcServer:    uforpcServer,
	}

	handlers.registerURPC()

	subFS, _ := fs.Sub(playgroundFS, "playground")
	parent.Group("/playground").StaticFS("", subFS)
	parent.POST("/urpc/:operationName", handlers.handleURPC)

	// Add other handlers not related to UFO RPC here
}

// registerURPC registra todos los handlers y hooks de UFO RPC
func (h *handlers) registerURPC() {
	h.registerProcBar()
}

func (h *handlers) handleURPC(c echo.Context) error {
	ctx := c.Request().Context()
	props := urpcProps{}

	operationName := c.Param("operationName")
	httpAdapter := uforpc.NewNetHTTPAdapter(c.Response(), c.Request())

	return h.uforpcServer.HandleRequest(ctx, props, operationName, httpAdapter)
}
