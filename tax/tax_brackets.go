package tax

import "math"

type TaxBracket struct {
	MinTotalIncome float64
	MaxTotalIncome float64
	TaxRate        float64
	TaxLevel       string
}

var taxBrackets = []TaxBracket{
	{MinTotalIncome: 0, MaxTotalIncome: 150000, TaxRate: 0, TaxLevel: "0-150,000"},
	{MinTotalIncome: 150000, MaxTotalIncome: 500000, TaxRate: 0.1, TaxLevel: "150,001-500,000"},
	{MinTotalIncome: 500000, MaxTotalIncome: 1000000, TaxRate: 0.15, TaxLevel: "500,001-1,000,000"},
	{MinTotalIncome: 1000000, MaxTotalIncome: 2000000, TaxRate: 0.2, TaxLevel: "1,000,001-2,000,000"},
	{MinTotalIncome: 2000000, MaxTotalIncome: math.MaxFloat64, TaxRate: 0.35, TaxLevel: "2,000,001 ขึ้นไป"},
}
