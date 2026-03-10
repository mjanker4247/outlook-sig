using OutlookSignatureInstaller.Common;
using System;
using System.IO;
using System.IO.Abstractions;
using System.Net.Http;
using System.Reflection;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;

namespace OutlookSignatureInstaller.Signature
{
    /// <summary>
    /// Mirrors pkg/signature/signature.go: the Data struct that carries user input into templates.
    /// </summary>
    public class SignatureData
    {
        public string Name { get; set; } = string.Empty;
        public string Title { get; set; } = string.Empty;
        public string Email { get; set; } = string.Empty;
        /// <summary>Human-readable phone, e.g. "+49 211 123456" (INTERNATIONAL format).</summary>
        public string PhoneDisplay { get; set; } = string.Empty;
        /// <summary>E164 phone for tel: links, e.g. "+49211123456".</summary>
        public string PhoneLink { get; set; } = string.Empty;
    }

    /// <summary>
    /// Mirrors pkg/signature/signature.go: core signature installation logic.
    ///
    /// Key C# adaptations:
    ///   - afero.Fs → IFileSystem (System.IO.Abstractions) for testable file operations.
    ///   - *slog.Logger removed; errors are returned via exceptions (idiomatic C#).
    ///   - Go's html/template engine replaced with regex-based placeholder substitution
    ///     (same approach Go uses for .txt; HTML values are HtmlEncoded for safety).
    ///   - Async I/O for HTTP downloads uses Task/await instead of Go's blocking HTTP client.
    /// </summary>
    public class Installer
    {
        private readonly string _templateBase;
        private readonly IFileSystem _fs;
        private readonly string _exeDir;
        private string? _sigDirOverride;

        /// <summary>Loaded configuration. Populated by LoadConfig().</summary>
        public Config? Config { get; private set; }

        /// <param name="templateBase">Path to the templates directory.</param>
        /// <param name="fs">File system abstraction (pass null for real FS, MockFileSystem in tests).</param>
        /// <param name="exeDir">
        ///     Directory containing the executable and config.yaml.
        ///     Defaults to the assembly location; injectable for unit tests.
        /// </param>
        public Installer(string templateBase, IFileSystem? fs = null, string? exeDir = null)
        {
            _templateBase = templateBase;
            _fs = fs ?? new FileSystem();
            _exeDir = exeDir
                ?? Path.GetDirectoryName(Assembly.GetExecutingAssembly().Location)
                ?? ".";
        }

        /// <summary>
        /// Allows tests to override the Outlook signature directory,
        /// mirroring the sigDir field in the Go Installer struct.
        /// </summary>
        public void SetSignatureDirectory(string sigDir) => _sigDirOverride = sigDir;

        /// <summary>
        /// Mirrors LoadConfig(): loads configuration with user-profile override priority.
        ///   1. User-profile config (if present and has a non-empty TemplateName)
        ///   2. Exe-adjacent config.yaml as fallback
        /// </summary>
        public void LoadConfig()
        {
            var userCfg = ConfigManager.LoadUserConfig(_fs);
            if (userCfg != null && !string.IsNullOrWhiteSpace(userCfg.TemplateName))
            {
                Config = userCfg;
                return;
            }

            string configPath = Path.Combine(_exeDir, "config.yaml");
            if (!_fs.File.Exists(configPath))
                throw new FileNotFoundException($"Config file not found: {configPath}");

            string yaml = _fs.File.ReadAllText(configPath);
            Config = ConfigManager.DeserializeConfig(yaml);
        }

        /// <summary>
        /// Mirrors DownloadWebTemplates(): downloads .htm and .txt template files from
        /// Config.BaseUrl + Config.TemplateName. Enforces FileSizeLimit (50 MB).
        ///
        /// C# change: async/await instead of Go's blocking http.Get + io.LimitReader.
        /// </summary>
        public async Task DownloadWebTemplatesAsync()
        {
            if (Config == null)
                throw new InvalidOperationException("Config not loaded. Call LoadConfig() first.");

            ValidateAndSanitizeTemplateName(Config.TemplateName);

            string baseUrl = Config.BaseUrl.TrimEnd('/') + '/';
            using var client = new HttpClient();

            foreach (string ext in new[] { ".htm", ".txt" })
            {
                string url = baseUrl + Config.TemplateName + ext;
                var response = await client.GetAsync(url);
                response.EnsureSuccessStatusCode();

                byte[] content = await response.Content.ReadAsByteArrayAsync();
                if (content.Length > Validation.FileSizeLimit)
                    throw new InvalidOperationException(
                        $"Template exceeds {Validation.FileSizeLimit / 1024 / 1024} MB limit: {url}");

                _fs.Directory.CreateDirectory(_templateBase);
                string destPath = Path.Combine(_templateBase, Config.TemplateName + ext);
                _fs.File.WriteAllBytes(destPath, content);
            }
        }

        /// <summary>
        /// Mirrors Install(): main installation orchestrator.
        /// </summary>
        public async Task InstallAsync(SignatureData data)
        {
            if (Config == null)
                throw new InvalidOperationException("Config not loaded. Call LoadConfig() first.");

            ValidateTemplateBase();

            if (Config.TemplateSource == "web")
                await DownloadWebTemplatesAsync();

            ValidateAndSanitizeTemplateName(Config.TemplateName);
            string sigDir = GetSignatureDirectory();
            _fs.Directory.CreateDirectory(sigDir);

            await InstallSignatureFilesAsync(Config.TemplateName, sigDir, data);
        }

        // --- Private helpers ---

        private void ValidateTemplateBase()
        {
            if (!_fs.Directory.Exists(_templateBase))
                throw new DirectoryNotFoundException(
                    $"Template directory not found: {_templateBase}");
        }

        /// <summary>
        /// Mirrors validateAndSanitizeTemplateName(): prevents path traversal attacks
        /// by ensuring the resolved template path stays within _templateBase.
        /// Mirrors Go's filepath.Rel() check.
        /// </summary>
        private void ValidateAndSanitizeTemplateName(string name)
        {
            if (string.IsNullOrWhiteSpace(name))
                throw new ArgumentException("Template name must not be empty.");

            // Reject names with path separators or parent-directory components
            if (name.Contains('/') || name.Contains('\\') || name.Contains(".."))
                throw new ArgumentException(
                    $"Invalid template name (path traversal detected): {name}");

            // Verify the resolved .htm path stays under _templateBase
            string resolved = Path.GetFullPath(Path.Combine(_templateBase, name + ".htm"));
            string baseNorm = Path.GetFullPath(_templateBase);

            if (!resolved.StartsWith(baseNorm + Path.DirectorySeparatorChar, StringComparison.OrdinalIgnoreCase)
                && !resolved.StartsWith(baseNorm + '/', StringComparison.OrdinalIgnoreCase))
            {
                throw new ArgumentException(
                    $"Template name resolves outside template directory: {name}");
            }
        }

        private string GetSignatureDirectory() =>
            _sigDirOverride ?? GetOutlookSignatureDir();

        /// <summary>
        /// Mirrors GetOutlookSignatureDir(): returns the platform-specific Outlook signature path.
        /// Windows: %APPDATA%\Microsoft\Signatures\
        /// macOS:   ~/Library/Group Containers/UBF8T346G9.Office/.../Signatures/
        /// </summary>
        public static string GetOutlookSignatureDir()
        {
            if (OperatingSystem.IsWindows())
            {
                string appData = Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData);
                return Path.Combine(appData, "Microsoft", "Signatures");
            }
            else if (OperatingSystem.IsMacOS())
            {
                string home = Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
                return Path.Combine(home, "Library", "Group Containers",
                    "UBF8T346G9.Office", "Outlook", "Outlook 15 Profiles",
                    "Main Profile", "Signatures");
            }
            else
            {
                throw new PlatformNotSupportedException(
                    "Only Windows and macOS are supported.");
            }
        }

        private async Task InstallSignatureFilesAsync(string sigName, string sigDir, SignatureData data)
        {
            foreach (string ext in new[] { ".htm", ".txt" })
                await InstallFileAsync(sigName, sigDir, ext, data);
        }

        private async Task InstallFileAsync(string sigName, string sigDir, string ext, SignatureData data)
        {
            string templatePath = Path.Combine(_templateBase, sigName + ext);
            string destPath = Path.Combine(sigDir, sigName + ext);

            if (ext == ".htm")
                await InstallHtmlFileAsync(templatePath, destPath, data);
            else
                await InstallTextFileAsync(templatePath, destPath, data);
        }

        /// <summary>
        /// Mirrors installHTMLFile(): renders the HTML template by replacing Go-style placeholders.
        ///
        /// Go change: Go uses html/template with automatic HTML escaping. Here we apply
        /// System.Net.WebUtility.HtmlEncode() to each substituted value for equivalent safety.
        /// </summary>
        private async Task InstallHtmlFileAsync(string templatePath, string destPath, SignatureData data)
        {
            string content = ReadWithSizeLimit(templatePath);
            content = ApplyPlaceholders(content, data, htmlEncode: true);
            // WriteAllText matches Go's os.WriteFile with UTF-8 encoding
            await Task.Run(() => _fs.File.WriteAllText(destPath, content, Encoding.UTF8));
        }

        /// <summary>
        /// Mirrors installTextFile(): regex-based placeholder substitution, no HTML encoding.
        /// </summary>
        private async Task InstallTextFileAsync(string templatePath, string destPath, SignatureData data)
        {
            string content = ReadWithSizeLimit(templatePath);
            content = ApplyPlaceholders(content, data, htmlEncode: false);
            await Task.Run(() => _fs.File.WriteAllText(destPath, content, Encoding.UTF8));
        }

        private string ReadWithSizeLimit(string path)
        {
            var info = _fs.FileInfo.New(path);
            if (info.Length > Validation.BufferSizeLimit)
                throw new InvalidOperationException(
                    $"Template file exceeds {Validation.BufferSizeLimit / 1024 / 1024} MB limit: {path}");

            return _fs.File.ReadAllText(path, Encoding.UTF8);
        }

        /// <summary>
        /// Mirrors the placeholders slice in signature.go: replaces all Go-template placeholders
        /// {{ .Name }}, {{ .Title }}, {{ .Email }}, {{ .PhoneDisplay }}, {{ .PhoneLink }}.
        ///
        /// The regex pattern mirrors Go's `{{\s*\.Name\s*}}` etc., so existing template files
        /// require no modification.
        /// </summary>
        private static string ApplyPlaceholders(string content, SignatureData data, bool htmlEncode)
        {
            string Encode(string s) => htmlEncode
                ? System.Net.WebUtility.HtmlEncode(s)
                : s;

            content = Regex.Replace(content, @"\{\{\s*\.Name\s*\}\}", Encode(data.Name));
            content = Regex.Replace(content, @"\{\{\s*\.Title\s*\}\}", Encode(data.Title));
            content = Regex.Replace(content, @"\{\{\s*\.Email\s*\}\}", Encode(data.Email));
            content = Regex.Replace(content, @"\{\{\s*\.PhoneDisplay\s*\}\}", Encode(data.PhoneDisplay));
            content = Regex.Replace(content, @"\{\{\s*\.PhoneLink\s*\}\}", Encode(data.PhoneLink));

            return content;
        }
    }
}
