using System.Runtime.InteropServices;

namespace OutlookSignatureInstaller
{
    /// <summary>
    /// Mirrors console_windows.go: attaches to the parent process's console when CLI
    /// arguments are present. The binary is built as WinExe (windowsgui) so no console
    /// opens by default — this enables CLI usage from cmd.exe / PowerShell.
    /// </summary>
    internal static class ConsoleHelper
    {
        [DllImport("kernel32.dll", SetLastError = true)]
        private static extern bool AttachConsole(uint dwProcessId);

        // ATTACH_PARENT_PROCESS = (DWORD)-1, mirrors the Go constant ^uintptr(0)
        private const uint AttachParentProcess = unchecked((uint)-1);

        public static void AttachIfNeeded(string[] args)
        {
            if (args.Length == 0)
                return;

            // Attach to parent console (cmd.exe, PowerShell, etc.).
            // If no parent console exists (e.g. launched from Explorer with args),
            // AttachConsole returns false and we do nothing — no new window is opened.
            AttachConsole(AttachParentProcess);
        }
    }
}
