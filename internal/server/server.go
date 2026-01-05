package server

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/engine"
	"github.com/varavelio/tribar/internal/eventbus"
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
	eventBus *eventbus.EventBus,
) *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	// Allow any origin coming from localhost (any port) or configured allowed origins
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// Check if origin is localhost/127.0.0.1
			if origin == "http://localhost" || strings.HasPrefix(origin, "http://localhost:") ||
				origin == "http://127.0.0.1" || strings.HasPrefix(origin, "http://127.0.0.1:") {
				return true, nil
			}

			// Check configured allowed origins
			settings := settingsManager.Get()
			for _, allowed := range settings.AllowedOrigins {
				if origin == allowed {
					return true, nil
				}
			}

			return false, nil
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
	api.MountRouter(apiGroup, logger, settingsManager, appState, eng, shortcutMgr, eventBus)

	// Mount Web UI routes
	subFS, _ := fs.Sub(webapp.BuildFS, "build")
	server.Group("/").StaticFS("", subFS)

	return server
}
