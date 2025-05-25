package gui

import (
	"fmt"
	"strings"

	"outlook-signature/pkg/common"
	"outlook-signature/pkg/signature"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/asaskevich/govalidator"
)

// validateName performs comprehensive validation on a name string
func validateName(name string) error {
	// Trim whitespace only at the beginning and end
	name = strings.TrimSpace(name)

	// Check if empty
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	// Check if name is too short (less than 2 characters)
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}

	// Check if name contains only whitespace characters
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot contain only spaces")
	}

	// Use govalidator to check if the name is valid
	if !govalidator.IsAlpha(name) && !govalidator.Matches(name, "^[a-zA-Z\\s\\.\\-']+$") {
		return fmt.Errorf("name can only contain letters, spaces, dots, hyphens, and apostrophes")
	}

	// Check for multiple consecutive spaces
	if strings.Contains(name, "  ") {
		return fmt.Errorf("name cannot contain multiple consecutive spaces")
	}

	return nil
}

// ShowGUI displays the signature installer GUI
func ShowGUI() {
	myApp := app.New()
	window := myApp.NewWindow("Outlook Signature Installer")

	// Create form fields
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Your full name")
	nameEntry.Validator = validateName
	nameEntry.OnChanged = func(s string) {
		if err := nameEntry.Validate(); err != nil {
			nameEntry.SetValidationError(err)
			// Force refresh of the widget
			nameEntry.Refresh()
		} else {
			nameEntry.SetValidationError(nil)
			nameEntry.Refresh()
		}
	}

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Your email address")
	emailEntry.Validator = func(s string) error {
		return common.ValidateEmail(s)
	}
	emailEntry.OnChanged = func(s string) {
		if err := emailEntry.Validate(); err != nil {
			emailEntry.SetValidationError(err)
			// Force refresh of the widget
			emailEntry.Refresh()
		} else {
			emailEntry.SetValidationError(nil)
			emailEntry.Refresh()
		}
	}

	phoneEntry := widget.NewEntry()
	phoneEntry.SetPlaceHolder("Your phone number")
	phoneEntry.Validator = func(s string) error {
		return common.ValidatePhoneNumber(s)
	}
	phoneEntry.OnChanged = func(s string) {
		if err := phoneEntry.Validate(); err != nil {
			phoneEntry.SetValidationError(err)
			// Force refresh of the widget
			phoneEntry.Refresh()
		} else {
			phoneEntry.SetValidationError(nil)
			phoneEntry.Refresh()
		}
	}

	// Get template base directory
	templateBase, err := common.GetTemplateBase()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Get available templates
	templates, err := common.GetAvailableTemplates()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Create template selection dropdown
	templateSelect := widget.NewSelect(templates, func(selected string) {
		// This function is called when a template is selected
	})
	if len(templates) > 0 {
		templateSelect.SetSelected(templates[0])
	}

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Phone", Widget: phoneEntry},
			{Text: "Template", Widget: templateSelect},
		},
		OnSubmit: func() {
			// Validate inputs
			if err := nameEntry.Validate(); err != nil {
				dialog.ShowError(err, window)
				return
			}

			if err := emailEntry.Validate(); err != nil {
				dialog.ShowError(err, window)
				return
			}

			if err := phoneEntry.Validate(); err != nil {
				dialog.ShowError(err, window)
				return
			}

			// Format phone number
			phoneDisplay, phoneLink, err := common.FormatPhoneNumber(phoneEntry.Text, "DE")
			if err != nil {
				dialog.ShowError(fmt.Errorf("could not format phone number: %v", err), window)
				return
			}

			// Create signature data
			data := signature.Data{
				Name:         nameEntry.Text,
				Email:        emailEntry.Text,
				PhoneDisplay: phoneDisplay,
				PhoneLink:    phoneLink,
			}

			// Install signature
			installer := signature.NewInstaller(templateBase)
			err = installer.Install(data, templateSelect.Selected)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			dialog.ShowInformation("Success", "Signature installed successfully!", window)
		},
	}

	// Create main container
	content := container.NewVBox(
		widget.NewLabel("Enter your signature details:"),
		form,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(500, 300))
	window.ShowAndRun()
}
