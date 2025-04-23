package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"outlook-signature/pkg/signature"

	"github.com/nyaruka/phonenumbers"
	"github.com/urfave/cli/v2"
)

// App creates and returns the CLI application
func App() *cli.App {
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

			data := signature.Data{
				Name:         name,
				Email:        email,
				PhoneDisplay: phoneDisplay,
				PhoneLink:    phoneLink,
			}

			exeDir, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %v", err)
			}
			templateBase := filepath.Join(filepath.Dir(exeDir), "templates")

			installer := signature.NewInstaller(templateBase)
			return installer.Install(data, sigName)
		},
	}
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
	if len(phone) < 5 {
		return fmt.Errorf("phone number is too short")
	}
	return nil
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
