package web

import (
	"net/http"
	"strconv"

	"minwa/internal/database"
	"minwa/web/templates"

	"github.com/labstack/echo/v4"
)

const checks = 50

func (hs *HTTPServer) routes() {
	hs.Server.GET("/", hs.indexHandler)
	hs.Server.POST("/", hs.addHandler)
	hs.Server.POST("/:id/delete", hs.deleteHandler)
}

func (hs *HTTPServer) deleteHandler(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return hs.Error(c, err)
	}

	if err := hs.Queries.EndpointsDelete(ctx, int64(id)); err != nil {
		return hs.Error(c, err)
	}

	if err := hs.Queries.ChecksDelete(ctx, int64(id)); err != nil {
		return hs.Error(c, err)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (hs *HTTPServer) indexHandler(c echo.Context) error {
	ctx := c.Request().Context()

	endpoints, err := hs.Queries.EndpointsList(c.Request().Context())
	if err != nil {
		return hs.Error(c, err)
	}

	var ec []templates.EndpointWithChecks

	for _, ep := range endpoints {
		checks, err := hs.Queries.ChecksForEndpoint(ctx, database.ChecksForEndpointParams{
			EndpointID: ep.ID,
			Limit:      checks,
		})
		if err != nil {
			return hs.Error(c, err)
		}

		ec = append(ec, templates.EndpointWithChecks{
			Endpoint: ep,
			Checks:   checks,
		})
	}

	return hs.Render(http.StatusOK, c, templates.Index(ec))
}

func (hs *HTTPServer) addHandler(c echo.Context) error {
	ctx := c.Request().Context()

	url := c.FormValue("url")

	sc, err := strconv.Atoi(c.FormValue("expected_status"))
	if err != nil {
		return hs.Error(c, err)
	}

	if err := hs.Queries.EndpointsCreate(ctx, database.EndpointsCreateParams{
		Url:            url,
		ExpectedStatus: int64(sc),
	}); err != nil {
		return hs.Error(c, err)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
