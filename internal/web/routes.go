package web

import (
	"net/http"

	"minwa/web/templates"

	"github.com/labstack/echo/v4"
)

func (hs *HTTPServer) routes() {
	hs.Server.GET("/", hs.indexHandler)
}

func (hs *HTTPServer) indexHandler(c echo.Context) error {
	return hs.Render(http.StatusOK, c, templates.Index())
}
