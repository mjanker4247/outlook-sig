// file: pkg/cli/cli.go
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"outlook-signature/pkg/common"
	"outlook-signature/pkg/gui"
	"outlook-signature/pkg/signature"

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
{{ .Title }}
{{ .Email }}
{{ .PhoneLink }}
{{ .PhoneDisplay }}
	
Images and .htm/.txt files are copied and filled automatically.
The template to use is configured in the config.yaml file.`,
		UsageText: "SignatureInstaller.exe [options]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "gui",
				Aliases: []string{"g"},
				Usage:   "Launch in GUI mode",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Your name",
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "Your profession or title (optional)",
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
				Name:    "template-source",
				Aliases: []string{"s"},
				Usage:   "Template source: 'local' or 'web' (overrides config.yaml)",
			},
		},
		Action: func(c *cli.Context) error {
			// Check if GUI mode is requested or no arguments are provided
			if c.Bool("gui") || len(os.Args) == 1 {
				gui.ShowGUI()
				return nil
			}

			return runCLIInstallation(c)
		},
	}
}

// runCLIInstallation handles the CLI installation flow
func runCLIInstallation(c *cli.Context) error {
	// Get and validate user input
	data, err := getUserInput(c)
	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	// Create and configure installer
	installer, err := createInstaller(c)
	if err != nil {
		return fmt.Errorf("failed to create installer: %v", err)
	}

	// Install signature
	return installer.Install(*data)
}

// getUserInput collects and validates user input
func getUserInput(c *cli.Context) (*signature.Data, error) {
	name, err := getOrPrompt(c.String("name"), "Enter your name: ")
	if err != nil {
		return nil, fmt.Errorf("failed to get name: %v", err)
	}
	if err := common.ValidateName(name); err != nil {
		return nil, fmt.Errorf("invalid name: %v", err)
	}

	// Title is optional:
	// - do NOT prompt if omitted on CLI
	// - do NOT complain if empty
	title := strings.TrimSpace(c.String("title"))
	if title != "" {
		if err := common.ValidateTitle(title); err != nil {
			return nil, fmt.Errorf("invalid title: %v", err)
		}
	}

	email, err := getOrPrompt(c.String("email"), "Enter your email: ")
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %v", err)
	}
	if err := common.ValidateEmail(email); err != nil {
		return nil, fmt.Errorf("invalid email: %v", err)
	}

	phone, err := getOrPrompt(c.String("phone"), "Enter your phone number: ")
	if err != nil {
		return nil, fmt.Errorf("failed to get phone: %v", err)
	}
	if err := common.ValidatePhoneNumber(phone); err != nil {
		return nil, fmt.Errorf("invalid phone number: %v", err)
	}

	phoneDisplay, phoneLink, err := common.FormatPhoneNumber(phone, "DE")
	if err != nil {
		fmt.Println("Warning: Could not format phone number. Using raw input.")
		phoneDisplay = phone
		phoneLink = phone
	}

	return &signature.Data{
		Name:         name,
		Title:        title, // may be empty string, that's fine
		Email:        email,
		PhoneDisplay: phoneDisplay,
		PhoneLink:    phoneLink,
	}, nil
}

// createInstaller creates and configures the signature installer
func createInstaller(c *cli.Context) (*signature.Installer, error) {
	templateBase, err := common.GetTemplateBase()
	if err != nil {
		return nil, err
	}

	installer := signature.NewInstaller(templateBase)

	// Override template source if specified via CLI flag
	if templateSource := c.String("template-source"); templateSource != "" {
		if installer.Config == nil {
			if err := installer.LoadConfig(); err != nil {
				return nil, fmt.Errorf("failed to load configuration: %v", err)
			}
		}
		installer.Config.TemplateSource = templateSource
	}

	return installer, nil
}

// getOrPrompt returns the trimmed flag value when non-empty, or prompts stdin and
// returns the trimmed user input.
func getOrPrompt(value, prompt string) (string, error) {
	if strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value), nil
	}
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}
	return strings.TrimSpace(input), nil
}
