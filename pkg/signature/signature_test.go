///////////////////////////////////////////////////////
// file: signature/signature_test.go
///////////////////////////////////////////////////////

package signature

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

// writeConfig is a helper that writes a YAML config into fs at the given path.
func writeConfig(t *testing.T, fs afero.Fs, path, yaml string) {
	t.Helper()
	if err := fs.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", path, err)
	}
	if err := afero.WriteFile(fs, path, []byte(yaml), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

// TestLoadConfig_UserProfilePriority verifies user-profile config wins over exe-adjacent.
func TestLoadConfig_UserProfilePriority(t *testing.T) {
	fs := afero.NewMemMapFs()

	userPath, err := UserConfigPath()
	if err != nil {
		t.Fatalf("UserConfigPath(): %v", err)
	}

	writeConfig(t, fs, userPath,
		"template_name: UserTemplate\ntemplate_source: web\nbase_url: http://user/\n")
	writeConfig(t, fs, "/app/config.yaml",
		"template_name: ExeTemplate\ntemplate_source: local\n")

	inst := &Installer{TemplateBase: "/app/templates", fs: fs}
	if err := inst.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if inst.Config.TemplateName != "UserTemplate" {
		t.Errorf("TemplateName = %q, want %q", inst.Config.TemplateName, "UserTemplate")
	}
}

// TestLoadConfig_FallsBackToExeAdjacent verifies exe-adjacent config is used when no user config exists.
func TestLoadConfig_FallsBackToExeAdjacent(t *testing.T) {
	fs := afero.NewMemMapFs()

	writeConfig(t, fs, "/app/config.yaml",
		"template_name: ExeTemplate\ntemplate_source: local\n")

	inst := &Installer{TemplateBase: "/app/templates", fs: fs}
	if err := inst.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if inst.Config.TemplateName != "ExeTemplate" {
		t.Errorf("TemplateName = %q, want %q", inst.Config.TemplateName, "ExeTemplate")
	}
}

// TestLoadConfig_BothAbsent verifies an error is returned when neither config source exists.
func TestLoadConfig_BothAbsent(t *testing.T) {
	fs := afero.NewMemMapFs()

	inst := &Installer{TemplateBase: "/app/templates", fs: fs}
	if err := inst.LoadConfig(); err == nil {
		t.Fatal("LoadConfig() error = nil, want non-nil when both sources are absent")
	}
}

// TestLoadConfig_CorruptUserProfile_Fallback verifies that a corrupt user config falls back to exe-adjacent.
func TestLoadConfig_CorruptUserProfile_Fallback(t *testing.T) {
	fs := afero.NewMemMapFs()

	userPath, err := UserConfigPath()
	if err != nil {
		t.Fatalf("UserConfigPath(): %v", err)
	}

	writeConfig(t, fs, userPath, "not: valid: yaml: [[")
	writeConfig(t, fs, "/app/config.yaml",
		"template_name: FallbackTemplate\ntemplate_source: local\n")

	inst := &Installer{TemplateBase: "/app/templates", fs: fs}
	if err := inst.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v, want fallback to succeed", err)
	}

	if inst.Config.TemplateName != "FallbackTemplate" {
		t.Errorf("TemplateName = %q, want %q", inst.Config.TemplateName, "FallbackTemplate")
	}
}

func TestDownloadWebTemplates_UsesTemplateNameForFilenames(t *testing.T) {
	fs := afero.NewMemMapFs()

	var requested []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = append(requested, r.URL.Path[1:])
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	inst := &Installer{
		TemplateBase: "/app",
		fs:           fs,
	}
	inst.Config = &Config{
		TemplateName:   "MyTemplate",
		TemplateSource: "web",
		BaseURL:        srv.URL + "/", 
	}

	if err := inst.DownloadWebTemplates(); err != nil {
		t.Fatalf("DownloadWebTemplates() returned error: %v", err)
	}

	wantNames := []string{"MyTemplate.htm", "MyTemplate.txt"}

	if len(requested) != len(wantNames) {
		t.Fatalf("expected %d requests, got %d (%v)", len(wantNames), len(requested), requested)
	}
	for i, want := range wantNames {
		if requested[i] != want {
			t.Errorf("request[%d]: expected %q, got %q", i, want, requested[i])
		}
	}

	// TemplateBase is the templates directory itself; files are written directly into it.
	for _, name := range wantNames {
		p := filepath.Join(inst.TemplateBase, name)
		if _, err := fs.Stat(p); err != nil {
			t.Errorf("expected file %q to be created, but stat failed: %v", p, err)
		}
	}
}