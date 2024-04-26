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

const (
	PersonalDeduction    = 60000.0
	MaxDonationAllowance = 100000.0
)

func Calculate(totalIncome, wht float64, allowances []Allowance) float64 {
	var tax float64 = 0

	donationAllowance := calculateDonationAllowance(allowances)
	taxableIncome := calculateTaxableIncome(totalIncome, donationAllowance)

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

	tax -= wht

	return formatCalculatedTax(tax)
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

func GetTaxLevels(totalIncome, wht float64, allowances []Allowance) []TaxLevel {
	var taxLevels []TaxLevel

	donationAllowance := calculateDonationAllowance(allowances)
	taxableIncome := calculateTaxableIncome(totalIncome, donationAllowance)

	for _, bracket := range taxBrackets {
		bracketRange := bracket.MaxTotalIncome - bracket.MinTotalIncome
		incomeInBracket := math.Min(taxableIncome, bracketRange)
		taxInBracket := incomeInBracket * bracket.TaxRate
		taxLevels = append(taxLevels, TaxLevel{Level: bracket.TaxLevel, Tax: formatCalculatedTax(taxInBracket)})
		taxableIncome -= incomeInBracket
	}

	return taxLevels
}

func calculateTaxableIncome(totalIncome float64, donationAllowance float64) float64 {
	return totalIncome - PersonalDeduction - donationAllowance
}

func calculateDonationAllowance(allowances []Allowance) float64 {
	var donationAllowance float64 = 0
	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			donationAllowance += allowance.Amount
		}
	}

	if donationAllowance > MaxDonationAllowance {
		return MaxDonationAllowance
	}

	return donationAllowance
}

func formatCalculatedTax(tax float64) float64 {
	return math.Round(tax*100) / 100
}
