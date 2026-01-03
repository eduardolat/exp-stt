package server

import (
	"io/fs"
	"net/http"
	"strings"

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

	// Add security headers
	server.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XContentTypeOptions: "nosniff",
		XFrameOptions:       "SAMEORIGIN", // Allow embedding in same origin (e.g. if we use iframes internally)
		ReferrerPolicy:      "strict-origin-when-cross-origin",
		XXSSProtection:      "1; mode=block", // Basic protection for older browsers
	}))

	// Configure CORS
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// Check if user has explicitly allowed all origins
			if settingsManager.Get().AllowExternalOrigins {
				return true, nil
			}

			// Default: Restrict to localhost/loopback
			// This prevents malicious websites from accessing the local API
			// via cross-site requests (CSRF/CORS attacks).
			// We allow dynamic ports (e.g. :5173) for development.
			if origin == "http://localhost" || strings.HasPrefix(origin, "http://localhost:") ||
				origin == "https://localhost" || strings.HasPrefix(origin, "https://localhost:") ||
				origin == "http://127.0.0.1" || strings.HasPrefix(origin, "http://127.0.0.1:") ||
				origin == "https://127.0.0.1" || strings.HasPrefix(origin, "https://127.0.0.1:") {
				return true, nil
			}

			// Also allow the app's own origin if it's serving from a different host/port binding
			// (though usually covered by localhost checks above if bound to localhost)

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
	api.MountRouter(apiGroup, logger, settingsManager, appState, eng, shortcutMgr)

	// Mount Web UI routes
	subFS, _ := fs.Sub(webapp.BuildFS, "build")
	server.Group("/").StaticFS("", subFS)

	return server
}
