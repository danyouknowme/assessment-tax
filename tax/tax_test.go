package tax

import "testing"

func TestCalculateTaxWithTotalIncomeOnly(t *testing.T) {
	testCases := []struct {
		name   string
		input  TaxCalculationRequest
		expect float64
	}{
		{
			name: "Total income 0, should return 0",
			input: TaxCalculationRequest{
				TotalIncome: 0.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 0.0,
		},
		{
			name: "Total income 30,000.0, should return 0",
			input: TaxCalculationRequest{
				TotalIncome: 30000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 0.0,
		},
		{
			name: "Total income 150,000.0, should return 0",
			input: TaxCalculationRequest{
				TotalIncome: 150000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 0.0,
		},
		{
			name: "Total income 150,001.0, should return 0",
			input: TaxCalculationRequest{
				TotalIncome: 150001.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 0,
		},
		{
			name: "Calculate tax with total income 500000.0",
			input: TaxCalculationRequest{
				TotalIncome: 500000.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0.00},
				},
			},
			expect: 29000.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Calculate(tc.input)

			if got != tc.expect {
				t.Errorf("Expected %v, got %v", tc.expect, got)
			}
		})
	}
}
