package common

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/nyaruka/phonenumbers"
)

// ValidationError represents a validation error with a user-friendly message
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateName performs comprehensive validation on a name string
func ValidateName(name string) error {
	// Trim whitespace only at the beginning and end
	name = strings.TrimSpace(name)

	// Check if empty
	if name == "" {
		return &ValidationError{
			Field:   "Name",
			Message: "cannot be empty",
		}
	}

	// Check if name is too short (less than 2 characters)
	if len(name) < 2 {
		return &ValidationError{
			Field:   "Name",
			Message: "must be at least 2 characters long",
		}
	}

	// Use govalidator to check if the name is valid
	if !govalidator.Matches(name, "^[a-zA-Z\\s\\.\\-']+$") {
		return &ValidationError{
			Field:   "Name",
			Message: "can only contain letters, spaces, dots, hyphens, and apostrophes",
		}
	}

	// Normalize multiple spaces into a single space
	name = strings.Join(strings.Fields(name), " ")

	// Check for multiple consecutive punctuation or special characters
	for i := range name[:len(name)-1] {
		current := name[i]
		next := name[i+1]
		if (current == '.' || current == '-' || current == '\'') && (next == '.' || next == '-' || next == '\'') {
			return &ValidationError{
				Field:   "Name",
				Message: "cannot contain multiple consecutive punctuation marks",
			}
		}
	}

	// Check for punctuation at the start or end
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "-") || strings.HasPrefix(name, "'") ||
		strings.HasSuffix(name, "-") || strings.HasSuffix(name, "'") {
		return &ValidationError{
			Field:   "Name",
			Message: "cannot start or end with punctuation marks (except for dots)",
		}
	}

	return nil
}

// ValidateEmail checks if the email address is valid
func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{
			Field:   "Email",
			Message: "cannot be empty",
		}
	}

	// Check for @ symbol
	atIndex := strings.Index(email, "@")
	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return &ValidationError{
			Field:   "Email",
			Message: "must contain a valid domain (e.g., example.com)",
		}
	}

	// Check for domain
	domain := email[atIndex+1:]
	if !strings.Contains(domain, ".") {
		return &ValidationError{
			Field:   "Email",
			Message: "must contain a valid domain (e.g., example.com)",
		}
	}

	// Check for valid characters
	if strings.ContainsAny(email, " ") {
		return &ValidationError{
			Field:   "Email",
			Message: "cannot contain spaces",
		}
	}

	return nil
}

// ValidatePhoneNumber checks if the phone number is valid
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return &ValidationError{
			Field:   "Phone",
			Message: "cannot be empty",
		}
	}

	if len(phone) < 5 {
		return &ValidationError{
			Field:   "Phone",
			Message: "must be at least 5 digits long",
		}
	}

	// Try to parse the phone number
	_, err := phonenumbers.Parse(phone, "DE")
	if err != nil {
		return &ValidationError{
			Field:   "Phone",
			Message: "must be a valid phone number (e.g., +49 123 456789)",
		}
	}

	return nil
}

// FormatPhoneNumber formats a phone number for display and link
func FormatPhoneNumber(phone string, countryCode string) (string, string, error) {
	num, err := phonenumbers.Parse(phone, countryCode)
	if err != nil {
		return phone, phone, &ValidationError{
			Field:   "Phone",
			Message: "could not be formatted",
		}
	}

	display := phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
	link := phonenumbers.Format(num, phonenumbers.E164)

	return display, link, nil
}

// ValidateSignatureName checks if the signature name is valid
func ValidateSignatureName(name string) error {
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return &ValidationError{
			Field:   "Template",
			Message: "contains invalid characters",
		}
	}
	return nil
}
