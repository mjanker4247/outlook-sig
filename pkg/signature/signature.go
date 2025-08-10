package signature

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// WebTemplateConfig represents web template configuration
type WebTemplateConfig struct {
	BaseURL       string   `yaml:"base_url"`
	TemplateFiles []string `yaml:"template_files"`
}

// Config represents the application configuration
type Config struct {
	TemplateName   string             `yaml:"template_name"`
	TemplateSource string             `yaml:"template_source"`
	WebTemplates   *WebTemplateConfig `yaml:"web_templates,omitempty"`
}

// Data represents the signature data structure
type Data struct {
	Name         string
	Email        string
	PhoneDisplay string
	PhoneLink    string
}

// Installer handles signature installation
type Installer struct {
	TemplateBase string
	Config       *Config
	sigDir       string // Optional override for signature directory
	fs           afero.Fs
}

// NewInstaller creates a new signature installer
func NewInstaller(templateBase string) *Installer {
	return &Installer{
		TemplateBase: templateBase,
		fs:           afero.NewOsFs(),
	}
}

// LoadConfig loads the configuration from the build root directory
func (i *Installer) LoadConfig() error {
	// Get the build root directory (parent of templates directory)
	buildRoot := filepath.Dir(i.TemplateBase)
	configPath := filepath.Join(buildRoot, "config.yaml")

	configData, err := afero.ReadFile(i.fs, configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	if config.TemplateName == "" {
		return fmt.Errorf("template_name is required in config file")
	}

	i.Config = &config
	return nil
}

// DownloadWebTemplates downloads templates from the configured web server
func (i *Installer) DownloadWebTemplates() error {
	if i.Config == nil || i.Config.WebTemplates == nil {
		return fmt.Errorf("web templates configuration not found")
	}

	// Create templates directory if it doesn't exist
	templatesDir := filepath.Join(i.TemplateBase, "templates")
	if err := i.fs.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Download each template file
	for _, filename := range i.Config.WebTemplates.TemplateFiles {
		url := i.Config.WebTemplates.BaseURL + filename
		filepath := filepath.Join(templatesDir, filename)

		if err := i.downloadFile(client, url, filepath); err != nil {
			return fmt.Errorf("failed to download %s: %v", filename, err)
		}
	}

	return nil
}

// downloadFile downloads a single file from a URL
func (i *Installer) downloadFile(client *http.Client, url, filepath string) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get %s: status %d", url, resp.StatusCode)
	}

	// Create the file
	file, err := i.fs.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filepath, err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %v", filepath, err)
	}

	return nil
}

// GetOutlookSignatureDir returns the path to the Outlook signatures directory
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
		return filepath.Join(homeDir, "Library", "Group Containers", "UBF8T346G9.Office", "Outlook", "Outlook 15 Profiles", "Main Profile", "Signatures"), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// Replaces placeholders like {{ .Name }} or {{.Name}} etc.
func (i *Installer) replacePlaceholders(templateOrPath string, data Data) (string, error) {
	var content string

	// Check if the input is a file
	if _, err := i.fs.Stat(templateOrPath); err == nil {
		// Read the template file
		templateContent, err := afero.ReadFile(i.fs, templateOrPath)
		if err != nil {
			return "", fmt.Errorf("Failed to read template file %s: %v", templateOrPath, err)
		}
		content = string(templateContent)
	} else {
		content = templateOrPath
	}

	values := map[string]string{
		"Name":         data.Name,
		"Email":        data.Email,
		"PhoneDisplay": data.PhoneDisplay,
		"PhoneLink":    data.PhoneLink,
	}

	for key, val := range values {
		pattern := regexp.MustCompile(`{{\s*\.` + regexp.QuoteMeta(key) + `\s*}}`)
		content = pattern.ReplaceAllString(content, val)
	}
	return content, nil
}

// Install installs a signature with the given data
func (i *Installer) Install(data Data) error {
	if _, err := i.fs.Stat(i.TemplateBase); os.IsNotExist(err) {
		return fmt.Errorf("Templates directory not found at %s", i.TemplateBase)
	}

	// Load configuration if not already loaded
	if i.Config == nil {
		if err := i.LoadConfig(); err != nil {
			return fmt.Errorf("Failed to load configuration: %v", err)
		}
	}

	// Check template source and download if needed
	if i.Config.TemplateSource == "web" {
		if err := i.DownloadWebTemplates(); err != nil {
			return fmt.Errorf("Failed to download web templates: %v", err)
		}
	}

	// Use the configured template name
	sigName := i.Config.TemplateName

	// Validate template name and sanitize it
	if sigName == "" {
		return fmt.Errorf("Template name cannot be empty")
	}
	if strings.ContainsAny(sigName, `/\:*?"<>|`) {
		return fmt.Errorf("Invalid template name: contains invalid characters")
	}
	// Additional security: Clean the path and check for traversal attempts
	cleanSigName := filepath.Clean(sigName)
	if cleanSigName != sigName || strings.Contains(cleanSigName, "..") {
		return fmt.Errorf("Invalid template name: potential path traversal detected")
	}

	var sigDir string
	var err error
	if i.sigDir != "" {
		sigDir = i.sigDir
	} else {
		sigDir, err = GetOutlookSignatureDir()
		if err != nil {
			return fmt.Errorf("Failed to get signature directory: %v", err)
		}
	}

	// Create the signature directory if it doesn't exist
	if err := i.fs.MkdirAll(sigDir, 0755); err != nil {
		return fmt.Errorf("Failed to create signature directory: %v", err)
	}

	fmt.Println("Installing signature to:", sigDir)
	extensions := []string{".htm", ".txt"}
	var errors []error

	for _, ext := range extensions {
		templatePath := filepath.Join(i.TemplateBase, cleanSigName+ext)
		// Additional security: Verify the resolved path is within TemplateBase
		if !strings.HasPrefix(templatePath, i.TemplateBase) {
			errors = append(errors, fmt.Errorf("Invalid template path: outside of template directory"))
			continue
		}

		destPath := filepath.Join(sigDir, cleanSigName+ext)
		// Additional security: Verify the resolved path is within sigDir
		if !strings.HasPrefix(destPath, sigDir) {
			errors = append(errors, fmt.Errorf("Invalid destination path: outside of signature directory"))
			continue
		}

		if _, err := i.fs.Stat(templatePath); os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("Template file not found: %s", templatePath))
			continue
		}

		if ext == ".htm" {
			// Use html/template with size limit
			tpl, err := template.New(filepath.Base(templatePath)).ParseFiles(templatePath)
			if err != nil {
				errors = append(errors, fmt.Errorf("Failed to parse %s: %v", templatePath, err))
				continue
			}

			// Use LimitedBuffer to prevent memory exhaustion
			buf := &LimitedBuffer{
				Buffer: bytes.Buffer{},
				limit:  5 * 1024 * 1024, // 5MB limit
			}
			if err := tpl.Execute(buf, data); err != nil {
				errors = append(errors, fmt.Errorf("Failed to execute template %s: %v", templatePath, err))
				continue
			}

			if err := afero.WriteFile(i.fs, destPath, buf.Bytes(), 0644); err != nil {
				errors = append(errors, fmt.Errorf("Failed to write %s: %v", destPath, err))
				continue
			}

			imageDirSrc := filepath.Join(i.TemplateBase, cleanSigName+"_files")
			imageDirDst := filepath.Join(sigDir, cleanSigName+"_files")

			// Additional security: Verify image directory paths
			if !strings.HasPrefix(imageDirSrc, i.TemplateBase) || !strings.HasPrefix(imageDirDst, sigDir) {
				errors = append(errors, fmt.Errorf("Invalid image directory path"))
				continue
			}

			if _, err := i.fs.Stat(imageDirSrc); err == nil {
				if err := i.copyDir(imageDirSrc, imageDirDst); err != nil {
					errors = append(errors, fmt.Errorf("Failed to copy image folder: %v", err))
				} else {
					fmt.Printf("Copied image assets to %s\n", imageDirDst)
				}
			}
		} else if ext == ".txt" {
			// Perform replacements with size limit
			result, err := i.replacePlaceholders(templatePath, data)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			// Save the new file
			err = afero.WriteFile(i.fs, destPath, []byte(result), 0644)
			if err != nil {
				errors = append(errors, fmt.Errorf("Failed to write file %s: %v", destPath, err))
				continue
			}
		} else {
			errors = append(errors, fmt.Errorf("Unsupported file extension: %s", ext))
		}

		fmt.Printf("Created: %s\n", destPath)
	}

	if len(errors) > 0 {
		return fmt.Errorf("Encountered %d errors during installation: %v", len(errors), errors)
	}

	return nil
}

// LimitedBuffer is a buffer with a size limit to prevent memory exhaustion
type LimitedBuffer struct {
	bytes.Buffer
	limit int
}

func (b *LimitedBuffer) Write(p []byte) (n int, err error) {
	if b.Buffer.Len()+len(p) > b.limit {
		return 0, fmt.Errorf("buffer size limit exceeded")
	}
	return b.Buffer.Write(p)
}

func (i *Installer) copyDir(src string, dst string) error {
	// Clean and verify paths
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if strings.Contains(src, "..") || strings.Contains(dst, "..") {
		return fmt.Errorf("path traversal detected")
	}

	return afero.Walk(i.fs, src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)
		// Verify the target path is within dst directory
		if !strings.HasPrefix(targetPath, dst) {
			return fmt.Errorf("invalid target path: outside of destination directory")
		}

		if info.IsDir() {
			return i.fs.MkdirAll(targetPath, 0755)
		}

		// Size limit for files
		if info.Size() > 50*1024*1024 { // 50MB limit
			return fmt.Errorf("file too large: %s", path)
		}

		srcFile, err := i.fs.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Create file with restricted permissions
		dstFile, err := i.fs.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// Use io.CopyN to limit the amount of data copied
		_, err = io.CopyN(dstFile, srcFile, 50*1024*1024) // 50MB limit
		if err != nil && err != io.EOF {
			return err
		}
		return nil
	})
}
