package common

import (
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "John Doe", false},
		{"valid with dot", "J. R. R. Tolkien", false},
		{"valid with hyphen", "Jean-Paul", false},
		{"valid with apostrophe", "O'Connor", false},
		{"empty", "", true},
		{"too short", "J", true},
		{"invalid chars", "John123", true},
		{"consecutive punctuation", "John--Doe", true},
		{"starts with punctuation", "-John", true},
		{"ends with punctuation", "John-", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

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
		{"empty", "", true},
		{"no @", "testexample.com", true},
		{"no domain", "test@", true},
		{"no username", "@example.com", true},
		{"invalid format", "test @example.com", true},
		{"multiple @", "test@@example.com", true},
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
			name:     "invalid number",
			input:    "invalid",
			country:  "DE",
			wantDisp: "invalid",
			wantLink: "invalid",
			wantErr:  true,
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
		{"invalid with slash", "My/Signature", true},
		{"invalid with backslash", "My\\Signature", true},
		{"invalid with colon", "My:Signature", true},
		{"invalid with star", "My*Signature", true},
		{"invalid with question", "My?Signature", true},
		{"invalid with quotes", "My\"Signature", true},
		{"invalid with angle brackets", "My<Signature>", true},
		{"invalid with pipe", "My|Signature", true},
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
