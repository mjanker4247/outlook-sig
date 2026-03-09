package signature

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// appName is the subdirectory name used inside the OS user config directory.
const appName = "OutlookSignatureInstaller"

// Config represents the application configuration.
type Config struct {
	TemplateName   string `yaml:"template_name"`
	TemplateSource string `yaml:"template_source"` // "local" or "web"
	BaseURL        string `yaml:"base_url"`        // used when TemplateSource == "web"
}

// UserConfigPath returns the path to the user-profile config file:
//
//	macOS:   ~/Library/Application Support/OutlookSignatureInstaller/config.yaml
//	Windows: %APPDATA%\OutlookSignatureInstaller\config.yaml
//
// It delegates to os.UserConfigDir() which is available since Go 1.13.
func UserConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user config directory: %v", err)
	}
	return filepath.Join(dir, appName, "config.yaml"), nil
}

// LoadUserConfig reads config from the user-profile location using fs.
// Returns (nil, nil) when the file does not exist so callers can distinguish
// "file absent" from a real I/O or parse error.
func LoadUserConfig(fs afero.Fs) (*Config, error) {
	path, err := UserConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := afero.ReadFile(fs, path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read user config at %s: %v", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse user config at %s: %v", path, err)
	}
	return &cfg, nil
}

// SaveUserConfig marshals cfg to YAML and writes it to the user-profile location.
// The directory is created with MkdirAll (mode 0o700) if it does not exist.
func SaveUserConfig(fs afero.Fs, cfg *Config) error {
	path, err := UserConfigPath()
	if err != nil {
		return err
	}

	if err := fs.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("failed to create user config directory: %v", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := afero.WriteFile(fs, path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write user config at %s: %v", path, err)
	}
	return nil
}
