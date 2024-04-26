package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) SettingPersonalDeduction(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "personal deduction set successfully",
	})
}
