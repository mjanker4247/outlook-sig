package signature

import (
	"testing"
)

func TestUnescapePhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTML entity &#43;",
			input:    "+49 123 4567890",
			expected: "+49 123 4567890",
		},
		{
			name:     "Escaped plus sign",
			input:    "\\+49 123 4567890",
			expected: "+49 123 4567890",
		},
		{
			name:     "HTML entity &plus;",
			input:    "&plus;49 123 4567890",
			expected: "+49 123 4567890",
		},
		{
			name:     "Multiple escaped characters",
			input:    "&#43;49 123 4567890",
			expected: "+49 123 4567890",
		},
		{
			name:     "No escaping needed",
			input:    "+49 123 4567890",
			expected: "+49 123 4567890",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unescapePhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("unescapePhoneNumber(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
