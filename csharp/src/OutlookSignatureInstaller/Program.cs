using OutlookSignatureInstaller.Cli;
using OutlookSignatureInstaller.Gui;
using System;
using System.Windows;

namespace OutlookSignatureInstaller
{
    /// <summary>
    /// Entry point. Mirrors cmd/signature-installer/main.go + console_windows.go:
    /// - No args (or --gui/-g): launch WPF GUI silently (no console window).
    /// - With CLI args: attach to parent console and run CLI mode.
    ///
    /// WPF change: [STAThread] is required for WPF's single-threaded apartment COM model.
    /// </summary>
    class Program
    {
        [STAThread]
        static int Main(string[] args)
        {
            // Mirrors console_windows.go: attach parent process console when CLI args are present.
            ConsoleHelper.AttachIfNeeded(args);

            bool launchGui = args.Length == 0
                || Array.Exists(args, a => a == "--gui" || a == "-g");

            if (launchGui)
            {
                // WPF requires an Application instance on the STA thread.
                var app = new Application();
                var window = new MainWindow();
                app.Run(window); // blocks until window closes
                return 0;
            }
            else
            {
                return CliApp.Run(args);
            }
        }
    }
}
