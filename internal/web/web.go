package web

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"minwa/internal/database"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HTTPServer is a wrapper around echo.Echo with a logger and database queries
type HTTPServer struct {
	Server *echo.Echo
	Logger *slog.Logger

	Queries database.Queries
	Pass    string
}

// NewHttpServer creates a new http server and sets up logging, routes and static files
func NewHttpServer(db *sql.DB, pass string) *HTTPServer {
	hs := &HTTPServer{
		Server:  echo.New(),
		Queries: *database.New(db),
		Logger:  slog.Default(),
		Pass:    pass,
	}

	hs.Server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				hs.Logger.LogAttrs(context.Background(), slog.LevelInfo, "req",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				hs.Logger.LogAttrs(context.Background(), slog.LevelError, "req",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	// setup routes
	hs.routes()

	// static files
	hs.Server.Static("/static", "web/static")

	// setup http basic auth
	hs.Server.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "user" && password == hs.Pass {
			return true, nil
		}
		return false, nil
	}))

	return hs
}

// Render renders a templ component with a status code
func (hs *HTTPServer) Render(status int, c echo.Context, view templ.Component) error {
	c.Response().Status = status
	return view.Render(context.Background(), c.Response().Writer)
}

// Error is a global error handler for the http server
func (hs *HTTPServer) Error(c echo.Context, err error) error {
	hs.Logger.Error("http handler failed", "error", err)
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"error":  "internal error",
		"status": "error",
	})
}
