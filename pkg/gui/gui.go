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

// ShowGUI displays the signature installer GUI
func ShowGUI() {
	myApp := app.New()
	window := myApp.NewWindow("Outlook Signature Installer")

	// Create form fields
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Your full name")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Your email address")

	phoneEntry := widget.NewEntry()
	phoneEntry.SetPlaceHolder("Your phone number")

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
			if nameEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("name cannot be empty"), window)
				return
			}

			if err := common.ValidateEmail(emailEntry.Text); err != nil {
				dialog.ShowError(err, window)
				return
			}

			if err := common.ValidatePhoneNumber(phoneEntry.Text); err != nil {
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
