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

	for _, name := range wantNames {
		p := filepath.Join(inst.TemplateBase, "templates", name)
		if _, err := fs.Stat(p); err != nil {
			t.Errorf("expected file %q to be created, but stat failed: %v", p, err)
		}
	}
}