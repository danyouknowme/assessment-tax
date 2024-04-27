package api

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
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

func (s *Server) acceptCSVExtension(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("taxFile")
		if err != nil {
			err := errors.New("missing file")
			return c.JSON(http.StatusBadRequest, errorResponse(err))
		}

		fmt.Println("File Header: ", file.Header.Get("Content-Type"))

		if !strings.HasPrefix(file.Header.Get("Content-Type"), "text/csv") {
			err := errors.New("invalid file format")
			return c.JSON(http.StatusBadRequest, errorResponse(err))
		}

		return next(c)
	}
}
