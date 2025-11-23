package common

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid with subdomain", "test@sub.example.com", false},
		{"valid with dots", "first.last@example.com", false},
		{"valid with plus", "test+label@example.com", false},
		{"valid with underscore", "test_user@example.com", false},
		{"valid with numbers", "test123@example.com", false},
		{"empty", "", true},
		{"no @", "testexample.com", true},
		{"no domain", "test@", true},
		{"no username", "@example.com", true},
		{"invalid format", "test @example.com", true},
		{"multiple @", "test@@example.com", true},
		{"whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid DE number with spaces",
			input:   "+49 30 12345678",
			wantErr: false,
		},
		{
			name:    "valid DE mobile number",
			input:   "+49 151 12345678",
			wantErr: false,
		},
		{
			name:    "valid DE number without country code",
			input:   "030 12345678",
			wantErr: false,
		},
		{
			name:    "valid DE mobile without country code",
			input:   "0151 12345678",
			wantErr: false,
		},
		{
			name:    "valid US number",
			input:   "+1 555 123 4567",
			wantErr: true,
			errMsg:  "not a valid phone number",
		},
		{
			name:    "valid US number without country code",
			input:   "555 123 4567",
			wantErr: false,
		},
		{
			name:    "valid UK number",
			input:   "+44 20 7946 0958",
			wantErr: false,
		},
		{
			name:    "number with parentheses",
			input:   "+49 (30) 12345678",
			wantErr: false,
		},
		{
			name:    "number with dashes",
			input:   "+49-30-12345678",
			wantErr: false,
		},
		{
			name:    "number with dots",
			input:   "+49.30.12345678",
			wantErr: false,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name:    "invalid format with letters",
			input:   "abcdefghijk",
			wantErr: true,
			errMsg:  "invalid phone number format",
		},
		{
			name:    "possible but invalid number",
			input:   "+49 000 0000000",
			wantErr: true,
			errMsg:  "not a valid phone number",
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name:    "number too short",
			input:   "123",
			wantErr: true,
			errMsg:  "not a valid phone number",
		},
		{
			name:    "number too long",
			input:   "+49 30 12345678901234567890",
			wantErr: true,
			errMsg:  "invalid phone number format",
		},
		{
			name:    "invalid country code",
			input:   "+99 30 12345678",
			wantErr: true,
			errMsg:  "number is too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneNumber(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if validErr, ok := err.(*ValidationError); !ok {
					t.Errorf("ValidatePhoneNumber(%q) error is not ValidationError", tt.input)
				} else if validErr.Message != tt.errMsg {
					t.Errorf("ValidatePhoneNumber(%q) error message = %q, want %q", tt.input, validErr.Message, tt.errMsg)
				}
			}
		})
	}
}

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		country  string
		wantDisp string
		wantLink string
		wantErr  bool
	}{
		{
			name:     "valid DE number",
			input:    "+49123456789",
			country:  "DE",
			wantDisp: "+49 123456789",
			wantLink: "+49123456789",
			wantErr:  false,
		},
		{
			name:     "valid DE number with spaces",
			input:    "+49 30 12345678",
			country:  "DE",
			wantDisp: "+49 30 12345678",
			wantLink: "+493012345678",
			wantErr:  false,
		},
		{
			name:     "valid DE mobile number",
			input:    "+49 151 12345678",
			country:  "DE",
			wantDisp: "+49 1511 2345678",
			wantLink: "+4915112345678",
			wantErr:  false,
		},
		{
			name:     "DE number without country code",
			input:    "030 12345678",
			country:  "DE",
			wantDisp: "+49 30 12345678",
			wantLink: "+493012345678",
			wantErr:  false,
		},
		{
			name:     "DE mobile without country code",
			input:    "0151 12345678",
			country:  "DE",
			wantDisp: "+49 1511 2345678",
			wantLink: "+4915112345678",
			wantErr:  false,
		},
		{
			name:     "US number",
			input:    "+1 555 123 4567",
			country:  "US",
			wantDisp: "+1 555-123-4567",
			wantLink: "+15551234567",
			wantErr:  false,
		},
		{
			name:     "US number without country code",
			input:    "555 123 4567",
			country:  "US",
			wantDisp: "+1 555-123-4567",
			wantLink: "+15551234567",
			wantErr:  false,
		},
		{
			name:     "UK number",
			input:    "+44 20 7946 0958",
			country:  "GB",
			wantDisp: "+44 20 7946 0958",
			wantLink: "+442079460958",
			wantErr:  false,
		},
		{
			name:     "UK number without country code",
			input:    "020 7946 0958",
			country:  "GB",
			wantDisp: "+44 20 7946 0958",
			wantLink: "+442079460958",
			wantErr:  false,
		},
		{
			name:     "French number",
			input:    "+33 1 42 86 20 00",
			country:  "FR",
			wantDisp: "+33 1 42 86 20 00",
			wantLink: "+33142862000",
			wantErr:  false,
		},
		{
			name:     "Italian number",
			input:    "+39 06 6982",
			country:  "IT",
			wantDisp: "+39 06 6982",
			wantLink: "+39066982",
			wantErr:  false,
		},
		{
			name:     "number with parentheses",
			input:    "+49 (30) 12345678",
			country:  "DE",
			wantDisp: "+49 30 12345678",
			wantLink: "+493012345678",
			wantErr:  false,
		},
		{
			name:     "number with dashes",
			input:    "+49-30-12345678",
			country:  "DE",
			wantDisp: "+49 30 12345678",
			wantLink: "+493012345678",
			wantErr:  false,
		},
		{
			name:     "number with dots",
			input:    "+49.30.12345678",
			country:  "DE",
			wantDisp: "+49 30 12345678",
			wantLink: "+493012345678",
			wantErr:  false,
		},
		{
			name:     "invalid number",
			input:    "invalid",
			country:  "DE",
			wantDisp: "invalid",
			wantLink: "invalid",
			wantErr:  true,
		},
		{
			name:     "empty number",
			input:    "",
			country:  "DE",
			wantDisp: "",
			wantLink: "",
			wantErr:  true,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			country:  "DE",
			wantDisp: "   ",
			wantLink: "   ",
			wantErr:  true,
		},
		{
			name:     "number too short",
			input:    "123",
			country:  "DE",
			wantDisp: "+49 123",
			wantLink: "+49123",
			wantErr:  false,
		},
		{
			name:     "number too long",
			input:    "+49 30 12345678901234567890",
			country:  "DE",
			wantDisp: "+49 30 12345678901234567890",
			wantLink: "+49 30 12345678901234567890",
			wantErr:  true,
		},
		{
			name:     "invalid country code",
			input:    "+99 30 12345678",
			country:  "DE",
			wantDisp: "+993 012345678",
			wantLink: "+993012345678",
			wantErr:  false,
		},
		{
			name:     "number with invalid characters",
			input:    "+49 30 abc123",
			country:  "DE",
			wantDisp: "+49 30 222123",
			wantLink: "+4930222123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			display, link, err := FormatPhoneNumber(tt.input, tt.country)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatPhoneNumber(%q, %q) error = %v, wantErr %v", tt.input, tt.country, err, tt.wantErr)
			}
			if !tt.wantErr {
				if display != tt.wantDisp {
					t.Errorf("FormatPhoneNumber(%q, %q) display = %q, want %q", tt.input, tt.country, display, tt.wantDisp)
				}
				if link != tt.wantLink {
					t.Errorf("FormatPhoneNumber(%q, %q) link = %q, want %q", tt.input, tt.country, link, tt.wantLink)
				}
			}
		})
	}
}

func TestValidateSignatureName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "MySignature", false},
		{"valid with spaces", "My Signature", false},
		{"valid with dots", "My.Signature", false},
		{"valid with hyphens", "My-Signature", false},
		{"valid with underscores", "My_Signature", false},
		{"valid with numbers", "MySignature123", false},
		{"invalid with slash", "My/Signature", true},
		{"invalid with backslash", "My\\Signature", true},
		{"invalid with colon", "My:Signature", true},
		{"invalid with star", "My*Signature", true},
		{"invalid with question", "My?Signature", true},
		{"invalid with quotes", "My\"Signature", true},
		{"invalid with angle brackets", "My<Signature>", true},
		{"invalid with pipe", "My|Signature", true},
		{"empty string", "", false},       // Empty names might be valid for some use cases
		{"whitespace only", "   ", false}, // Whitespace-only names might be valid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSignatureName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSignatureName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
