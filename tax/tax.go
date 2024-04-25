package tax

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

func Calculate(req TaxCalculationRequest) float64 {
	income := req.TotalIncome - PersonalAllowance
	if income <= 150000 {
		return 0
	} else if income <= 500000 {
		return (income - 150000) * 0.1
	}

	return 0
}
