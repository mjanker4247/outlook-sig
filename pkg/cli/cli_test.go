package cli

import (
	"testing"
)

func TestCLIHelp(t *testing.T) {
	app := App()
	if app.Name != "Outlook Signature Installer" {
		t.Errorf("Expected app name 'Outlook Signature Installer', got '%s'", app.Name)
	}

	// Test flag definitions
	nameFlag := app.Flags[1] // --name flag
	if nameFlag.Names()[0] != "name" {
		t.Errorf("Expected name flag, got %s", nameFlag.Names()[0])
	}
}
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
		{
			name:     "returns existing value with whitespace",
			value:    "  test-value  ",
			prompt:   "Enter value:",
			expected: "  test-value  ",
		},
		// Note: Testing empty/whitespace-only values would require stdin input
		// which is not suitable for automated testing
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getOrPrompt(tt.value, tt.prompt)
			if err != nil {
				t.Errorf("getOrPrompt(%q, %q) returned error: %v", tt.value, tt.prompt, err)
			}
			if result != tt.expected {
				t.Errorf("getOrPrompt(%q, %q) = %q, want %q", tt.value, tt.prompt, result, tt.expected)
			}
		})
	}
}
