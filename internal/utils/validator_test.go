package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate_Success(t *testing.T) {
	validator := NewValidator()

	type TestStruct struct {
		Name      string    `validate:"required"`
		Email     string    `validate:"required,email"`
		Password  string    `validate:"required,min=8"`
		BirthDate time.Time `validate:"required"`
	}

	testData := TestStruct{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		BirthDate: time.Now(),
	}

	err := validator.Validate(testData)
	assert.NoError(t, err)
}

func TestValidator_Validate_Failure(t *testing.T) {
	validator := NewValidator()

	type TestStruct struct {
		Name      string    `validate:"required"`
		Email     string    `validate:"required,email"`
		Password  string    `validate:"required,min=8"`
		BirthDate time.Time `validate:"required"`
	}

	testData := TestStruct{
		Name:      "",          // Required field missing
		Email:     "invalid",   // Invalid email
		Password:  "123",       // Too short
		BirthDate: time.Time{}, // Zero value
	}

	err := validator.Validate(testData)
	assert.Error(t, err)

	// Check that we get validation errors
	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok)
	assert.Len(t, validationErrors, 4)

	// Check specific error messages
	errorMap := make(map[string]string)
	for _, validationError := range validationErrors {
		errorMap[validationError.Field] = validationError.Message
	}

	assert.Contains(t, errorMap, "Name")
	assert.Contains(t, errorMap, "Email")
	assert.Contains(t, errorMap, "Password")
	assert.Contains(t, errorMap, "BirthDate")
}
