package utils

import (
	"regexp"
	"unicode"
)

// SanitizeInput removes potentially dangerous characters from input
func SanitizeInput(input string) string {
	// Remove script tags (case insensitive)
	re := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = re.ReplaceAllString(input, "")

	// Remove javascript: and vbscript: protocols
	re = regexp.MustCompile(`(?i)(javascript:|vbscript:)`)
	input = re.ReplaceAllString(input, "")

	// Remove on-event handlers like onclick, onload, etc.
	re = regexp.MustCompile(`(?i)on\w+\s*=`)
	input = re.ReplaceAllString(input, "")

	return input
}

// SanitizeHTML removes dangerous HTML elements and attributes
func SanitizeHTML(html string) string {
	// This is a simplified sanitizer - in production, you'd want to use a more robust library
	// Remove potentially dangerous tags
	tagsToRemove := []string{
		"script", "iframe", "object", "embed", "form", "input", "button", 
		"link", "meta", "base", "frameset", "frame",
	}

	for _, tag := range tagsToRemove {
		// Remove opening tags
		re := regexp.MustCompile(`(?i)<\s*` + tag + `[^>]*>`)
		html = re.ReplaceAllString(html, "")
		
		// Remove closing tags
		re = regexp.MustCompile(`(?i)<\s*/\s*` + tag + `\s*>`)
		html = re.ReplaceAllString(html, "")
	}

	return html
}

// ValidatePasswordStrength validates that a password meets security requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return ValidationError("password must be at least 8 characters long")
	}

	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case !unicode.IsLetter(char) && !unicode.IsNumber(char):
			hasSpecial = true
		}
	}

	if !hasMinLen || !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ValidationError("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}



// ValidateSafeString ensures the string doesn't contain dangerous patterns
func ValidateSafeString(input string, fieldName string) error {
	// Check for SQL injection patterns
	sqlPatterns := []string{
		`(?i)union\s+select`,
		`(?i)drop\s+table`,
		`(?i)drop\s+database`,
		`(?i)exec\s*\(`,
		`(?i)insert\s+into`,
		`(?i)delete\s+from`,
		`(?i)update\s+\w+\s+set`,
		`'\s*(or|and)\s*1\s*=\s*1`,
	}

	for _, pattern := range sqlPatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(input) {
			return ValidationError(fieldName + " contains potential SQL injection patterns")
		}
	}

	// Check for XSS patterns
	xssPatterns := []string{
		`(?i)<script`,
		`(?i)javascript:`,
		`(?i)vbscript:`,
		`(?i)on\w+\s*=`,
	}

	for _, pattern := range xssPatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(input) {
			return ValidationError(fieldName + " contains potential XSS patterns")
		}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

// ValidateCSRFToken validates a CSRF token (in a real implementation, you'd verify the token against a stored value)
func ValidateCSRFToken(token, expectedToken string) bool {
	// In a real implementation, this would verify the token against a stored value
	return token == expectedToken
}