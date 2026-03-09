using System;
using System.IO.Abstractions;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;
using OutlookSignatureInstaller.Common;
using OutlookSignatureInstaller.Signature;

namespace OutlookSignatureInstaller.Gui
{
    /// <summary>
    /// Mirrors pkg/gui/gui.go: WPF code-behind for the two-tab main window.
    ///
    /// Fyne → WPF mapping:
    ///   fyne.Window       → Window (this class)
    ///   container.AppTabs → TabControl (in XAML)
    ///   widget.Entry      → TextBox
    ///   widget.Label      → TextBlock
    ///   widget.Button     → Button
    ///   widget.Select     → ComboBox
    ///   Real-time validation via Entry.Validator → TextChanged event + error TextBlock
    /// </summary>
    public partial class MainWindow : Window
    {
        private Config? _config;

        public MainWindow()
        {
            InitializeComponent();
            LoadAndPopulateConfig();
        }

        /// <summary>
        /// Mirrors ShowGUI()'s initial config load: reads config and pre-fills the Settings tab.
        /// </summary>
        private void LoadAndPopulateConfig()
        {
            try
            {
                var installer = new Installer(Paths.GetTemplateBase());
                installer.LoadConfig();
                _config = installer.Config;

                if (_config != null)
                {
                    TemplateNameBox.Text = _config.TemplateName;
                    BaseUrlBox.Text = _config.BaseUrl;

                    // Select the matching item in TemplateSourceBox
                    foreach (ComboBoxItem item in TemplateSourceBox.Items)
                    {
                        if ((string)item.Content == _config.TemplateSource)
                        {
                            TemplateSourceBox.SelectedItem = item;
                            break;
                        }
                    }
                }

                if (TemplateSourceBox.SelectedItem == null)
                    TemplateSourceBox.SelectedIndex = 0;
            }
            catch
            {
                // If config is missing or unreadable, start with defaults — same as Go GUI.
                _config = new Config();
                TemplateSourceBox.SelectedIndex = 0;
            }
        }

        // ─── Real-time validation handlers ─────────────────────────────────────────
        // Mirrors createValidatedEntry()'s inline validator callback in gui.go.

        private void OnNameChanged(object sender, TextChangedEventArgs e) =>
            ShowError(NameError, Validation.ValidateName(NameBox.Text)?.Message);

        private void OnEmailChanged(object sender, TextChangedEventArgs e) =>
            ShowError(EmailError, Validation.ValidateEmail(EmailBox.Text)?.Message);

        private void OnPhoneChanged(object sender, TextChangedEventArgs e) =>
            ShowError(PhoneError, Validation.ValidatePhoneNumber(PhoneBox.Text)?.Message);

        private void OnBaseUrlChanged(object sender, TextChangedEventArgs e) =>
            ShowError(UrlError, Validation.ValidateURL(BaseUrlBox.Text)?.Message);

        // ─── Install button ──────────────────────────────────────────────────────────

        /// <summary>
        /// Mirrors the Install button handler in buildSignatureTab(): validates all fields,
        /// builds SignatureData, calls Installer.Install(), and shows success/error feedback.
        /// </summary>
        private async void OnInstallClicked(object sender, RoutedEventArgs e)
        {
            // Final validation pass before installing
            var nameErr  = Validation.ValidateName(NameBox.Text);
            var emailErr = Validation.ValidateEmail(EmailBox.Text);
            var phoneErr = Validation.ValidatePhoneNumber(PhoneBox.Text);

            ShowError(NameError,  nameErr?.Message);
            ShowError(EmailError, emailErr?.Message);
            ShowError(PhoneError, phoneErr?.Message);

            if (nameErr != null || emailErr != null || phoneErr != null)
                return;

            InstallButton.IsEnabled = false;
            SetStatus(StatusLabel, "Installing...", Brushes.Gray);

            try
            {
                var (phoneDisplay, phoneLink) = Validation.FormatPhoneNumber(PhoneBox.Text);

                var data = new SignatureData
                {
                    Name         = NameBox.Text,
                    Title        = TitleBox.Text,
                    Email        = EmailBox.Text,
                    PhoneDisplay = phoneDisplay,
                    PhoneLink    = phoneLink,
                };

                var installer = new Installer(Paths.GetTemplateBase());
                installer.LoadConfig();
                await installer.InstallAsync(data);

                SetStatus(StatusLabel, "Signature installed successfully.", Brushes.Green);
            }
            catch (Exception ex)
            {
                SetStatus(StatusLabel, $"Error: {ex.Message}", Brushes.Red);
            }
            finally
            {
                InstallButton.IsEnabled = true;
            }
        }

        // ─── Settings save button ────────────────────────────────────────────────────

        /// <summary>
        /// Mirrors the Save button in buildSettingsTab(): validates URL, serialises and
        /// writes the config via ConfigManager.SaveUserConfig().
        /// </summary>
        private void OnSaveSettingsClicked(object sender, RoutedEventArgs e)
        {
            var urlErr = Validation.ValidateURL(BaseUrlBox.Text);
            ShowError(UrlError, urlErr?.Message);
            if (urlErr != null) return;

            try
            {
                var cfg = new Config
                {
                    TemplateName   = TemplateNameBox.Text,
                    TemplateSource = (TemplateSourceBox.SelectedItem as ComboBoxItem)
                                         ?.Content?.ToString() ?? "local",
                    BaseUrl        = BaseUrlBox.Text,
                };

                ConfigManager.SaveUserConfig(new FileSystem(), cfg);
                _config = cfg;

                SetStatus(SettingsStatusLabel, "Settings saved.", Brushes.Green);
            }
            catch (Exception ex)
            {
                SetStatus(SettingsStatusLabel, $"Error saving settings: {ex.Message}", Brushes.Red);
            }
        }

        // ─── UI helpers ─────────────────────────────────────────────────────────────

        private static void ShowError(TextBlock label, string? message)
        {
            if (string.IsNullOrEmpty(message))
            {
                label.Text       = string.Empty;
                label.Visibility = Visibility.Collapsed;
            }
            else
            {
                label.Text       = message;
                label.Visibility = Visibility.Visible;
            }
        }

        private static void SetStatus(TextBlock label, string message, Brush colour)
        {
            label.Text       = message;
            label.Foreground = colour;
        }
    }
}
