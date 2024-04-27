package api

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"strconv"

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

type CalculateTaxForCSVResponse struct {
	Taxes []TaxCSV `json:"taxes"`
}

type TaxCSV struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

func (s *Server) CalculateTaxForCSV(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		err := errors.New("missing file")
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	src, err := file.Open()
	if err != nil {
		err := errors.New("failed to open file")
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	defer src.Close()

	reader := csv.NewReader(src)

	header, err := reader.Read()
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := validateCSVHeader(header); err != nil {
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

	var taxes []TaxCSV
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return c.JSON(http.StatusBadRequest, errorResponse(err))
		}

		req, err := validateCSVBodyRequest(record)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errorResponse(err))
		}

		taxVal, _ := tax.Calculate(defaultDeductions, req)
		taxes = append(taxes, TaxCSV{
			TotalIncome: req.TotalIncome,
			Tax:         taxVal,
		})
	}

	return c.JSON(http.StatusOK, CalculateTaxForCSVResponse{
		Taxes: taxes,
	})
}

func validateCSVHeader(header []string) error {
	if len(header) != 3 {
		return errors.New("invalid csv header")
	}

	if header[0] != "totalIncome" || header[1] != "wht" || header[2] != "donation" {
		return errors.New("invalid csv header")
	}

	return nil
}

func validateCSVBodyRequest(record []string) (tax.CalculationRequest, error) {
	if len(record) != 3 {
		return tax.CalculationRequest{}, errors.New("invalid csv body")
	}

	totalIncome, err := strconv.ParseFloat(record[0], 64)
	if err != nil {
		return tax.CalculationRequest{}, errors.New("invalid total income")
	}

	wht, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return tax.CalculationRequest{}, errors.New("invalid wht")
	}

	donation, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return tax.CalculationRequest{}, errors.New("invalid donation")
	}

	return tax.CalculationRequest{
		TotalIncome: totalIncome,
		Wht:         wht,
		Allowances: []tax.Allowance{
			{
				AllowanceType: "donation",
				Amount:        donation,
			},
		},
	}, nil
}
