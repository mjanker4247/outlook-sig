package signature

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"outlook-signature/pkg/common"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// placeholder maps a pre-compiled regex to its corresponding Data field.
type placeholder struct {
	re    *regexp.Regexp
	value func(Data) string
}

// placeholders holds pre-compiled regex patterns for .txt template substitution.
var placeholders = []placeholder{
	{regexp.MustCompile(`{{\s*\.Name\s*}}`), func(d Data) string { return d.Name }},
	{regexp.MustCompile(`{{\s*\.Title\s*}}`), func(d Data) string { return d.Title }},
	{regexp.MustCompile(`{{\s*\.Email\s*}}`), func(d Data) string { return d.Email }},
	{regexp.MustCompile(`{{\s*\.PhoneDisplay\s*}}`), func(d Data) string { return d.PhoneDisplay }},
	{regexp.MustCompile(`{{\s*\.PhoneLink\s*}}`), func(d Data) string { return d.PhoneLink }},
}

// Data represents the signature data structure.
type Data struct {
	Name         string
	Title        string
	Email        string
	PhoneDisplay string
	PhoneLink    string
}

// Installer handles signature installation.
type Installer struct {
	TemplateBase string   // path to templates directory
	Config       *Config  // loaded configuration
	sigDir       string   // optional override for signature directory
	fs           afero.Fs // filesystem abstraction
	logger       *slog.Logger
}

// NewInstaller creates a new signature installer.
func NewInstaller(templateBase string) *Installer {
	return &Installer{
		TemplateBase: templateBase,
		fs:           afero.NewOsFs(),
		logger:       slog.Default(),
	}
}

// Logger returns the configured logger, defaulting to slog.Default when nil.
func (i *Installer) Logger() *slog.Logger {
	if i.logger == nil {
		i.logger = slog.Default()
	}
	return i.logger
}

// LoadConfig loads configuration using priority order:
//  1. User-profile config (~/…/OutlookSignatureInstaller/config.yaml)
//  2. Exe-adjacent config (<templateBase>/../config.yaml)
//
// The first source that exists and parses successfully is used. A corrupt or
// invalid user-profile config is logged and skipped in favour of the fallback.
// An error is returned only when both sources fail.
func (i *Installer) LoadConfig() error {
	// 1. Try user-profile config.
	cfg, err := LoadUserConfig(i.fs)
	if err != nil {
		i.Logger().Warn("failed to load user config, falling back to exe-adjacent", slog.String("error", err.Error()))
	} else if cfg != nil {
		if cfg.TemplateName == "" {
			i.Logger().Warn("user config has empty template_name, falling back to exe-adjacent")
		} else {
			i.Config = cfg
			i.Logger().Info("configuration loaded from user profile",
				slog.String("template_name", cfg.TemplateName),
				slog.String("source", cfg.TemplateSource),
			)
			return nil
		}
	}

	// 2. Fall back to exe-adjacent config.
	buildRoot := filepath.Dir(i.TemplateBase)
	configPath := filepath.Join(buildRoot, "config.yaml")
	return i.loadConfigFromPath(configPath)
}

// loadConfigFromPath reads and validates a config file at the given path.
func (i *Installer) loadConfigFromPath(path string) error {
	i.Logger().Info("loading configuration", slog.String("path", path))

	configData, err := afero.ReadFile(i.fs, path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config file %s: %v", path, err)
	}

	if config.TemplateName == "" {
		return fmt.Errorf("template_name is required in config file %s", path)
	}

	i.Config = &config
	i.Logger().Info(
		"configuration loaded",
		slog.String("template_name", config.TemplateName),
		slog.String("source", config.TemplateSource),
		slog.String("base_url", config.BaseURL),
	)
	return nil
}

// DownloadWebTemplates downloads templates from the configured web server.
// It derives template file names from TemplateName (e.g. <name>.htm, <name>.txt).
func (i *Installer) DownloadWebTemplates() error {
	if i.Config == nil {
		return fmt.Errorf("configuration not loaded")
	}
	if i.Config.TemplateSource != "web" {
		return fmt.Errorf("template_source is %q, expected \"web\"", i.Config.TemplateSource)
	}
	if i.Config.BaseURL == "" {
		return fmt.Errorf("base_url must be set in config for web templates")
	}

	templateName := i.Config.TemplateName
	if templateName == "" {
		return fmt.Errorf("template_name must be set to download web templates")
	}

	i.Logger().Info("downloading web templates", slog.String("base_url", i.Config.BaseURL))

	// TemplateBase is the templates directory itself.
	templatesDir := i.TemplateBase
	if err := i.fs.MkdirAll(templatesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create templates directory: %v", err)
	}

	extensions := []string{".htm", ".txt"}
	var filenames []string
	for _, ext := range extensions {
		filenames = append(filenames, templateName+ext)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for _, filename := range filenames {
		url := i.Config.BaseURL + filename
		targetPath := filepath.Join(templatesDir, filename)

		i.Logger().Info(
			"downloading template file",
			slog.String("filename", filename),
			slog.String("url", url),
			slog.String("target", targetPath),
		)

		if err := i.downloadFile(client, url, targetPath); err != nil {
			return fmt.Errorf("failed to download %s: %v", filename, err)
		}

		i.Logger().Info(
			"downloaded template file",
			slog.String("filename", filename),
			slog.String("path", targetPath),
		)
	}

	i.Logger().Info("web templates download complete", slog.Int("file_count", len(filenames)))
	return nil
}

// downloadFile downloads a single file from a URL into the fs.
// The response body is limited to common.FileSizeLimit bytes via io.LimitReader
// to prevent memory exhaustion from oversized responses.
func (i *Installer) downloadFile(client *http.Client, url, path string) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get %s: status %d", url, resp.StatusCode)
	}

	file, err := i.fs.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", path, err)
	}
	defer file.Close()

	limited := io.LimitReader(resp.Body, common.FileSizeLimit)
	if _, err = io.Copy(file, limited); err != nil {
		return fmt.Errorf("failed to write to file %s: %v", path, err)
	}

	return nil
}

// GetOutlookSignatureDir returns the path to the Outlook signatures directory.
func GetOutlookSignatureDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not found")
		}
		return filepath.Join(appData, "Microsoft", "Signatures"), nil
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %v", err)
		}
		return filepath.Join(
			homeDir,
			"Library",
			"Group Containers",
			"UBF8T346G9.Office",
			"Outlook",
			"Outlook 15 Profiles",
			"Main Profile",
			"Signatures",
		), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// replacePlaceholders replaces {{ .Field }} placeholders in a .txt template string
// using pre-compiled regex patterns. File reading is the caller's responsibility.
func replacePlaceholders(content string, data Data) string {
	for _, p := range placeholders {
		content = p.re.ReplaceAllString(content, p.value(data))
	}
	return content
}

// Install installs a signature with the given data.
func (i *Installer) Install(data Data) error {
	if err := i.validateTemplateBase(); err != nil {
		return err
	}

	if err := i.ensureConfigLoaded(); err != nil {
		return err
	}

	if err := i.handleWebTemplates(); err != nil {
		return err
	}

	sigName, err := i.validateAndSanitizeTemplateName()
	if err != nil {
		return err
	}

	sigDir, err := i.getSignatureDirectory()
	if err != nil {
		return err
	}

	return i.installSignatureFiles(sigName, sigDir, data)
}

// validateTemplateBase checks if the template base directory exists.
// Any stat error — including permission-denied — is returned, not just "not exist".
func (i *Installer) validateTemplateBase() error {
	if _, err := i.fs.Stat(i.TemplateBase); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("templates directory not found at %s", i.TemplateBase)
		}
		return fmt.Errorf("failed to access templates directory %s: %v", i.TemplateBase, err)
	}
	return nil
}

// ensureConfigLoaded loads configuration if not already loaded.
func (i *Installer) ensureConfigLoaded() error {
	if i.Config == nil {
		if err := i.LoadConfig(); err != nil {
			return fmt.Errorf("failed to load configuration: %v", err)
		}
	}
	return nil
}

// handleWebTemplates downloads web templates if needed.
// BaseURL validation is delegated to DownloadWebTemplates.
func (i *Installer) handleWebTemplates() error {
	if i.Config.TemplateSource == "web" {
		if err := i.DownloadWebTemplates(); err != nil {
			return fmt.Errorf("failed to download web templates: %v", err)
		}
	}
	return nil
}

// validateAndSanitizeTemplateName validates and sanitizes the template name.
func (i *Installer) validateAndSanitizeTemplateName() (string, error) {
	sigName := i.Config.TemplateName
	if sigName == "" {
		return "", fmt.Errorf("template name cannot be empty")
	}

	if strings.ContainsAny(sigName, `/\:*?"<>|`) {
		return "", fmt.Errorf("invalid template name: contains invalid characters")
	}

	cleanSigName := filepath.Clean(sigName)
	if cleanSigName != sigName || strings.Contains(cleanSigName, "..") {
		return "", fmt.Errorf("invalid template name: potential path traversal detected")
	}

	return cleanSigName, nil
}

// getSignatureDirectory returns the signature directory path.
func (i *Installer) getSignatureDirectory() (string, error) {
	if i.sigDir != "" {
		return i.sigDir, nil
	}

	sigDir, err := GetOutlookSignatureDir()
	if err != nil {
		return "", fmt.Errorf("failed to get signature directory: %v", err)
	}

	if err := i.fs.MkdirAll(sigDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create signature directory: %v", err)
	}

	return sigDir, nil
}

// installSignatureFiles installs the signature files with the given data.
func (i *Installer) installSignatureFiles(sigName, sigDir string, data Data) error {
	fmt.Println("Installing signature to:", sigDir)

	extensions := []string{".htm", ".txt"}
	var errs []error

	for _, ext := range extensions {
		if err := i.installFile(sigName, sigDir, ext, data); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during installation: %v", len(errs), errs)
	}

	return nil
}

// installFile installs a single signature file.
func (i *Installer) installFile(sigName, sigDir, ext string, data Data) error {
	templatePath := filepath.Join(i.TemplateBase, sigName+ext)

	// Security: ensure templatePath is contained within TemplateBase.
	// Use filepath.Rel to avoid prefix-collision bugs (e.g. /templates-evil vs /templates).
	if rel, err := filepath.Rel(i.TemplateBase, templatePath); err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("invalid template path: outside of template directory")
	}

	destPath := filepath.Join(sigDir, sigName+ext)

	// Security: ensure destPath is contained within sigDir.
	if rel, err := filepath.Rel(sigDir, destPath); err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("invalid destination path: outside of signature directory")
	}

	if _, err := i.fs.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("template file not found: %s", templatePath)
	}

	switch ext {
	case ".htm":
		return i.installHTMLFile(templatePath, destPath, data)
	case ".txt":
		return i.installTextFile(templatePath, destPath, data)
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
}

// installHTMLFile installs an HTML signature file.
// The template is read via afero so the filesystem abstraction is respected in tests.
func (i *Installer) installHTMLFile(templatePath, destPath string, data Data) error {
	raw, err := afero.ReadFile(i.fs, templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %v", templatePath, err)
	}

	tpl, err := htmltemplate.New(filepath.Base(templatePath)).Parse(string(raw))
	if err != nil {
		return fmt.Errorf("failed to parse %s: %v", templatePath, err)
	}

	var buf bytes.Buffer

	i.Logger().Info("rendering HTML template", slog.String("template", templatePath))

	if err := tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %v", templatePath, err)
	}

	if buf.Len() > common.BufferSizeLimit {
		return fmt.Errorf("rendered template exceeds buffer size limit of %d bytes", common.BufferSizeLimit)
	}

	if err := afero.WriteFile(i.fs, destPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %v", destPath, err)
	}

	fmt.Printf("Created: %s\n", destPath)
	i.Logger().Info("HTML template rendered", slog.String("destination", destPath))
	return nil
}

// installTextFile installs a text signature file.
func (i *Installer) installTextFile(templatePath, destPath string, data Data) error {
	raw, err := afero.ReadFile(i.fs, templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %v", templatePath, err)
	}

	result := replacePlaceholders(string(raw), data)

	if err := afero.WriteFile(i.fs, destPath, []byte(result), 0o644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", destPath, err)
	}

	fmt.Printf("Created: %s\n", destPath)
	i.Logger().Info("text template rendered", slog.String("destination", destPath))
	return nil
}
