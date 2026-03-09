package cli

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCLIHelp(t *testing.T) {
	app := App()
	if app.Name != "Outlook Signature Installer" {
		t.Errorf("Expected app name 'Outlook Signature Installer', got '%s'", app.Name)
	}

	// Test that --name flag is defined (search by name to avoid fragile index assumptions).
	var nameFlag cli.Flag
	for _, f := range app.Flags {
		if f.Names()[0] == "name" {
			nameFlag = f
			break
		}
	}
	if nameFlag == nil {
		t.Errorf("Expected --name flag to be defined, but it was not found")
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
			name:     "trims surrounding whitespace from pre-supplied value",
			value:    "  test-value  ",
			prompt:   "Enter value:",
			expected: "test-value",
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

func TestGetUserInputValidation(t *testing.T) {
	tests := []struct {
		name        string
		cliName     string
		cliTitle    string
		expectError bool
	}{
		{
			name:        "valid input passes",
			cliName:     "Jane Doe",
			cliTitle:    "Senior Engineer",
			expectError: false,
		},
		{
			name:        "invalid name fails",
			cliName:     "!nvalid",
			cliTitle:    "Engineer",
			expectError: true,
		},
		// Titles are intentionally free-form; no validation error expected for special chars.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(tt.cliName, tt.cliTitle, "user@example.com", "+4915123456789")
			data, err := getUserInput(ctx)

			if tt.expectError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !tt.expectError {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if data.Name != tt.cliName {
					t.Fatalf("expected name %q, got %q", tt.cliName, data.Name)
				}

				if data.Title != tt.cliTitle {
					t.Fatalf("expected title %q, got %q", tt.cliTitle, data.Title)
				}
			}
		})
	}
}

func newTestContext(name, title, email, phone string) *cli.Context {
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.String("name", name, "")
	set.String("title", title, "")
	set.String("email", email, "")
	set.String("phone", phone, "")
	return cli.NewContext(App(), set, nil)
}
