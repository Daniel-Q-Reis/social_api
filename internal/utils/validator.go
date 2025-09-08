// Package utils provides utility functions for the social media API
package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors
func (ve ValidationErrors) Error() string {
	var errs []string
	for _, err := range ve {
		errs = append(errs, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(errs, "; ")
}

// Validator wraps the go-playground/validator to provide validation functionality
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	v := &Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}

	// Register custom validation rules if needed
	// v.validate.RegisterValidation("custom_rule", customRuleFunc)

	return v
}

// Validate validates a struct based on its tags
func (v *Validator) Validate(s interface{}) error {
	err := v.validate.Struct(s)
	if err != nil {
		// Convert validator.ValidationErrors to our ValidationErrors
		var validationErrors ValidationErrors
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: v.getErrorMessage(err),
			})
		}
		return validationErrors
	}
	return nil
}

// getErrorMessage returns a human-readable error message for a validation error
func (v *Validator) getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("This field must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("This field must be at most %s characters long", err.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", err.Tag())
	}
}

// RegisterTagNameFunc registers a function to get the name of a field
func (v *Validator) RegisterTagNameFunc(fn validator.TagNameFunc) {
	v.validate.RegisterTagNameFunc(fn)
}
