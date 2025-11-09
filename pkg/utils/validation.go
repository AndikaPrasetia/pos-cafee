package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateRequired validates that a string is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateEmail validates that a string is a valid email format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateMinLength validates that a string has at least the specified length
func ValidateMinLength(value string, minLength int, fieldName string) error {
	if len(value) < minLength {
		return fmt.Errorf("%s must be at least %d characters", fieldName, minLength)
	}
	return nil
}

// ValidateMaxLength validates that a string has at most the specified length
func ValidateMaxLength(value string, maxLength int, fieldName string) error {
	if len(value) > maxLength {
		return fmt.Errorf("%s must be at most %d characters", fieldName, maxLength)
	}
	return nil
}

// ValidateBetweenLength validates that a string's length is between min and max values
func ValidateBetweenLength(value string, minLength, maxLength int, fieldName string) error {
	if len(value) < minLength || len(value) > maxLength {
		return fmt.Errorf("%s must be between %d and %d characters", fieldName, minLength, maxLength)
	}
	return nil
}

// ValidatePositiveNumber validates that a number is positive
func ValidatePositiveNumber(value float64, fieldName string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than 0", fieldName)
	}
	return nil
}

// ValidateNonNegativeNumber validates that a number is non-negative
func ValidateNonNegativeNumber(value float64, fieldName string) error {
	if value < 0 {
		return fmt.Errorf("%s must be greater than or equal to 0", fieldName)
	}
	return nil
}

// ValidateUUID validates that a string is a valid UUID
func ValidateUUID(uuid string, fieldName string) error {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	if !uuidRegex.MatchString(uuid) {
		return fmt.Errorf("%s must be a valid UUID", fieldName)
	}
	return nil
}

// ValidateEnum validates that a string value is one of the allowed values
func ValidateEnum(value string, allowedValues []string, fieldName string) error {
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	
	// Build error message with allowed values
	allowedStr := strings.Join(allowedValues, ", ")
	return fmt.Errorf("%s must be one of: %s", fieldName, allowedStr)
}

// ValidateDecimalPrecision validates decimal precision
func ValidateDecimalPrecision(value string, totalDigits, decimalPlaces int, fieldName string) error {
	// This function would check if the decimal value fits within the specified precision
	// For now, it's a basic implementation
	if len(value) > totalDigits {
		return fmt.Errorf("%s exceeds the maximum allowed length", fieldName)
	}
	
	// Check for decimal places if needed
	if strings.Contains(value, ".") {
		parts := strings.Split(value, ".")
		if len(parts) == 2 && len(parts[1]) > decimalPlaces {
			return fmt.Errorf("%s exceeds the maximum allowed decimal places (%d)", fieldName, decimalPlaces)
		}
	}
	return nil
}