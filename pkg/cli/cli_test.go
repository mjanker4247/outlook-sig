package cli

import (
	"fmt"
	"outlook-signature/pkg/common"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected error
	}{
		{"test@example.com", nil},
		{"test@example.co.uk", nil},
		{"test@example", fmt.Errorf("invalid email format")},
		{"testexample.com", fmt.Errorf("invalid email format")},
		{"@example.com", fmt.Errorf("invalid email format")},
		{"", fmt.Errorf("invalid email format")},
	}

	for _, test := range tests {
		err := common.ValidateEmail(test.email)
		if (err == nil && test.expected != nil) || (err != nil && test.expected == nil) {
			t.Errorf("ValidateEmail(%q) = %v, want %v", test.email, err, test.expected)
		}
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		phone    string
		expected error
	}{
		{"+49123456789", nil},
		{"+49 123 456789", nil},
		{"0123456789", nil},
		{"1234", fmt.Errorf("phone number is too short")},
		{"", fmt.Errorf("phone number cannot be empty")},
	}

	for _, test := range tests {
		err := common.ValidatePhoneNumber(test.phone)
		if (err == nil && test.expected != nil) || (err != nil && test.expected == nil) {
			t.Errorf("ValidatePhoneNumber(%q) = %v, want %v", test.phone, err, test.expected)
		}
	}
}

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		phone       string
		countryCode string
		display     string
		link        string
		shouldError bool
	}{
		{
			phone:       "+49123456789",
			countryCode: "DE",
			display:     "+49 123456789",
			link:        "+49123456789",
			shouldError: false,
		},
		{
			phone:       "0123456789",
			countryCode: "DE",
			display:     "+49 123456789",
			link:        "+49123456789",
			shouldError: false,
		},
		{
			phone:       "invalid",
			countryCode: "DE",
			display:     "invalid",
			link:        "invalid",
			shouldError: true,
		},
	}

	for _, test := range tests {
		display, link, err := common.FormatPhoneNumber(test.phone, test.countryCode)
		if test.shouldError {
			if err == nil {
				t.Errorf("FormatPhoneNumber(%q, %q) should have returned an error", test.phone, test.countryCode)
			}
		} else {
			if err != nil {
				t.Errorf("FormatPhoneNumber(%q, %q) returned error: %v", test.phone, test.countryCode, err)
			}
			if display != test.display {
				t.Errorf("FormatPhoneNumber(%q, %q) display = %q, want %q", test.phone, test.countryCode, display, test.display)
			}
			if link != test.link {
				t.Errorf("FormatPhoneNumber(%q, %q) link = %q, want %q", test.phone, test.countryCode, link, test.link)
			}
		}
	}
}
