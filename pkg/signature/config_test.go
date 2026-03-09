package signature

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

// TestUserConfigPath verifies the returned path has the expected structure.
func TestUserConfigPath(t *testing.T) {
	path, err := UserConfigPath()
	if err != nil {
		t.Fatalf("UserConfigPath() returned error: %v", err)
	}
	if path == "" {
		t.Fatal("UserConfigPath() returned empty string")
	}
	if !strings.HasSuffix(path, filepath.Join(appName, "config.yaml")) {
		t.Errorf("UserConfigPath() = %q, want path ending in %q",
			path, filepath.Join(appName, "config.yaml"))
	}
}

// TestLoadUserConfig_FileAbsent verifies (nil, nil) is returned when no file exists.
func TestLoadUserConfig_FileAbsent(t *testing.T) {
	fs := afero.NewMemMapFs()

	cfg, err := LoadUserConfig(fs)
	if err != nil {
		t.Fatalf("LoadUserConfig() error = %v, want nil", err)
	}
	if cfg != nil {
		t.Fatalf("LoadUserConfig() cfg = %+v, want nil", cfg)
	}
}

// TestLoadUserConfig_ValidFile verifies a valid YAML file is parsed correctly.
func TestLoadUserConfig_ValidFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	path, err := UserConfigPath()
	if err != nil {
		t.Fatalf("UserConfigPath() error: %v", err)
	}

	if err := fs.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("MkdirAll error: %v", err)
	}
	yaml := "template_name: MyTemplate\ntemplate_source: web\nbase_url: http://srv/\n"
	if err := afero.WriteFile(fs, path, []byte(yaml), 0o600); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg, err := LoadUserConfig(fs)
	if err != nil {
		t.Fatalf("LoadUserConfig() error = %v, want nil", err)
	}
	if cfg == nil {
		t.Fatal("LoadUserConfig() returned nil cfg, want non-nil")
	}
	if cfg.TemplateName != "MyTemplate" {
		t.Errorf("TemplateName = %q, want %q", cfg.TemplateName, "MyTemplate")
	}
	if cfg.TemplateSource != "web" {
		t.Errorf("TemplateSource = %q, want %q", cfg.TemplateSource, "web")
	}
	if cfg.BaseURL != "http://srv/" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "http://srv/")
	}
}

// TestLoadUserConfig_InvalidYAML verifies a parse error is returned for corrupt files.
func TestLoadUserConfig_InvalidYAML(t *testing.T) {
	fs := afero.NewMemMapFs()

	path, err := UserConfigPath()
	if err != nil {
		t.Fatalf("UserConfigPath() error: %v", err)
	}

	if err := fs.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("MkdirAll error: %v", err)
	}
	if err := afero.WriteFile(fs, path, []byte("not: valid: yaml: [["), 0o600); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg, err := LoadUserConfig(fs)
	if err == nil {
		t.Fatalf("LoadUserConfig() error = nil, want non-nil for invalid YAML")
	}
	if cfg != nil {
		t.Errorf("LoadUserConfig() cfg = %+v, want nil on error", cfg)
	}
}

// TestSaveUserConfig_CreatesDir verifies the config directory is created when absent.
func TestSaveUserConfig_CreatesDir(t *testing.T) {
	fs := afero.NewMemMapFs()

	cfg := &Config{TemplateName: "T", TemplateSource: "local", BaseURL: ""}
	if err := SaveUserConfig(fs, cfg); err != nil {
		t.Fatalf("SaveUserConfig() error = %v", err)
	}

	path, _ := UserConfigPath()
	if _, err := fs.Stat(filepath.Dir(path)); err != nil {
		t.Errorf("config directory not created: %v", err)
	}
	if _, err := fs.Stat(path); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}

// TestSaveUserConfig_RoundTrip verifies that save followed by load returns equal data.
func TestSaveUserConfig_RoundTrip(t *testing.T) {
	fs := afero.NewMemMapFs()

	want := &Config{
		TemplateName:   "RoundTrip",
		TemplateSource: "web",
		BaseURL:        "https://example.com/templates/",
	}

	if err := SaveUserConfig(fs, want); err != nil {
		t.Fatalf("SaveUserConfig() error = %v", err)
	}

	got, err := LoadUserConfig(fs)
	if err != nil {
		t.Fatalf("LoadUserConfig() error = %v", err)
	}
	if got == nil {
		t.Fatal("LoadUserConfig() returned nil after save")
	}

	if got.TemplateName != want.TemplateName {
		t.Errorf("TemplateName: got %q, want %q", got.TemplateName, want.TemplateName)
	}
	if got.TemplateSource != want.TemplateSource {
		t.Errorf("TemplateSource: got %q, want %q", got.TemplateSource, want.TemplateSource)
	}
	if got.BaseURL != want.BaseURL {
		t.Errorf("BaseURL: got %q, want %q", got.BaseURL, want.BaseURL)
	}
}
