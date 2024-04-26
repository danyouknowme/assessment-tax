package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) basicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, p, ok := c.Request().BasicAuth()
		if !ok {
			err := errors.New("missing basic auth")
			return c.JSON(http.StatusUnauthorized, errorResponse(err))
		}

		if u != s.config.AdminUsername || p != s.config.AdminPassword {
			err := errors.New("invalid username or password")
			return c.JSON(http.StatusUnauthorized, errorResponse(err))
		}

		return next(c)
	}
}
