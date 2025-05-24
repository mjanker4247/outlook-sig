package common

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// ValidateEmail checks if the email address is valid
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// Check for @ symbol
	atIndex := strings.Index(email, "@")
	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return fmt.Errorf("invalid email format")
	}

	// Check for domain
	domain := email[atIndex+1:]
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("invalid email format")
	}

	// Check for valid characters
	if strings.ContainsAny(email, " ") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidatePhoneNumber checks if the phone number is valid
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}
	if len(phone) < 5 {
		return fmt.Errorf("phone number is too short")
	}
	return nil
}

// FormatPhoneNumber formats a phone number for display and link
func FormatPhoneNumber(phone string, countryCode string) (string, string, error) {
	num, err := phonenumbers.Parse(phone, countryCode)
	if err != nil {
		return phone, phone, err
	}

	display := phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
	link := phonenumbers.Format(num, phonenumbers.E164)

	return display, link, nil
}

// ValidateSignatureName checks if the signature name is valid
func ValidateSignatureName(name string) error {
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return fmt.Errorf("invalid signature name: contains invalid characters")
	}
	return nil
}
