package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/danyouknowme/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

type CalculateTaxResponse struct {
	Tax       float64        `json:"tax"`
	TaxRefund float64        `json:"taxRefund"`
	TaxLevel  []tax.TaxLevel `json:"taxLevel"`
}

func (s *Server) CalculateTax(c echo.Context) error {
	var req tax.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	defaultDeductions, err := s.store.GetAllDeductions(c.Request().Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := errors.New("invalid deduction type not found")
			return c.JSON(http.StatusNotFound, errorResponse(err))
		}

		err := errors.New("failed to get deductions")
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	taxVal, taxRefund := tax.Calculate(defaultDeductions, req)
	taxLevels := tax.GetTaxLevels(defaultDeductions, req)

	if taxRefund > 0 {
		return c.JSON(http.StatusOK, CalculateTaxResponse{
			Tax:       0,
			TaxRefund: taxRefund,
		})
	}

	return c.JSON(http.StatusOK, CalculateTaxResponse{
		Tax:      taxVal,
		TaxLevel: taxLevels,
	})
}
