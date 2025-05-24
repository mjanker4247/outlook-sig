package common

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetTemplateBase returns the base directory for templates
func GetTemplateBase() (string, error) {
	exeDir, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	return filepath.Join(filepath.Dir(exeDir), "templates"), nil
}

// GetAvailableTemplates returns a list of available HTML templates
func GetAvailableTemplates() ([]string, error) {
	templateBase, err := GetTemplateBase()
	if err != nil {
		return nil, err
	}

	// Read the templates directory
	entries, err := os.ReadDir(templateBase)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %v", err)
	}

	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".htm" {
			// Remove the .htm extension for display
			name := entry.Name()[:len(entry.Name())-4]
			templates = append(templates, name)
		}
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no HTML templates found in %s", templateBase)
	}

	return templates, nil
}
