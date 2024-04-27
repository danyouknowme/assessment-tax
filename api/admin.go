package api

import (
	"database/sql"
	"errors"
	"github.com/danyouknowme/assessment-tax/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SettingPersonalDeductionRequest struct {
	Amount float64 `json:"amount" validate:"required,min=10000.0,max=100000.0"`
}

func (s *Server) SettingPersonalDeduction(c echo.Context) error {
	var req SettingPersonalDeductionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := c.Validate(req); err != nil {
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
		if errors.Is(err, sql.ErrNoRows) {
			err := errors.New("personal deduction not found")
			return c.JSON(http.StatusNotFound, errorResponse(err))
		}

		err := errors.New("failed to update personal deduction")
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, map[string]float64{
		"personalDeduction": deduction.Amount,
	})
}

type SettingKReceiptDeductionRequest struct {
	Amount float64 `json:"amount" validate:"required,min=0.0,max=100000.0"`
}

func (s *Server) SettingKReceiptDeduction(c echo.Context) error {
	var req SettingKReceiptDeductionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	deduction, err := s.store.UpdateDeductionByType(
		c.Request().Context(),
		"k-receipt",
		db.UpdateDeductionParams{
			Amount: req.Amount,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := errors.New("k-receipt deduction not found")
			return c.JSON(http.StatusNotFound, errorResponse(err))
		}

		err := errors.New("failed to update k-receipt deduction")
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, map[string]float64{
		"kReceipt": deduction.Amount,
	})
}
