using System;
using System.IO;
using System.IO.Abstractions;
using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;

namespace OutlookSignatureInstaller.Signature
{
    /// <summary>
    /// Mirrors pkg/signature/config.go: the config.yaml structure.
    /// YamlMember aliases mirror the yaml struct tags in Go (template_name, template_source, base_url).
    /// </summary>
    public class Config
    {
        [YamlMember(Alias = "template_name")]
        public string TemplateName { get; set; } = "Standard";

        /// <summary>"local" or "web" — mirrors TemplateSource in Go.</summary>
        [YamlMember(Alias = "template_source")]
        public string TemplateSource { get; set; } = "local";

        /// <summary>Base URL used when TemplateSource == "web".</summary>
        [YamlMember(Alias = "base_url")]
        public string BaseUrl { get; set; } = string.Empty;
    }

    /// <summary>
    /// Mirrors the config loading/saving functions in pkg/signature/config.go.
    ///
    /// C# change: static class replaces Go package-level functions.
    /// IFileSystem (System.IO.Abstractions) replaces afero.Fs for testability.
    /// </summary>
    public static class ConfigManager
    {
        private const string AppName = "OutlookSignatureInstaller";
        private const string ConfigFileName = "config.yaml";

        /// <summary>
        /// Mirrors UserConfigPath(): returns the platform-specific user-profile config path.
        /// Windows: %APPDATA%\OutlookSignatureInstaller\config.yaml
        /// macOS:   ~/Library/Application Support/OutlookSignatureInstaller/config.yaml
        /// </summary>
        public static string UserConfigPath()
        {
            string baseDir;
            if (OperatingSystem.IsWindows())
            {
                baseDir = Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData);
            }
            else if (OperatingSystem.IsMacOS())
            {
                string home = Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
                baseDir = Path.Combine(home, "Library", "Application Support");
            }
            else
            {
                // Linux fallback: use XDG_CONFIG_HOME or ~/.config
                string xdg = Environment.GetEnvironmentVariable("XDG_CONFIG_HOME")
                    ?? Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".config");
                baseDir = xdg;
            }

            return Path.Combine(baseDir, AppName, ConfigFileName);
        }

        /// <summary>
        /// Mirrors LoadUserConfig(): loads the user-profile config.
        /// Returns null if the file does not exist or cannot be parsed — same as Go returning nil.
        /// </summary>
        public static Config? LoadUserConfig(IFileSystem fs)
        {
            string path = UserConfigPath();
            if (!fs.File.Exists(path))
                return null;

            try
            {
                string yaml = fs.File.ReadAllText(path);
                return DeserializeConfig(yaml);
            }
            catch
            {
                // Corrupt or unreadable config — fall through to exe-adjacent config, same as Go.
                return null;
            }
        }

        /// <summary>
        /// Mirrors SaveUserConfig(): writes the config to the user-profile path,
        /// creating the directory if it doesn't exist.
        /// </summary>
        public static void SaveUserConfig(IFileSystem fs, Config cfg)
        {
            string path = UserConfigPath();
            string dir = fs.Path.GetDirectoryName(path)
                ?? throw new InvalidOperationException("Cannot determine config directory.");
            fs.Directory.CreateDirectory(dir);
            fs.File.WriteAllText(path, SerializeConfig(cfg));
        }

        /// <summary>Deserializes a YAML string into a Config object.</summary>
        public static Config DeserializeConfig(string yaml)
        {
            var deserializer = new DeserializerBuilder()
                .WithNamingConvention(UnderscoredNamingConvention.Instance)
                .IgnoreUnmatchedProperties()
                .Build();
            return deserializer.Deserialize<Config>(yaml) ?? new Config();
        }

        /// <summary>Serializes a Config object to YAML.</summary>
        public static string SerializeConfig(Config cfg)
        {
            var serializer = new SerializerBuilder()
                .WithNamingConvention(UnderscoredNamingConvention.Instance)
                .Build();
            return serializer.Serialize(cfg);
        }
    }
}
