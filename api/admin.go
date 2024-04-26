package api

import (
	"github.com/danyouknowme/assessment-tax/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SettingPersonalDeductionRequest struct {
	Amount float64 `json:"amount"`
}

func (s *Server) SettingPersonalDeduction(c echo.Context) error {
	var req SettingPersonalDeductionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	deduction, err := s.store.UpdateDeductionByType(
		c.Request().Context(),
		"personal",
		db.UpdateDeductionParams{
			Amount: req.Amount,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, map[string]float64{
		"personalDeduction": deduction.Amount,
	})
}
