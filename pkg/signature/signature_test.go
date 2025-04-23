package signature

import (
	"testing"
)

func TestPhoneNumberUnescape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"&#43;49123456789", "+49123456789"},
		{"+49123456789", "+49123456789"},
		{"&#43;49 123 456789", "+49 123 456789"},
		{"", ""},
	}

	for _, test := range tests {
		result := unescapePhoneNumber(test.input)
		if result != test.expected {
			t.Errorf("unescapePhoneNumber(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}
