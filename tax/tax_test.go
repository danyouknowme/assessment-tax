package tax

import (
	"reflect"
	"testing"

	"github.com/danyouknowme/assessment-tax/db"
)

type testCase struct {
	name         string
	input        CalculationRequest
	expectTax    float64
	expectRefund float64
}

var defaultDeductions = []db.Deduction{
	{Type: "personal", Amount: 60000.0},
	{Type: "donation", Amount: 100000.0},
	{Type: "k-receipt", Amount: 50000.0},
}

func TestCalculateTaxWithTotalIncomeOnly(t *testing.T) {
	testCases := []testCase{
		{
			name: "Total income 0, should return 0",
			input: CalculationRequest{
				TotalIncome: 0.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 0.0,
		},
		{
			name: "Total income 30,000.0, should return 0",
			input: CalculationRequest{
				TotalIncome: 30000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 0.0,
		},
		{
			name: "Total income 150,000.0, should return 0",
			input: CalculationRequest{
				TotalIncome: 150000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 0.0,
		},
		{
			name: "Total income 150,001.0, should return 0",
			input: CalculationRequest{
				TotalIncome: 150001.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 0.0,
		},
		{
			name: "Total income 500,000 should return 29,000",
			input: CalculationRequest{
				TotalIncome: 500000.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0.00},
				},
			},
			expectTax: 29000.0,
		},
		{
			name: "Total income 500,001 should return 29,000.1",
			input: CalculationRequest{
				TotalIncome: 500001.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0.00},
				},
			},
			expectTax: 29000.1,
		},
		{
			name: "Total income 1,000,000 should return 101,000",
			input: CalculationRequest{
				TotalIncome: 1000000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 101000.0,
		},
		{
			name: "Total income 1,000,001 should return 101,000.15",
			input: CalculationRequest{
				TotalIncome: 1000001.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 101000.15,
		},
		{
			name: "Total income 2,000,000 should return 298,000",
			input: CalculationRequest{
				TotalIncome: 2000000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 298000.0,
		},
		{
			name: "Total income 2,000,001 should return 298,000.2",
			input: CalculationRequest{
				TotalIncome: 2000001.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 298000.2,
		},
		{
			name: "Total income 4,000,000 should return 989,000",
			input: CalculationRequest{
				TotalIncome: 4000000.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 989000.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tax, refund := Calculate(defaultDeductions, tc.input)

			if tax != tc.expectTax {
				t.Errorf("Expected %v, got %v", tc.expectTax, tax)
			}

			if refund != 0 {
				t.Errorf("Expected 0, got %v", refund)
			}
		})
	}
}

func TestCalculateTaxWithWht(t *testing.T) {
	testCases := []testCase{
		{
			name: "Total income 0 and WHT 0 should return 0",
			input: CalculationRequest{
				TotalIncome: 0.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 0.0,
		},
		{
			name: "Total income 500,000.0 and WHT 0.0 should return 29,000",
			input: CalculationRequest{
				TotalIncome: 500000.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0.00},
				},
			},
			expectTax: 29000.0,
		},
		{
			name: "Total income 500,000.0 and WHT 25,000.0 should return 4,000",
			input: CalculationRequest{
				TotalIncome: 500000.0,
				Wht:         25000.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0.00},
				},
			},
			expectTax: 4000.0,
		},
		{
			name: "Total income 500,000.0 and WHT 25,000.0 should return 4,000",
			input: CalculationRequest{
				TotalIncome: 500000.0,
				Wht:         29000.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0.00},
				},
			},
			expectTax: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tax, refund := Calculate(defaultDeductions, tc.input)

			if tax != tc.expectTax {
				t.Errorf("Expected %v, got %v", tc.expectTax, tax)
			}

			if refund != 0 {
				t.Errorf("Expected 0, got %v", refund)
			}
		})
	}
}

func TestCalculateTaxWithDonationAllowances(t *testing.T) {
	testCases := []testCase{
		{
			name: "Total income 0, WHT 0.0 and no allowances, should return 0",
			input: CalculationRequest{
				TotalIncome: 0.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expectTax: 0.0,
		},
		{
			name: "Total income 500,000.0 and WHT 0.0 and donation allowance 20,000 should return 39,000",
			input: CalculationRequest{
				TotalIncome: 500000.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 10000.00},
					{AllowanceType: "donation", Amount: 20000.00},
				},
			},
			expectTax: 26000.0,
		},
		{
			name: "Total income 500,000.0 and WHT 0.0 and donation allowance 200,000 should return 19,000",
			input: CalculationRequest{
				TotalIncome: 500000.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 200000.00},
				},
			},
			expectTax: 19000.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tax, refund := Calculate(defaultDeductions, tc.input)

			if tax != tc.expectTax {
				t.Errorf("Expected %v, got %v", tc.expectTax, tax)
			}

			if refund != 0 {
				t.Errorf("Expected 0, got %v", refund)
			}
		})
	}
}

func TestCalculateTaxAndRefund(t *testing.T) {
	testCases := []testCase{
		{
			name: "Total income 500,000.0 and WHT 100,000.0 and donation allowance 200,000 should return refund 81,000",
			input: CalculationRequest{
				TotalIncome: 500000.0,
				Wht:         100000.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 200000.00},
				},
			},
			expectTax:    0.0,
			expectRefund: 81000.0,
		},
		{
			name: "Total income 1,000,000.0 and WHT 200,000.0 and donation allowance 150,000 should return refund 114,000",
			input: CalculationRequest{
				TotalIncome: 1000000.0,
				Wht:         200000.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 150000.00},
				},
			},
			expectTax:    0.0,
			expectRefund: 114000.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tax, refund := Calculate(defaultDeductions, tc.input)

			if tax != tc.expectTax {
				t.Errorf("Expected %v, got %v", tc.expectTax, tax)
			}

			if refund != tc.expectRefund {
				t.Errorf("Expected %v, got %v", tc.expectRefund, refund)
			}
		})
	}
}

func TestGetTaxLevels(t *testing.T) {
	testCases := []struct {
		name   string
		input  CalculationRequest
		expect []TaxLevel
	}{
		{
			name: "Total income 0, should return 0 of all levels",
			input: CalculationRequest{
				TotalIncome: 0.0,
				Wht:         0.0,
				Allowances:  []Allowance{},
			},
			expect: []TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			},
		},
		{
			name: "Total income 500,000 and donation allowance 200,000, should return 19,000 in level 150,000-500,000",
			input: CalculationRequest{
				TotalIncome: 500000.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 200000.00},
				},
			},
			expect: []TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 19000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetTaxLevels(defaultDeductions, tc.input)

			if !reflect.DeepEqual(got, tc.expect) {
				t.Errorf("Expected %v, got %v", tc.expect, got)
			}
		})
	}
}
