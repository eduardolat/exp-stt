package server

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/engine"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/server/api"
	"github.com/varavelio/tribar/internal/shortcut"
	"github.com/varavelio/tribar/internal/state"
	"github.com/varavelio/tribar/webapp"
)

func NewServer(
	logger logger.Logger,
	settingsManager *config.SettingsManager,
	appState *state.Instance,
	eng *engine.Engine,
	shortcutMgr *shortcut.Manager,
) *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	// Allow any origin coming from localhost (any port)
	// Very permissive since it's for local use
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Connection",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Mount API routes
	apiGroup := server.Group("/api/v1")
	api.MountRouter(apiGroup, logger, settingsManager, appState, eng, shortcutMgr)

	// Mount Web UI routes
	subFS, _ := fs.Sub(webapp.BuildFS, "build")
	server.Group("/").StaticFS("", subFS)

	return server
}
