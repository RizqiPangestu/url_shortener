package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RizqiPangestu/url_shortener/internal/adapter"
	"github.com/RizqiPangestu/url_shortener/internal/app"
	"github.com/RizqiPangestu/url_shortener/internal/core"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

func main() {
	c := dig.New()
	registerDependencies(c)
	if err := c.Invoke(startServer); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

type application struct {
	dig.In
	APIController app.APIController
	Validator     echo.Validator
}

func startServer(a application) {
	configureLogger()
	ec := echo.New()
	ec.Validator = a.Validator
	ec.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: middleware.DefaultSkipper,
		Handler: bodyDumpHandler,
	}))
	ec.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	a.APIController.RegisterRoutes(ec)

	ec.StartServer(&http.Server{
		Addr:         os.Getenv("ADDRESS"),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

func configureLogger() {
	logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				a.Value = slog.StringValue(fmt.Sprintf("%s: %s:%d", source.Function, source.File, source.Line))
			}
			return a
		},
	})
	slog.SetDefault(slog.New(logger))
}

func registerDependencies(c *dig.Container) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	c.Provide(newAPIController)
	c.Provide(core.NewURLService)
	c.Provide(core.NewTrackerService)
	c.Provide(adapter.NewURLPostgresAdapter)
	c.Provide(adapter.NewTrackerPostgresAdapter)
	c.Provide(NewValidator)
	c.Provide(adapter.NewPostgresDB)
}

func newAPIController(urlService core.URLService, trackerService core.TrackerService) (app.APIController, error) {
	baseDomain := os.Getenv("BASE_DOMAIN")
	if !strings.HasPrefix(baseDomain, "http") {
		return nil, errors.New("BASE_DOMAIN is not set")
	}
	return app.NewAPIController(urlService, baseDomain, trackerService), nil
}

type validate struct {
	instance *validator.Validate
}

func NewValidator() echo.Validator {
	return &validate{instance: validator.New(validator.WithRequiredStructEnabled())}
}

func (v *validate) Validate(object interface{}) error {
	if err := v.instance.Struct(object); err != nil {
		var errMsgs []string
		for _, errValidation := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("`%s`", errValidation.Field()))
		}
		return err
	}

	return nil
}

func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	req := c.Request()
	ctx := req.Context()
	slog.InfoContext(ctx, "HTTP request received",
		slog.Any("request", struct {
			Url    string      `json:"url"`
			Method string      `json:"method"`
			Header http.Header `json:"header"`
			Body   string      `json:"body,omitempty"`
		}{
			Url:    req.URL.String(),
			Method: req.Method,
			Header: req.Header,
			Body:   string(reqBody),
		}),
		slog.Any("response", struct {
			Header http.Header `json:"header"`
			Body   string      `json:"body,omitempty"`
			Status int         `json:"status"`
		}{
			Header: c.Response().Header(),
			Body:   string(resBody),
			Status: c.Response().Status,
		}),
	)
}
