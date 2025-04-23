package signature

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

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
	sigDir       string // Optional override for signature directory
}

// NewInstaller creates a new signature installer
func NewInstaller(templateBase string) *Installer {
	return &Installer{
		TemplateBase: templateBase,
	}
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

// Install installs a signature with the given data
func (i *Installer) Install(data Data, sigName string) error {
	if _, err := os.Stat(i.TemplateBase); os.IsNotExist(err) {
		return fmt.Errorf("templates directory not found at %s", i.TemplateBase)
	}

	var sigDir string
	var err error
	if i.sigDir != "" {
		sigDir = i.sigDir
	} else {
		sigDir, err = GetOutlookSignatureDir()
		if err != nil {
			return fmt.Errorf("failed to get signature directory: %v", err)
		}
	}

	// Create the signature directory if it doesn't exist
	if err := os.MkdirAll(sigDir, 0755); err != nil {
		return fmt.Errorf("failed to create signature directory: %v", err)
	}

	fmt.Println("Installing signature to:", sigDir)
	extensions := []string{".htm", ".txt"}
	var errors []error

	for _, ext := range extensions {
		templatePath := filepath.Join(i.TemplateBase, sigName+ext)
		destPath := filepath.Join(sigDir, sigName+ext)

		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("template file not found: %s", templatePath))
			continue
		}

		funcMap := template.FuncMap{
			"unescape": unescapePhoneNumber,
		}

		// Use html/template for both file types to ensure consistent escaping
		tpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to parse %s: %v", templatePath, err))
			continue
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, data); err != nil {
			errors = append(errors, fmt.Errorf("failed to execute template %s: %v", templatePath, err))
			continue
		}

		if err := os.WriteFile(destPath, buf.Bytes(), 0644); err != nil {
			errors = append(errors, fmt.Errorf("failed to write %s: %v", destPath, err))
			continue
		}

		if ext == ".htm" {
			imageDirSrc := filepath.Join(i.TemplateBase, sigName+"_files")
			imageDirDst := filepath.Join(sigDir, sigName+"_files")
			if _, err := os.Stat(imageDirSrc); err == nil {
				if err := copyDir(imageDirSrc, imageDirDst); err != nil {
					errors = append(errors, fmt.Errorf("failed to copy image folder: %v", err))
				} else {
					fmt.Printf("Copied image assets to %s\n", imageDirDst)
				}
			}
		}

		fmt.Printf("Created: %s\n", destPath)
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during installation: %v", len(errors), errors)
	}

	return nil
}

func unescapePhoneNumber(phone string) string {
	// First replace HTML entity
	phone = strings.ReplaceAll(phone, "&#43;", "+")
	// Then ensure any remaining + signs are not escaped
	return strings.ReplaceAll(phone, "+", "+")
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
