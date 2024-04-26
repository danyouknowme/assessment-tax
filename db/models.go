package db

type Deduction struct {
	Type   string
	Amount float64
}

type UpdateDeductionParams struct {
	Amount float64 `json:"amount"`
}
