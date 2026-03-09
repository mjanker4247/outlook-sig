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
