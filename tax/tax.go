package tax

import (
	"math"

	"github.com/danyouknowme/assessment-tax/db"
)

type CalculationRequest struct {
	TotalIncome float64     `json:"totalIncome" validate:"required,min=0.0"`
	Wht         float64     `json:"wht" validate:"wht_custom_validation"`
	Allowances  []Allowance `json:"allowances" validate:"dive"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType" validate:"allowance_type_custom_validation"`
	Amount        float64 `json:"amount" validate:"min=0.0"`
}

func Calculate(defaultDeductions []db.Deduction, req CalculationRequest) (float64, float64) {
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

	if tax < 0 {
		return 0, formatCalculatedTax(tax * -1)
	}

	return formatCalculatedTax(tax), 0
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

func GetTaxLevels(defaultDeductions []db.Deduction, req CalculationRequest) []TaxLevel {
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
