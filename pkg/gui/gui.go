//go:build !windows
// +build !windows

package gui

import (
	"fmt"

	"outlook-signature/pkg/common"
	"outlook-signature/pkg/signature"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// createValidatedEntry creates an entry widget with validation
func createValidatedEntry(placeholder string, validator func(string) error) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	entry.Validator = validator
	entry.OnChanged = func(s string) {
		if err := entry.Validate(); err != nil {
			entry.SetValidationError(err)
		} else {
			entry.SetValidationError(nil)
		}
		entry.Refresh()
	}
	return entry
}

// ShowGUI displays the signature installer GUI
func ShowGUI() {
	myApp := app.New()
	window := myApp.NewWindow("Outlook Signature Installer")

	// Create form fields with validation
	nameEntry := createValidatedEntry("Your full name", common.ValidateName)
	titleEntry := createValidatedEntry("Your profession or title (optional)", common.ValidateName)
	emailEntry := createValidatedEntry("Your email address", common.ValidateEmail)
	phoneEntry := createValidatedEntry("Your phone number", common.ValidatePhoneNumber)

	// Get template base directory
	templateBase, err := common.GetTemplateBase()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to find templates: %v", err), window)
		return
	}

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Title", Widget: titleEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Phone", Widget: phoneEntry},
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
				dialog.ShowError(err, window)
				return
			}

			// Create signature data
			data := signature.Data{
				Name:         nameEntry.Text,
				Title:        titleEntry.Text,
				Email:        emailEntry.Text,
				PhoneDisplay: phoneDisplay,
				PhoneLink:    phoneLink,
			}

			// Install signature
			installer := signature.NewInstaller(templateBase)
			if err := installer.LoadConfig(); err != nil {
				dialog.ShowError(fmt.Errorf("failed to load configuration: %v", err), window)
				return
			}

			if err := installer.Install(data); err != nil {
				dialog.ShowError(fmt.Errorf("failed to install signature: %v", err), window)
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
	window.Resize(fyne.NewSize(500, 250))
	window.ShowAndRun()
}
