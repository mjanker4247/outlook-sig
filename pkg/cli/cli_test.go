package cli

import (
	"testing"
)

func TestGetOrPrompt(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		prompt   string
		expected string
	}{
		{
			name:     "returns existing value",
			value:    "test-value",
			prompt:   "Enter value:",
			expected: "test-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getOrPrompt(tt.value, tt.prompt)
			if result != tt.expected {
				t.Errorf("getOrPrompt(%q, %q) = %q, want %q", tt.value, tt.prompt, result, tt.expected)
			}
		})
	}
}
