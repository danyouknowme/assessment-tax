package tax

import "testing"

func TestCalculateTax(t *testing.T) {
	testCases := []struct {
		name   string
		input  TaxCalculationRequest
		expect float64
	}{
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
