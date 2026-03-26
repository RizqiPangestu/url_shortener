package app

import (
	"errors"
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
	URLService core.URLService
}

func NewAPIController(urlService core.URLService) APIController {
	return &apiController{
		URLService: urlService,
	}
}

func (c *apiController) RegisterRoutes(ec *echo.Echo) {
	ec.POST("/shorten", c.Shorten)
	ec.GET("u/:short_path", c.Redirect)
}

func (c *apiController) Shorten(ec echo.Context) error {
	var req ShortenRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := c.URLService.Shorten(req.URL)
	if err != nil {
		if errors.Is(err, core.ErrURLAlreadyExists) { // obfuscate already exists error
			return ec.JSON(http.StatusInternalServerError, map[string]string{"error": core.ErrSystemError.Error()})
		}

		return ec.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ec.JSON(http.StatusOK, ShortenResponse{
		ShortURL: result,
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

	result, err := c.URLService.Expand(req.ShortPath)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ec.Redirect(http.StatusPermanentRedirect, result)
}
