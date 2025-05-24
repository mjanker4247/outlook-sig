package gui

import (
	"fmt"
	"path/filepath"

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

	templateEntry := widget.NewEntry()
	templateEntry.SetPlaceHolder("Template name")
	templateEntry.SetText("OutlookSignature")

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Phone", Widget: phoneEntry},
			{Text: "Template", Widget: templateEntry},
		},
		OnSubmit: func() {
			// Get executable directory for templates
			exeDir, err := filepath.Abs(".")
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to get executable path: %v", err), window)
				return
			}
			templateBase := filepath.Join(exeDir, "templates")

			// Create signature data
			data := signature.Data{
				Name:         nameEntry.Text,
				Email:        emailEntry.Text,
				PhoneDisplay: phoneEntry.Text,
				PhoneLink:    phoneEntry.Text,
			}

			// Install signature
			installer := signature.NewInstaller(templateBase)
			err = installer.Install(data, templateEntry.Text)
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
	window.Resize(fyne.NewSize(400, 300))
	window.ShowAndRun()
}
