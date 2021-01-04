package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func MainHandler(c echo.Context) error {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("API working from %s, process %d", h, os.Getpid()))
}

func HealthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Health check working")
}
