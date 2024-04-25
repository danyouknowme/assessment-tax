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
			expect: 0.0,
		},
		{
			name: "Total income 500,000 should return 29,000",
			input: TaxCalculationRequest{
				TotalIncome: 500000.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0.00},
				},
			},
			expect: 29000.0,
		},
		{
			name: "Total income 500,001 should return 29,000.1",
			input: TaxCalculationRequest{
				TotalIncome: 500001.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0.00},
				},
			},
			expect: 29000.1,
		},
		{
			name: "Total income 1,000,000 should return 101,000",
			input: TaxCalculationRequest{
				TotalIncome: 1000000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 101000.0,
		},
		{
			name: "Total income 1,000,001 should return 101,000.15",
			input: TaxCalculationRequest{
				TotalIncome: 1000001.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 101000.15,
		},
		{
			name: "Total income 2,000,000 should return 298,000",
			input: TaxCalculationRequest{
				TotalIncome: 2000000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 298000.0,
		},
		{
			name: "Total income 2,000,001 should return 298,000.2",
			input: TaxCalculationRequest{
				TotalIncome: 2000001.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 298000.2,
		},
		{
			name: "Total income 4,000,000 should return 989,000",
			input: TaxCalculationRequest{
				TotalIncome: 4000000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: 989000.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			totalIncome := tc.input.TotalIncome
			got := Calculate(totalIncome)

			if got != tc.expect {
				t.Errorf("Expected %v, got %v", tc.expect, got)
			}
		})
	}
}
