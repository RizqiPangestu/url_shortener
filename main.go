package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RizqiPangestu/url_shortener/internal/adapter"
	"github.com/RizqiPangestu/url_shortener/internal/app"
	"github.com/RizqiPangestu/url_shortener/internal/core"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
	ec := echo.New()
	ec.Validator = a.Validator
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

func registerDependencies(c *dig.Container) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	c.Provide(app.NewAPIController)
	c.Provide(core.NewURLService)
	c.Provide(adapter.NewURLMongoAdapter)
	c.Provide(NewValidator)
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
