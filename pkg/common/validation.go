package common

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/nyaruka/phonenumbers"
)

const (
	// MinNameLength is the minimum required length for a name
	MinNameLength = 2
	// BufferSizeLimit is the maximum buffer size to prevent memory exhaustion
	BufferSizeLimit = 5 * 1024 * 1024 // 5MB
	// FileSizeLimit is the maximum file size for copying
	FileSizeLimit = 50 * 1024 * 1024 // 50MB
)

// Validation error messages
const (
	ErrNameEmpty            = "cannot be empty"
	ErrNameTooShort         = "must contain at least one non-empty line"
	ErrNameInvalidChars     = "can only contain letters, spaces, dots, hyphens, and apostrophes"
	ErrNameConsecutivePunct = "cannot contain multiple consecutive punctuation marks"
	ErrNameInvalidPunctPos  = "cannot start or end with punctuation marks (except for dots)"
	ErrEmailEmpty           = "cannot be empty"
	ErrEmailInvalid         = "invalid email format"
	ErrPhoneEmpty           = "cannot be empty"
	ErrPhoneInvalid         = "invalid phone number format"
	ErrPhoneNotValid        = "not a valid phone number"
	ErrPhoneCountryCode     = "invalid country code"
	ErrPhoneTooShort        = "number is too short"
	ErrPhoneTooLong         = "number is too long"
	ErrPhoneInvalidLength   = "number has invalid length for the country"
	ErrPhoneFormat          = "could not be formatted"
	ErrSignatureInvalid     = "contains invalid characters"
)

// ValidationError represents a validation error with a user-friendly message
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// newValidationError creates a new validation error with consistent formatting
func newValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// CleanLineBreaks removes multiple consecutive line breaks and normalizes them to single line breaks
func CleanLineBreaks(input string) string {
	// Split by line breaks
	lines := strings.Split(input, "\n")

	// Filter out empty lines and trim whitespace
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	// Join back with single line breaks
	return strings.Join(cleanLines, "\n")
}

// ValidateName performs comprehensive validation on a name string
func ValidateName(name string) error {
	// Trim whitespace only at the beginning and end
	name = strings.TrimSpace(name)

	// Check if empty
	if name == "" {
		return newValidationError("Name", ErrNameEmpty)
	}

	// Split into lines and filter out empty/non-visible lines
	validLines := getValidLines(name)
	if len(validLines) == 0 {
		return newValidationError("Name", ErrNameTooShort)
	}

	// Check if name is too short (less than MinNameLength characters)
	if len(validLines[0]) < MinNameLength {
		return newValidationError("Name", fmt.Sprintf("first line must be at least %d characters long", MinNameLength))
	}

	// Validate each valid line
	for i, line := range validLines {
		if err := validateNameLine(line, i+1); err != nil {
			return err
		}
	}

	return nil
}

// getValidLines splits name into lines and filters out empty ones
func getValidLines(name string) []string {
	var validLines []string
	lines := strings.Split(name, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			validLines = append(validLines, line)
		}
	}
	return validLines
}

// validateNameLine validates a single line of the name
func validateNameLine(line string, lineNumber int) error {
	// Use govalidator to check if each line is valid
	if !govalidator.Matches(line, `^[a-zA-Z\s\.\-']+$`) {
		return newValidationError("Name", fmt.Sprintf("line %d %s", lineNumber, ErrNameInvalidChars))
	}

	// Normalize multiple spaces into a single space
	line = strings.Join(strings.Fields(line), " ")

	// Check for multiple consecutive punctuation or special characters
	if hasConsecutivePunctuation(line) {
		return newValidationError("Name", fmt.Sprintf("line %d %s", lineNumber, ErrNameConsecutivePunct))
	}

	// Check for punctuation at the start or end
	if hasInvalidPunctuationPosition(line) {
		return newValidationError("Name", fmt.Sprintf("line %d %s", lineNumber, ErrNameInvalidPunctPos))
	}

	return nil
}

// hasConsecutivePunctuation checks if a string contains consecutive punctuation marks
func hasConsecutivePunctuation(line string) bool {
	if len(line) < 2 {
		return false
	}
	for j := range line[:len(line)-1] {
		current := line[j]
		next := line[j+1]
		if isPunctuation(current) && isPunctuation(next) {
			return true
		}
	}
	return false
}

// hasInvalidPunctuationPosition checks if a string starts or ends with invalid punctuation
func hasInvalidPunctuationPosition(line string) bool {
	punctuationStart := strings.HasPrefix(line, ".") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "'")
	punctuationEnd := strings.HasSuffix(line, "-") || strings.HasSuffix(line, "'")
	return punctuationStart || punctuationEnd
}

// isPunctuation checks if a character is a punctuation mark
func isPunctuation(char byte) bool {
	return char == '.' || char == '-' || char == '\''
}

// ValidateEmail checks if the email address is valid
func ValidateEmail(email string) error {
	if email == "" {
		return newValidationError("Email", ErrEmailEmpty)
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return newValidationError("Email", ErrEmailInvalid)
	}

	return nil
}

// ValidatePhoneNumber checks if the phone number is valid using the phonenumbers library
func ValidatePhoneNumber(phone string) error {
	// Trim whitespace first
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return newValidationError("Phone", ErrPhoneEmpty)
	}

	// Try to parse the phone number (defaulting to DE as fallback)
	num, err := phonenumbers.Parse(phone, "DE")
	if err != nil {
		return newValidationError("Phone", ErrPhoneInvalid)
	}

	return validatePhoneNumberReason(num)
}

// validatePhoneNumberReason validates the phone number based on the reason code
func validatePhoneNumberReason(num *phonenumbers.PhoneNumber) error {
	reason := phonenumbers.IsPossibleNumberWithReason(num)

	switch reason {
	case phonenumbers.IS_POSSIBLE:
		// Only check IsValidNumber if the number is possible
		if !phonenumbers.IsValidNumber(num) {
			return newValidationError("Phone", ErrPhoneNotValid)
		}
		return nil
	case phonenumbers.INVALID_COUNTRY_CODE:
		return newValidationError("Phone", ErrPhoneCountryCode)
	case phonenumbers.TOO_SHORT:
		return newValidationError("Phone", ErrPhoneTooShort)
	case phonenumbers.TOO_LONG:
		return newValidationError("Phone", ErrPhoneTooLong)
	case phonenumbers.INVALID_LENGTH:
		return newValidationError("Phone", ErrPhoneInvalidLength)
	default:
		return newValidationError("Phone", ErrPhoneNotValid)
	}
}

// FormatPhoneNumber formats a phone number for display and link
func FormatPhoneNumber(phone string, countryCode string) (string, string, error) {
	num, err := phonenumbers.Parse(phone, countryCode)
	if err != nil {
		return phone, phone, newValidationError("Phone", ErrPhoneFormat)
	}

	display := phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
	link := phonenumbers.Format(num, phonenumbers.E164)

	return display, link, nil
}

// ValidateSignatureName checks if the signature name is valid
func ValidateSignatureName(name string) error {
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return newValidationError("Template", ErrSignatureInvalid)
	}
	return nil
}
