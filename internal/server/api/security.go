package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// requireJSONContentType enforces that the request has a Content-Type header
// starting with "application/json". This helps prevent CSRF attacks via
// simple requests (like text/plain forms).
func requireJSONContentType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ct := c.Request().Header.Get("Content-Type")
		if !strings.HasPrefix(strings.ToLower(ct), "application/json") {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{
				"error": "Content-Type must be application/json",
			})
		}
		return next(c)
	}
}
