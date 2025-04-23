package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nyaruka/phonenumbers"
	"github.com/urfave/cli/v2"
)

type SignatureData struct {
	Name         string
	Email        string
	PhoneDisplay string
	PhoneLink    string
}

// Custom function to "unescape" the phone number (remove HTML-encoded entities like &#43;)
func unescapePhoneNumber(phone string) string {
	// Replace HTML encoded plus sign (&#43;) with the actual plus sign (+)
	return strings.ReplaceAll(phone, "&#43;", "+")
}

func getOrPrompt(value, prompt string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	return strings.TrimSpace(input)
}

func getExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

func formatPhoneNumber(phone string, countryCode string) (string, string, error) {
	num, err := phonenumbers.Parse(phone, countryCode)
	if err != nil {
		return phone, phone, err
	}

	display := phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
	link := phonenumbers.Format(num, phonenumbers.E164)

	return display, link, nil
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

		// It's a file
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

// Get Outlook signature directory
func getOutlookSignatureDir() (string, error) {
	var sigDir string
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not found")
		}
		sigDir = filepath.Join(appData, "Microsoft", "Signatures")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %v", err)
		}
		sigDir = filepath.Join(homeDir, "Library", "Group Containers", "UBF8T346G9.Office", "Outlook", "Outlook 15 Profiles", "Main Profile", "Signatures")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(sigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create signature directory: %v", err)
	}

	return sigDir, nil
}

func installSignature(data SignatureData, sigName string) error {
	exeDir, err := getExecutableDir()
	if err != nil {
		return fmt.Errorf("failed to get executable directory: %v", err)
	}

	templateBase := filepath.Join(exeDir, "templates")
	if _, err := os.Stat(templateBase); os.IsNotExist(err) {
		return fmt.Errorf("templates directory not found at %s", templateBase)
	}

	sigDir, err := getOutlookSignatureDir()
	if err != nil {
		return fmt.Errorf("failed to get signature directory: %v", err)
	}

	fmt.Println("Installing signature to:", sigDir)
	extensions := []string{".htm", ".txt"}
	var errors []error

	for _, ext := range extensions {
		templatePath := filepath.Join(templateBase, sigName+ext)
		destPath := filepath.Join(sigDir, sigName+ext)

		// Check if template exists
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("template file not found: %s", templatePath))
			continue
		}

		// Create a new template and register the custom function
		funcMap := template.FuncMap{
			"unescape": unescapePhoneNumber,
		}

		tpl, err := template.New(sigName).Funcs(funcMap).ParseFiles(templatePath)
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
			imageDirSrc := filepath.Join(templateBase, sigName+"_files")
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

func validateEmail(email string) error {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}
	// Basic validation - you might want to add more specific rules
	if len(phone) < 5 {
		return fmt.Errorf("phone number is too short")
	}
	return nil
}

func createCLIApp() *cli.App {
	return &cli.App{
		Name:  "Outlook Signature Installer",
		Usage: "Install a predefined Outlook signature with personal info",
		Description: `This tool installs Microsoft Outlook signatures from templates.
Templates may include these placeholders

{{ .Name }}
{{ .Email }}
{{ .PhoneLink }}
{{ .PhoneDisplay }}
	
Images and .htm/.txt files are copied and filled automatically.`,
		UsageText: "signature-installer.exe [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Your full name",
			},
			&cli.StringFlag{
				Name:    "email",
				Aliases: []string{"e"},
				Usage:   "Your email address",
			},
			&cli.StringFlag{
				Name:    "phone",
				Aliases: []string{"p"},
				Usage:   "Your phone number",
			},
			&cli.StringFlag{
				Name:    "sig",
				Aliases: []string{"s"},
				Usage:   "Signature base filename",
				Value:   "OutlookSignature",
			},
		},
		Action: func(c *cli.Context) error {
			name := getOrPrompt(c.String("name"), "Enter your name: ")
			if name == "" {
				return fmt.Errorf("name cannot be empty")
			}

			email := getOrPrompt(c.String("email"), "Enter your email: ")
			if err := validateEmail(email); err != nil {
				return fmt.Errorf("invalid email: %v", err)
			}

			phone := getOrPrompt(c.String("phone"), "Enter your phone number: ")
			if err := validatePhoneNumber(phone); err != nil {
				return fmt.Errorf("invalid phone number: %v", err)
			}

			sigName := c.String("sig")
			if strings.ContainsAny(sigName, `/\:*?"<>|`) {
				return fmt.Errorf("invalid signature name: contains invalid characters")
			}

			phoneDisplay, phoneLink, err := formatPhoneNumber(phone, "DE")
			if err != nil {
				fmt.Println("Warning: Could not format phone number. Using raw input.")
				phoneDisplay = phone
				phoneLink = phone
			}

			data := SignatureData{
				Name:         name,
				Email:        email,
				PhoneDisplay: phoneDisplay,
				PhoneLink:    phoneLink,
			}

			return installSignature(data, sigName)
		},
	}
}

func main() {
	app := createCLIApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
