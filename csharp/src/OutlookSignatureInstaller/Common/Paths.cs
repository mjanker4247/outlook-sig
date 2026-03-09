using System;
using System.IO;
using System.Reflection;

namespace OutlookSignatureInstaller.Common
{
    /// <summary>
    /// Mirrors pkg/common/paths.go: resolves the templates directory relative to the executable.
    /// </summary>
    public static class Paths
    {
        /// <summary>
        /// Returns the path to the templates/ directory next to the application executable.
        /// Mirrors GetTemplateBase() in Go, which uses os.Executable() + filepath.Dir().
        ///
        /// C# change: Assembly.GetExecutingAssembly().Location is used instead of os.Executable().
        /// In .NET, single-file publish may return an empty string — handle that gracefully.
        /// </summary>
        public static string GetTemplateBase()
        {
            string? exePath = Assembly.GetExecutingAssembly().Location;
            if (string.IsNullOrEmpty(exePath))
                exePath = Environment.ProcessPath
                    ?? throw new InvalidOperationException("Failed to determine executable path.");

            string exeDir = Path.GetDirectoryName(exePath)
                ?? throw new InvalidOperationException("Failed to determine executable directory.");

            return Path.Combine(exeDir, "templates");
        }
    }
}
