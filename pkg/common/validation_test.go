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
		{"valid multiline", "John Doe\nSoftware Engineer\nSenior Developer", false},
		{"empty", "", true},
		{"too short", "J", true},
		{"invalid chars", "John123", true},
		{"consecutive punctuation", "John--Doe", true},
		{"starts with punctuation", "-John", true},
		{"ends with punctuation", "John-", true},
		{"starts with dot", ".John", true},
		{"ends with dot", "John.", false},                  // Dots are allowed at the end
		{"multiple spaces", "John   Doe", false},           // Multiple spaces are normalized
		{"empty lines", "John\n\nDoe", false},              // Empty lines are filtered
		{"whitespace only lines", "John\n   \nDoe", false}, // Whitespace-only lines are filtered
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

// Test helper functions
func TestHasConsecutivePunctuation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"no consecutive", "John Doe", false},
		{"consecutive dots", "John..Doe", true},
		{"consecutive hyphens", "John--Doe", true},
		{"consecutive apostrophes", "John''Doe", true},
		{"mixed consecutive", "John.-Doe", true},
		{"single punctuation", "John-Doe", false},
		{"empty string", "", false},
		{"single character", "J", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasConsecutivePunctuation(tt.input)
			if result != tt.expected {
				t.Errorf("hasConsecutivePunctuation(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHasInvalidPunctuationPosition(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid start and end", "John Doe", false},
		{"starts with dot", ".John Doe", true},
		{"starts with hyphen", "-John Doe", true},
		{"starts with apostrophe", "'John Doe", true},
		{"ends with hyphen", "John Doe-", true},
		{"ends with apostrophe", "John Doe'", true},
		{"ends with dot", "John Doe.", false}, // Dots are allowed at the end
		{"empty string", "", false},
		{"single character", "J", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasInvalidPunctuationPosition(tt.input)
			if result != tt.expected {
				t.Errorf("hasInvalidPunctuationPosition(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPunctuation(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"dot", '.', true},
		{"hyphen", '-', true},
		{"apostrophe", '\'', true},
		{"letter", 'a', false},
		{"number", '1', false},
		{"space", ' ', false},
		{"other", '!', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPunctuation(tt.input)
			if result != tt.expected {
				t.Errorf("isPunctuation(%c) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
