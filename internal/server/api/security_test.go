package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRequireJSONContentType(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock handler that returns 200 OK
	h := requireJSONContentType(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	tests := []struct {
		name           string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "Valid Content-Type",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid Content-Type with Charset",
			contentType:    "application/json; charset=utf-8",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Content-Type Text",
			contentType:    "text/plain",
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Invalid Content-Type Form",
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Empty Content-Type",
			contentType:    "",
			expectedStatus: http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req.Header.Set("Content-Type", tt.contentType)
			rec = httptest.NewRecorder()
			c = e.NewContext(req, rec)

			// We need to reset the request in the context if we reuse it,
			// but creating a new context is cleaner, however `e.NewContext` is cheap.
			// Actually `c.Request().Header` is what matters.

			if err := h(c); err != nil {
				// Middleware might return error (e.g. echo.NewHTTPError)
				// echo.Context.JSON returns error
				if he, ok := err.(*echo.HTTPError); ok {
					if he.Code != tt.expectedStatus {
						t.Errorf("expected status %d, got %d", tt.expectedStatus, he.Code)
					}
				} else {
					// If it's not HTTPError, check if we handled it manually in middleware
					// In our case we return c.JSON which returns nil if successful or error if failed to write
					// Wait, c.JSON returns error if it fails to write response?
					// No, it writes status code to response and returns error only on IO failure.
					// But we should check rec.Code
				}
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
