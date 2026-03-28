package main

import (
	"errors"
	"os"
	"strings"

	"github.com/RizqiPangestu/url_shortener/internal/adapter"
	"github.com/RizqiPangestu/url_shortener/internal/app"
	"github.com/RizqiPangestu/url_shortener/internal/core"
	"github.com/go-pg/pg/v10"
	"go.uber.org/dig"
)

func registerDependencies(c *dig.Container) {
	c.Provide(newAPIController)
	c.Provide(core.NewURLService)
	c.Provide(core.NewTrackerService)
	c.Provide(adapter.NewURLPostgresAdapter)
	c.Provide(adapter.NewTrackerPostgresAdapter)
	c.Provide(NewValidator)
	c.Provide(newPostgresDB)
}

func newAPIController(urlService core.URLService, trackerService core.TrackerService) (app.APIController, error) {
	baseDomain := os.Getenv(ConfigBaseDomain)
	if !strings.HasPrefix(baseDomain, "http") {
		return nil, errors.New("BASE_DOMAIN is not set")
	}
	return app.NewAPIController(urlService, baseDomain, trackerService), nil
}

func newPostgresDB() *pg.DB {
	return adapter.NewPostgresDB(
		os.Getenv(ConfigPostgresHost),
		os.Getenv(ConfigPostgresDatabase),
		os.Getenv(ConfigPostgresUser),
		os.Getenv(ConfigPostgresPassword),
	)
}
