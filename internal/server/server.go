package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/server/api"
	"github.com/varavelio/tribar/internal/state"
)

func NewServer(
	logger logger.Logger,
	settingsManager *config.SettingsManager,
	appState *state.Instance,
) *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	apiGroup := server.Group("/api/v1")
	api.MountRouter(apiGroup, logger, settingsManager, appState)

	return server
}
