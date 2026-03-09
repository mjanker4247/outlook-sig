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
	"github.com/spf13/afero"
)

// createValidatedEntry creates an entry widget with real-time validation.
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

// ShowGUI displays the tabbed Outlook Signature Installer window.
// It performs a best-effort config load at startup to pre-populate the Settings
// tab. Template loading is deferred to submission so the window is visible before
// any error dialogs are displayed.
func ShowGUI() {
	myApp := app.New()
	window := myApp.NewWindow("Outlook Signature Installer")

	// Best-effort startup config load — errors are non-fatal here because the
	// window is not yet visible and we cannot show a dialog yet. The Settings tab
	// will surface any issue when the user tries to save.
	cfg, _ := signature.LoadUserConfig(afero.NewOsFs())
	if cfg == nil {
		cfg = &signature.Config{TemplateName: "Standard", TemplateSource: "local"}
	}

	tabs := container.NewAppTabs(
		buildSignatureTab(window, cfg),
		buildSettingsTab(window, cfg),
	)

	window.SetContent(tabs)
	window.Resize(fyne.NewSize(520, 380))
	window.ShowAndRun()
}

// buildSignatureTab constructs the "Signature" install tab.
// The shared cfg pointer is read at submit time so it reflects any settings
// the user may have saved in the Settings tab during the same session.
func buildSignatureTab(window fyne.Window, cfg *signature.Config) *container.TabItem {
	nameEntry := createValidatedEntry("Your full name", common.ValidateName)
	titleEntry := createValidatedEntry("Your profession or title (optional)", common.ValidateTitle)
	emailEntry := createValidatedEntry("Your email address", common.ValidateEmail)
	phoneEntry := createValidatedEntry("Your phone number", common.ValidatePhoneNumber)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Title", Widget: titleEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Phone", Widget: phoneEntry},
		},
		OnSubmit: func() {
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

			// Resolve template directory at submission time so the window is
			// already visible when any error dialog is shown.
			templateBase, err := common.GetTemplateBase()
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to find templates: %v", err), window)
				return
			}

			phoneDisplay, phoneLink, err := common.FormatPhoneNumber(phoneEntry.Text, "DE")
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			data := signature.Data{
				Name:         nameEntry.Text,
				Title:        titleEntry.Text,
				Email:        emailEntry.Text,
				PhoneDisplay: phoneDisplay,
				PhoneLink:    phoneLink,
			}

			installer := signature.NewInstaller(templateBase)
			// Seed the installer with any in-session config changes so the user
			// does not have to restart after editing settings.
			installer.Config = cfg
			if err := installer.Install(data); err != nil {
				dialog.ShowError(fmt.Errorf("failed to install signature: %v", err), window)
				return
			}

			dialog.ShowInformation("Success", "Signature installed successfully!", window)
		},
	}

	content := container.NewVBox(
		widget.NewLabel("Enter your signature details:"),
		form,
	)
	return container.NewTabItem("Signature", content)
}

// buildSettingsTab constructs the "Settings" configuration tab.
// Fields are pre-populated from the shared cfg pointer. When the user saves,
// cfg is updated in-place and written to the user-profile location so the next
// launch picks up the changes automatically.
func buildSettingsTab(window fyne.Window, cfg *signature.Config) *container.TabItem {
	templateNameEntry := widget.NewEntry()
	templateNameEntry.SetText(cfg.TemplateName)
	templateNameEntry.SetPlaceHolder("e.g. Standard")

	baseURLEntry := widget.NewEntry()
	baseURLEntry.SetText(cfg.BaseURL)
	baseURLEntry.SetPlaceHolder("e.g. http://server/templates/")

	// Disable the URL field when source is "local".
	applySourceState := func(source string) {
		if source == "web" {
			baseURLEntry.Enable()
		} else {
			baseURLEntry.Disable()
		}
	}

	sourceSelect := widget.NewSelect([]string{"local", "web"}, applySourceState)
	sourceSelect.SetSelected(cfg.TemplateSource)
	applySourceState(cfg.TemplateSource) // apply initial state

	saveBtn := widget.NewButton("Save Settings", func() {
		name := templateNameEntry.Text
		if name == "" {
			dialog.ShowError(fmt.Errorf("Template Name: cannot be empty"), window)
			return
		}
		if err := common.ValidateSignatureName(name); err != nil {
			dialog.ShowError(err, window)
			return
		}

		source := sourceSelect.Selected
		rawURL := baseURLEntry.Text
		if source == "web" {
			if rawURL == "" {
				dialog.ShowError(fmt.Errorf("Base URL: cannot be empty when source is \"web\""), window)
				return
			}
			if err := common.ValidateURL(rawURL); err != nil {
				dialog.ShowError(err, window)
				return
			}
		}

		// Update shared config in-place so the Signature tab picks it up immediately.
		cfg.TemplateName = name
		cfg.TemplateSource = source
		cfg.BaseURL = rawURL

		if err := signature.SaveUserConfig(afero.NewOsFs(), cfg); err != nil {
			dialog.ShowError(fmt.Errorf("failed to save settings: %v", err), window)
			return
		}

		dialog.ShowInformation("Saved", "Settings saved successfully!", window)
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Template Name", Widget: templateNameEntry},
			{Text: "Template Source", Widget: sourceSelect},
			{Text: "Base URL", Widget: baseURLEntry},
		},
	}

	content := container.NewVBox(
		widget.NewLabel("Application Settings"),
		form,
		saveBtn,
	)
	return container.NewTabItem("Settings", content)
}
