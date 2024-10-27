package utils

import "github.com/go-playground/validator/v10"

// XValidator is a struct that wraps the validator from the go-playground/validator/v10 package.
type XValidator struct {
	validator *validator.Validate
}

// InitValidator init validator
func InitValidator(v *validator.Validate) IXValidator {
	return &XValidator{validator: v}
}

// IXValidator is an interface defining validation operations.
type IXValidator interface {
	Validate(data interface{}) []ErrorResponse
}

// ErrorResponse represents details of a validation error.
type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

// Validate performs validation on the provided data using the wrapped validator.
func (v *XValidator) Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}
