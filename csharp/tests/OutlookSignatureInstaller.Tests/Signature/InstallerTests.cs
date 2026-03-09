using System.IO;
using System.IO.Abstractions.TestingHelpers;
using System.Reflection;
using System.Threading.Tasks;
using Xunit;
using OutlookSignatureInstaller.Signature;

namespace OutlookSignatureInstaller.Tests.Signature
{
    /// <summary>
    /// Mirrors pkg/signature/signature_test.go and config_test.go.
    ///
    /// MockFileSystem (System.IO.Abstractions.TestingHelpers) replaces afero.NewMemMapFs(),
    /// providing an in-memory file system without touching the real disk.
    /// </summary>
    public class InstallerTests
    {
        // Stable directory used as the fake exe dir / template base in tests
        private const string FakeExeDir      = @"C:\app";
        private const string FakeTemplateDir = @"C:\app\templates";
        private const string FakeSigDir      = @"C:\signatures";

        // ─── Config priority ─────────────────────────────────────────────────────────

        /// <summary>
        /// Mirrors TestInstall/config priority: user-profile config wins over exe-adjacent config.
        /// </summary>
        [Fact]
        public void LoadConfig_UserConfigExists_TakesPriorityOverExeConfig()
        {
            var fs = new MockFileSystem();

            // Exe-adjacent config
            fs.AddFile(Path.Combine(FakeExeDir, "config.yaml"),
                new MockFileData("template_name: ExeConfig\ntemplate_source: local\nbase_url: ''\n"));

            // User-profile config — should win
            fs.AddFile(ConfigManager.UserConfigPath(),
                new MockFileData("template_name: UserConfig\ntemplate_source: web\nbase_url: http://example.com/\n"));

            var installer = new Installer(FakeTemplateDir, fs, FakeExeDir);
            installer.LoadConfig();

            Assert.Equal("UserConfig", installer.Config?.TemplateName);
            Assert.Equal("web",        installer.Config?.TemplateSource);
        }

        [Fact]
        public void LoadConfig_NoUserConfig_UsesExeAdjacentConfig()
        {
            var fs = new MockFileSystem();
            fs.AddFile(Path.Combine(FakeExeDir, "config.yaml"),
                new MockFileData("template_name: ExeConfig\ntemplate_source: local\nbase_url: ''\n"));

            var installer = new Installer(FakeTemplateDir, fs, FakeExeDir);
            installer.LoadConfig();

            Assert.Equal("ExeConfig", installer.Config?.TemplateName);
            Assert.Equal("local",     installer.Config?.TemplateSource);
        }

        [Fact]
        public void LoadConfig_UserConfigEmptyTemplateName_FallsBackToExeConfig()
        {
            var fs = new MockFileSystem();

            // User config present but TemplateName is empty → should fall through
            fs.AddFile(ConfigManager.UserConfigPath(),
                new MockFileData("template_name: ''\ntemplate_source: web\nbase_url: http://example.com/\n"));

            fs.AddFile(Path.Combine(FakeExeDir, "config.yaml"),
                new MockFileData("template_name: ExeConfig\ntemplate_source: local\nbase_url: ''\n"));

            var installer = new Installer(FakeTemplateDir, fs, FakeExeDir);
            installer.LoadConfig();

            Assert.Equal("ExeConfig", installer.Config?.TemplateName);
        }

        // ─── Install (local mode) ────────────────────────────────────────────────────

        [Fact]
        public async Task InstallAsync_LocalTemplate_WritesBothSignatureFiles()
        {
            var fs = new MockFileSystem();

            // Template files
            fs.AddFile(Path.Combine(FakeTemplateDir, "Standard.htm"),
                new MockFileData("<html><body>Hello {{ .Name }}, {{ .Email }}</body></html>"));
            fs.AddFile(Path.Combine(FakeTemplateDir, "Standard.txt"),
                new MockFileData("Hello {{ .Name }}, {{ .Email }}"));

            // Config
            fs.AddFile(Path.Combine(FakeExeDir, "config.yaml"),
                new MockFileData("template_name: Standard\ntemplate_source: local\nbase_url: ''\n"));

            var installer = new Installer(FakeTemplateDir, fs, FakeExeDir);
            installer.SetSignatureDirectory(FakeSigDir);
            installer.LoadConfig();

            var data = new SignatureData
            {
                Name         = "John Doe",
                Title        = "Engineer",
                Email        = "john@example.com",
                PhoneDisplay = "+49 211 123456",
                PhoneLink    = "+49211123456",
            };

            await installer.InstallAsync(data);

            string htmPath = Path.Combine(FakeSigDir, "Standard.htm");
            string txtPath = Path.Combine(FakeSigDir, "Standard.txt");

            Assert.True(fs.File.Exists(htmPath));
            Assert.True(fs.File.Exists(txtPath));

            string htmContent = fs.File.ReadAllText(htmPath);
            Assert.Contains("John Doe",          htmContent);
            Assert.Contains("john@example.com",  htmContent);
            // Placeholders must be fully replaced
            Assert.DoesNotContain("{{ .Name }}",  htmContent);
            Assert.DoesNotContain("{{ .Email }}", htmContent);

            string txtContent = fs.File.ReadAllText(txtPath);
            Assert.Contains("John Doe",         txtContent);
            Assert.Contains("john@example.com", txtContent);
        }

        [Fact]
        public async Task InstallAsync_HtmlTemplate_HtmlEncodesSpecialChars()
        {
            var fs = new MockFileSystem();

            fs.AddFile(Path.Combine(FakeTemplateDir, "Standard.htm"),
                new MockFileData("<p>{{ .Name }}</p>"));
            fs.AddFile(Path.Combine(FakeTemplateDir, "Standard.txt"),
                new MockFileData("{{ .Name }}"));
            fs.AddFile(Path.Combine(FakeExeDir, "config.yaml"),
                new MockFileData("template_name: Standard\ntemplate_source: local\nbase_url: ''\n"));

            var installer = new Installer(FakeTemplateDir, fs, FakeExeDir);
            installer.SetSignatureDirectory(FakeSigDir);
            installer.LoadConfig();

            // Name with HTML-special characters — should be encoded in .htm, raw in .txt
            var data = new SignatureData { Name = "A & B <Test>", Email = "a@b.com" };
            await installer.InstallAsync(data);

            string htmContent = fs.File.ReadAllText(Path.Combine(FakeSigDir, "Standard.htm"));
            Assert.Contains("A &amp; B &lt;Test&gt;", htmContent);

            string txtContent = fs.File.ReadAllText(Path.Combine(FakeSigDir, "Standard.txt"));
            Assert.Contains("A & B <Test>", txtContent); // no encoding for .txt
        }

        // ─── SaveUserConfig / LoadUserConfig ─────────────────────────────────────────

        [Fact]
        public void SaveAndLoadUserConfig_RoundTrip_PreservesAllFields()
        {
            var fs = new MockFileSystem();

            var cfg = new Config
            {
                TemplateName   = "MyTemplate",
                TemplateSource = "web",
                BaseUrl        = "http://example.com/",
            };

            ConfigManager.SaveUserConfig(fs, cfg);

            var loaded = ConfigManager.LoadUserConfig(fs);
            Assert.NotNull(loaded);
            Assert.Equal("MyTemplate",          loaded!.TemplateName);
            Assert.Equal("web",                 loaded.TemplateSource);
            Assert.Equal("http://example.com/", loaded.BaseUrl);
        }

        [Fact]
        public void LoadUserConfig_FileMissing_ReturnsNull()
        {
            var fs = new MockFileSystem(); // empty FS — no config file
            Assert.Null(ConfigManager.LoadUserConfig(fs));
        }
    }
}
