package api

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() (*CustomValidator, error) {
	validate := validator.New()

	// Register custom validation
	if err := registerAllowanceTypeValidation(validate); err != nil {
		return nil, err
	}
	if err := registerWhtValidation(validate); err != nil {
		return nil, err
	}

	return &CustomValidator{validator: validate}, nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func registerAllowanceTypeValidation(v *validator.Validate) error {
	return v.RegisterValidation("allowance_type_custom_validation", func(fl validator.FieldLevel) bool {
		allowanceType := fl.Field().String()
		return allowanceType == "donation" || allowanceType == "k-receipt"
	})
}

func registerWhtValidation(v *validator.Validate) error {
	return v.RegisterValidation("wht_custom_validation", func(fl validator.FieldLevel) bool {
		wht := fl.Field().Float()
		totalIncome := fl.Parent().FieldByName("TotalIncome").Float()
		return wht > 0 && wht < totalIncome
	})
}
