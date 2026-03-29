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

	// Add security headers
	server.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN", // Allow embedding in same origin (e.g. if we use iframes internally)
		ReferrerPolicy:     "strict-origin-when-cross-origin",
		XSSProtection:      "1; mode=block", // Basic protection for older browsers
		// CSP:
		// - default-src 'self': Default to only allowing resources from the same origin
		// - script-src 'self' 'unsafe-inline': Allow local scripts and inline scripts (required for SvelteKit hydration)
		// - style-src 'self' 'unsafe-inline' https://raw.githubusercontent.com: Allow local styles, inline styles (for dynamic themes), and DaisyUI themes
		// - img-src 'self' data:: Allow local images and data URIs (base64)
		// - font-src 'self' data:: Allow local fonts and data URIs
		// - connect-src 'self' http: https:: Allow connections to self and external AI providers (http/https)
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' https://raw.githubusercontent.com; img-src 'self' data:; font-src 'self' data:; connect-src 'self' http: https:;",
	}))

	// Configure CORS
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// Always allow localhost/loopback for local development and usage.
			// This prevents malicious websites from accessing the local API
			// via cross-site requests (CSRF/CORS attacks).
			// We allow dynamic ports (e.g. :5173) for development.
			if origin == "http://localhost" || strings.HasPrefix(origin, "http://localhost:") ||
				origin == "https://localhost" || strings.HasPrefix(origin, "https://localhost:") ||
				origin == "http://127.0.0.1" || strings.HasPrefix(origin, "http://127.0.0.1:") ||
				origin == "https://127.0.0.1" || strings.HasPrefix(origin, "https://127.0.0.1:") {
				return true, nil
			}

			// Check user-configured allowed origins
			for _, allowed := range settingsManager.Get().AllowedCORSOrigins {
				// Wildcard allows all origins
				if allowed == "*" {
					return true, nil
				}
				// Exact match
				if allowed == origin {
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
