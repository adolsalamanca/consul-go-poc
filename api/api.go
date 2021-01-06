package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func MainHandler(c echo.Context, h string, pid int) error {
	return c.String(http.StatusOK, fmt.Sprintf("API working from %s, process %d", h, pid))
}

func HealthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Health check working")
}
