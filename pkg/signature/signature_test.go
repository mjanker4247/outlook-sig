package signature

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"outlook-signature/pkg/common"
)

func TestDataToHTMLData(t *testing.T) {
	tests := []struct {
		name     string
		input    Data
		expected HTMLData
	}{
		{
			name: "single line name",
			input: Data{
				Name:         "John Doe",
				Email:        "john@example.com",
				PhoneDisplay: "+49 123 456789",
				PhoneLink:    "+49123456789",
			},
			expected: HTMLData{
				Name:         "John Doe",
				Email:        "john@example.com",
				PhoneDisplay: "+49 123 456789",
				PhoneLink:    "+49123456789",
			},
		},
		{
			name: "multiline name",
			input: Data{
				Name:         "John Doe\nSoftware Engineer\nSenior Developer",
				Email:        "john@example.com",
				PhoneDisplay: "+49 123 456789",
				PhoneLink:    "+49123456789",
			},
			expected: HTMLData{
				Name:         "John Doe<br>Software Engineer<br>Senior Developer",
				Email:        "john@example.com",
				PhoneDisplay: "+49 123 456789",
				PhoneLink:    "+49123456789",
			},
		},
		{
			name: "name with empty lines",
			input: Data{
				Name:         "John Doe\n\nSoftware Engineer\n\n\nSenior Developer",
				Email:        "john@example.com",
				PhoneDisplay: "+49 123 456789",
				PhoneLink:    "+49123456789",
			},
			expected: HTMLData{
				Name:         "John Doe<br>Software Engineer<br>Senior Developer",
				Email:        "john@example.com",
				PhoneDisplay: "+49 123 456789",
				PhoneLink:    "+49123456789",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToHTMLData()

			if result.Name != tt.expected.Name {
				t.Errorf("Name conversion failed:\nExpected: %q\nGot:      %q", tt.expected.Name, result.Name)
			}

			if result.Email != tt.expected.Email {
				t.Errorf("Email mismatch:\nExpected: %q\nGot:      %q", tt.expected.Email, result.Email)
			}

			if result.PhoneDisplay != tt.expected.PhoneDisplay {
				t.Errorf("PhoneDisplay mismatch:\nExpected: %q\nGot:      %q", tt.expected.PhoneDisplay, result.PhoneDisplay)
			}

			if result.PhoneLink != tt.expected.PhoneLink {
				t.Errorf("PhoneLink mismatch:\nExpected: %q\nGot:      %q", tt.expected.PhoneLink, result.PhoneLink)
			}
		})
	}
}

func TestMultilineNameHTMLConversion(t *testing.T) {
	// Test that newlines are properly converted to <br> tags
	input := "John Doe\nSoftware Engineer\nSenior Developer"
	expected := "John Doe<br>Software Engineer<br>Senior Developer"

	result := strings.ReplaceAll(input, "\n", "<br>")

	if result != expected {
		t.Errorf("Newline to <br> conversion failed:\nExpected: %q\nGot:      %q", expected, result)
	}
}

func TestCleanLineBreaks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line",
			input:    "John Doe",
			expected: "John Doe",
		},
		{
			name:     "two lines",
			input:    "John Doe\nSoftware Engineer",
			expected: "John Doe\nSoftware Engineer",
		},
		{
			name:     "multiple consecutive line breaks",
			input:    "John Doe\n\n\nSoftware Engineer\n\n\n\nSenior Developer",
			expected: "John Doe\nSoftware Engineer\nSenior Developer",
		},
		{
			name:     "lines with only whitespace",
			input:    "John Doe\n   \n\t\nSoftware Engineer\n  \n\nSenior Developer",
			expected: "John Doe\nSoftware Engineer\nSenior Developer",
		},
		{
			name:     "empty lines at start and end",
			input:    "\n\nJohn Doe\nSoftware Engineer\n\n",
			expected: "John Doe\nSoftware Engineer",
		},
		{
			name:     "mixed whitespace and content",
			input:    "  John Doe  \n  Software Engineer  \n  Senior Developer  ",
			expected: "John Doe\nSoftware Engineer\nSenior Developer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.CleanLineBreaks(tt.input)
			if result != tt.expected {
				t.Errorf("CleanLineBreaks failed:\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, tt.expected, result)
			}
		})
	}
}

func TestHTMLTemplateProcessing(t *testing.T) {
	// Create test data with multiline name
	data := Data{
		Name:         "John Doe\nSoftware Engineer\nSenior Developer",
		Email:        "john@example.com",
		PhoneDisplay: "+49 123 456789",
		PhoneLink:    "+49123456789",
	}

	// Test the toHTMLData conversion
	htmlData := data.ToHTMLData()
	expectedName := "John Doe<br>Software Engineer<br>Senior Developer"

	if string(htmlData.Name) != expectedName {
		t.Errorf("HTML conversion failed:\nExpected: %q\nGot:      %q", expectedName, string(htmlData.Name))
	}

	// Test that the safeHTML function would work correctly
	// This simulates what happens in the template
	funcMap := template.FuncMap{
		"safeHTML": func(s template.HTML) template.HTML {
			return s
		},
	}

	// Create a simple test template
	tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{ .Name | safeHTML }}`)
	if err != nil {
		t.Fatalf("Failed to parse test template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, htmlData); err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	result := buf.String()
	if result != expectedName {
		t.Errorf("Template execution failed:\nExpected: %q\nGot:      %q", expectedName, result)
	}

	// Verify that <br> tags are not escaped
	if strings.Contains(result, "&lt;br&gt;") {
		t.Errorf("HTML tags are still escaped: %q", result)
	}

	if !strings.Contains(result, "<br>") {
		t.Errorf("HTML tags are missing: %q", result)
	}
}

func TestHTMLTemplateProcessingWithMultipleLineBreaks(t *testing.T) {
	// Test with input that has multiple consecutive line breaks
	data := Data{
		Name:         "John Doe\n\n\nSoftware Engineer\n\n\n\nSenior Developer",
		Email:        "john@example.com",
		PhoneDisplay: "+49 123 456789",
		PhoneLink:    "+49123456789",
	}

	// Test the toHTMLData conversion - should clean up multiple line breaks
	htmlData := data.ToHTMLData()
	expectedName := "John Doe<br>Software Engineer<br>Senior Developer"

	if string(htmlData.Name) != expectedName {
		t.Errorf("HTML conversion with multiple line breaks failed:\nExpected: %q\nGot:      %q", expectedName, string(htmlData.Name))
	}

	// Verify that the output contains exactly 2 <br> tags (3 lines - 1)
	brCount := strings.Count(string(htmlData.Name), "<br>")
	if brCount != 2 {
		t.Errorf("Expected 2 <br> tags, got %d", brCount)
	}
}
