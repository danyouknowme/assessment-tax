package api

import (
	"net/http"

	"github.com/danyouknowme/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func CalculateTax(c echo.Context) error {
	var req tax.TaxCalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	taxVal := tax.Calculate(req.TotalIncome, req.Wht, req.Allowances)
	taxLevels := tax.GetTaxLevels(req.TotalIncome, req.Wht, req.Allowances)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tax":      taxVal,
		"taxLevel": taxLevels,
	})
}
