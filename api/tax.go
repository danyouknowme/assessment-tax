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

	totalIncome := req.TotalIncome

	taxVal := tax.Calculate(totalIncome, req.Wht)

	return c.JSON(http.StatusOK, map[string]float64{
		"tax": taxVal,
	})
}
