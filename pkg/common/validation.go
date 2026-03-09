package common

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"unicode"

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
	ErrNameTooShort         = "must be at least 2 characters long"
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
	ErrURLEmpty             = "cannot be empty"
	ErrURLInvalid           = "must be a valid HTTP or HTTPS URL"
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

// isNamePunct reports whether r is a name-allowed punctuation mark.
func isNamePunct(r rune) bool {
	return r == '.' || r == '-' || r == '\''
}

// ValidateName validates a name, supporting Unicode letters (e.g. German umlauts).
// Allowed characters: Unicode letters, spaces, dots, hyphens, apostrophes.
func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return newValidationError("Name", ErrNameEmpty)
	}

	runes := []rune(name)
	if len(runes) < MinNameLength {
		return newValidationError("Name", ErrNameTooShort)
	}

	// Allow Unicode letters, spaces, dots, hyphens, and apostrophes.
	for _, r := range runes {
		if !unicode.IsLetter(r) && r != ' ' && !isNamePunct(r) {
			return &ValidationError{Field: "Name", Message: ErrNameInvalidChars}
		}
	}

	// Normalize multiple spaces before further checks.
	name = strings.Join(strings.Fields(name), " ")
	runes = []rune(name)

	// Consecutive punctuation marks are not allowed (e.g. "..", "--").
	for i := 0; i < len(runes)-1; i++ {
		if isNamePunct(runes[i]) && isNamePunct(runes[i+1]) {
			return &ValidationError{Field: "Name", Message: ErrNameConsecutivePunct}
		}
	}

	// Names cannot start or end with hyphens or apostrophes (trailing dot is OK, e.g. "Dr.").
	first, last := runes[0], runes[len(runes)-1]
	if first == '-' || first == '\'' || last == '-' || last == '\'' {
		return &ValidationError{Field: "Name", Message: ErrNameInvalidPunctPos}
	}

	return nil
}

// ValidateEmail checks if the email address is valid.
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
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

// ValidateTitle is intentionally permissive — job titles are free-form text.
func ValidateTitle(_ string) error {
	return nil
}

// ValidateURL checks that rawURL is a syntactically valid HTTP or HTTPS URL.
// An empty (or whitespace-only) value returns nil — callers that require a
// value must enforce non-emptiness separately, enabling this to be reused for
// optional fields. Intranet hostnames (e.g. "http://server/path/") are accepted.
func ValidateURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil
	}
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return newValidationError("URL", ErrURLInvalid)
	}
	return nil
}
