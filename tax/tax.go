package tax

import (
	"math"

	"github.com/danyouknowme/assessment-tax/db"
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

func Calculate(defaultDeductions []db.Deduction, req TaxCalculationRequest) float64 {
	var tax float64 = 0

	donationAllowance := calculateDonationAllowance(getDeductionByType(defaultDeductions, "donation").Amount, req.Allowances)
	taxableIncome := calculateTaxableIncome(req.TotalIncome, getDeductionByType(defaultDeductions, "personal").Amount, donationAllowance)

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

	tax -= req.Wht

	return formatCalculatedTax(tax)
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

func GetTaxLevels(defaultDeductions []db.Deduction, req TaxCalculationRequest) []TaxLevel {
	var taxLevels []TaxLevel

	donationAllowance := calculateDonationAllowance(getDeductionByType(defaultDeductions, "donation").Amount, req.Allowances)
	taxableIncome := calculateTaxableIncome(req.TotalIncome, getDeductionByType(defaultDeductions, "personal").Amount, donationAllowance)

	for _, bracket := range taxBrackets {
		bracketRange := bracket.MaxTotalIncome - bracket.MinTotalIncome
		incomeInBracket := math.Min(taxableIncome, bracketRange)
		taxInBracket := incomeInBracket * bracket.TaxRate
		taxLevels = append(taxLevels, TaxLevel{Level: bracket.TaxLevel, Tax: formatCalculatedTax(taxInBracket)})
		taxableIncome -= incomeInBracket
	}

	return taxLevels
}

func calculateTaxableIncome(totalIncome, personalDeduction, donationAllowance float64) float64 {
	return totalIncome - personalDeduction - donationAllowance
}

func calculateDonationAllowance(maxDonationDeduction float64, allowances []Allowance) float64 {
	var donationAllowance float64 = 0
	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			donationAllowance += allowance.Amount
		}
	}

	if donationAllowance > maxDonationDeduction {
		return maxDonationDeduction
	}

	return donationAllowance
}

func formatCalculatedTax(tax float64) float64 {
	return math.Round(tax*100) / 100
}

func getDeductionByType(deductions []db.Deduction, deductionType string) db.Deduction {
	for _, deduction := range deductions {
		if deduction.Type == deductionType {
			return deduction
		}
	}

	return db.Deduction{}
}
