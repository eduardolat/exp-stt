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

	// Allow any origin coming from localhost (any port) or configured origins
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// Check user configured origins
			for _, o := range settingsManager.Get().AllowedOrigins {
				if origin == o {
					return true, nil
				}
			}

			// Check for localhost/127.0.0.1 (common for local UI)
			// We check for "http://localhost", "http://127.0.0.1", etc.
			// Or just prefix. Ideally we parse the URL but this is simpler for the middleware func.
			// The webapp runs on localhost, likely on a random port during dev or specific port in prod.
			// Since it's a desktop app, the browser might be opening the UI from localhost.
			// A simple check is if it starts with http://localhost or http://127.0.0.1
			// But to be safer, we can try to be more specific if possible.
			// However, SvelteKit dev server uses different ports.

			// For this fix, let's allow all localhost origins as intended, but NOT everything else.
			// We check strictly for "localhost" or "127.0.0.1" domain.
			// Allowed formats:
			// http://localhost
			// http://localhost:port
			// https://localhost
			// https://localhost:port
			// http://127.0.0.1
			// http://127.0.0.1:port
			// https://127.0.0.1
			// https://127.0.0.1:port

			allowedPrefixes := []string{
				"http://localhost",
				"https://localhost",
				"http://127.0.0.1",
				"https://127.0.0.1",
			}

			for _, prefix := range allowedPrefixes {
				if origin == prefix {
					return true, nil
				}
				if strings.HasPrefix(origin, prefix+":") {
					return true, nil
				}
			}

			// Block everything else
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
