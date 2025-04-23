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
	input, _ := reader.ReadString('\n')
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
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("APPDATA environment variable not found")
	}
	return filepath.Join(appData, "Microsoft", "Signatures"), nil
}

func installSignature(data SignatureData, sigName string) error {
	exeDir, err := getExecutableDir()
	if err != nil {
		return err
	}
	templateBase := filepath.Join(exeDir, "templates")
	sigDir, err := getOutlookSignatureDir()
	if err != nil {
		return err
	}
	fmt.Println("Installing signature to:", sigDir)
	extensions := []string{".htm", ".txt"}

	for _, ext := range extensions {
		templatePath := filepath.Join(templateBase, sigName+ext)
		destPath := filepath.Join(sigDir, sigName+ext)

		// Create a new template and register the custom function
		funcMap := template.FuncMap{
			"unescape": unescapePhoneNumber, // Register unescape function
		}

		tpl, err := template.New(sigName).Funcs(funcMap).ParseFiles(templatePath)
		if err != nil {
			fmt.Printf("Failed to parse %s: %v\n", templatePath, err)
			continue
		}

		var buf bytes.Buffer
		err = tpl.Execute(&buf, data)
		if err != nil {
			fmt.Printf("Failed to execute template %s: %v\n", templatePath, err)
			continue
		}

		content := buf.String()
		
		err = os.WriteFile(destPath, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Failed to write %s: %v\n", destPath, err)
			continue
		}
		
		if ext == ".htm" {
			imageDirSrc := filepath.Join(templateBase, sigName+"_files")
			imageDirDst := filepath.Join(sigDir, sigName+"_files")
			if _, err := os.Stat(imageDirSrc); err == nil {
				err := copyDir(imageDirSrc, imageDirDst)
				if err != nil {
					fmt.Printf("Failed to copy image folder: %v\n", err)
				} else {
					fmt.Printf("Copied image assets to %s\n", imageDirDst)
				}
			}
		}

		fmt.Printf("Created: %s\n", destPath)
	}

	return nil
}


func main() {
	app := &cli.App{
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
			email := getOrPrompt(c.String("email"), "Enter your email: ")
			phone := getOrPrompt(c.String("phone"), "Enter your phone number: ")
			sigName := c.String("sig")

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

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
