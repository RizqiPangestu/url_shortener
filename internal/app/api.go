package app

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/RizqiPangestu/url_shortener/internal/core"
	"github.com/labstack/echo/v4"
)

type APIController interface {
	RegisterRoutes(ec *echo.Echo)
	Shorten(ec echo.Context) error
	Redirect(ec echo.Context) error
}

type apiController struct {
	urlService     core.URLService
	baseDomain     string
	trackerService core.TrackerService
}

func NewAPIController(urlService core.URLService, baseDomain string, trackerService core.TrackerService) APIController {
	return &apiController{
		urlService:     urlService,
		baseDomain:     baseDomain,
		trackerService: trackerService,
	}
}

func (c *apiController) RegisterRoutes(ec *echo.Echo) {
	ec.POST("/shorten", c.Shorten)
	ec.GET("u/:short_path", c.Redirect, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			err := next(ec)
			shortPath := ec.Param("short_path")
			if err := c.trackerService.Track(shortPath); err != nil {
				slog.WarnContext(ec.Request().Context(), "error tracking", "error", err)
			}

			return err
		}
	})
}

func (c *apiController) Shorten(ec echo.Context) error {
	var req ShortenRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := c.urlService.Shorten(req.URL)
	if err != nil {
		if errors.Is(err, core.ErrURLAlreadyExists) { // obfuscate already exists error
			return ec.JSON(http.StatusInternalServerError, map[string]string{"error": core.ErrSystemError.Error()})
		}

		return ec.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ec.JSON(http.StatusOK, ShortenResponse{
		ShortURL: c.baseDomain + "/u/" + result,
	})
}

func (c *apiController) Redirect(ec echo.Context) error {
	var req RedirectRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := ec.Validate(req)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := c.urlService.Expand(req.ShortPath)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// use found status code to prevent browser from caching the redirect
	// so we can track the redirect for analytics
	// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/302
	return ec.Redirect(http.StatusFound, result)
}
