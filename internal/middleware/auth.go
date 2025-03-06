package imiddleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"testProject/internal/models"
)

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !c.IsTLS() {
			return echo.NewHTTPError(http.StatusBadRequest, "Connection not secured")
		}
		if c.Request().Body == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		var userRequest models.LoginRequest
		if err := c.Bind(&userRequest); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		c.Set("user", userRequest)
		return next(c)
	}
}
