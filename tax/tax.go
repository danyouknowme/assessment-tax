package tax

import (
	"math"
)

type TaxCalculationRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	Wht         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

const PersonalAllowance = 60000.0

type TaxBracket struct {
	MinTotalIncome float64
	MaxTotalIncome float64
	TaxRate        float64
}

var taxBrackets = []TaxBracket{
	{MinTotalIncome: 0, MaxTotalIncome: 150000, TaxRate: 0},
	{MinTotalIncome: 150000, MaxTotalIncome: 500000, TaxRate: 0.1},
	{MinTotalIncome: 500000, MaxTotalIncome: 1000000, TaxRate: 0.15},
	{MinTotalIncome: 1000000, MaxTotalIncome: 2000000, TaxRate: 0.2},
	{MinTotalIncome: 2000000, MaxTotalIncome: math.MaxFloat64, TaxRate: 0.35},
}

func Calculate(netIncome float64) float64 {
	var tax float64 = 0
	taxableIncome := netIncome - PersonalAllowance
	for _, bracket := range taxBrackets {
		if taxableIncome <= 0 {
			break
		}

		bracketRange := bracket.MaxTotalIncome - bracket.MinTotalIncome
		incomeInBracket := math.Min(taxableIncome, bracketRange)
		taxInBracket := incomeInBracket * bracket.TaxRate
		tax += taxInBracket
		taxableIncome -= incomeInBracket
	}
	return formatCalculatedTax(tax)
}

func formatCalculatedTax(tax float64) float64 {
	return math.Round(tax*100) / 100
}
