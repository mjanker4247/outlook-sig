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

{{ .Name }} (supports multiline: name, profession, title)
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
				Usage:   "Your name (can include profession/title on separate lines)",
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

			name, err := getOrPrompt(c.String("name"), "Enter your name: ")
			if err != nil {
				return fmt.Errorf("failed to get name: %v", err)
			}
			if name == "" {
				return fmt.Errorf("name cannot be empty")
			}
			// Convert literal \n to actual newlines
			name = strings.ReplaceAll(name, "\\n", "\n")

			email, err := getOrPrompt(c.String("email"), "Enter your email: ")
			if err != nil {
				return fmt.Errorf("failed to get email: %v", err)
			}
			if err := common.ValidateEmail(email); err != nil {
				return fmt.Errorf("invalid email: %v", err)
			}

			phone, err := getOrPrompt(c.String("phone"), "Enter your phone number: ")
			if err != nil {
				return fmt.Errorf("failed to get phone: %v", err)
			}
			if err := common.ValidatePhoneNumber(phone); err != nil {
				return fmt.Errorf("invalid phone number: %v", err)
			}

			phoneDisplay, phoneLink, err := common.FormatPhoneNumber(phone, "DE")
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

			templateBase, err := common.GetTemplateBase()
			if err != nil {
				return err
			}

			installer := signature.NewInstaller(templateBase)

			// Override template source if specified via CLI flag
			if templateSource := c.String("template-source"); templateSource != "" {
				if installer.Config == nil {
					if err := installer.LoadConfig(); err != nil {
						return fmt.Errorf("failed to load configuration: %v", err)
					}
				}
				installer.Config.TemplateSource = templateSource
			}

			return installer.Install(data)
		},
	}
}

func getOrPrompt(value, prompt string) (string, error) {
	if strings.TrimSpace(value) != "" {
		return value, nil
	}
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}
	return strings.TrimSpace(input), nil
}
