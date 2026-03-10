using OutlookSignatureInstaller.Common;
using OutlookSignatureInstaller.Signature;
using System;
using System.Collections.Generic;

namespace OutlookSignatureInstaller.Cli
{
    /// <summary>
    /// Mirrors pkg/cli/cli.go: command-line interface for the signature installer.
    ///
    /// C# change: urfave/cli replaced with simple manual arg parsing.
    /// The flag API is intentionally close to the original:
    ///   --name/-n, --title/-t, --email/-e, --phone/-p, --template-source/-s, --gui/-g
    /// </summary>
    public static class CliApp
    {
        /// <summary>
        /// Mirrors App() + runCLIInstallation(): entry point for CLI mode.
        /// Returns exit code (0 = success, 1 = error).
        /// </summary>
        public static int Run(string[] args)
        {
            if (ContainsFlag(args, "--help", "-h"))
            {
                PrintHelp();
                return 0;
            }

            try
            {
                var flags = ParseArgs(args);
                var data = GetUserInput(flags);
                RunInstallation(flags, data);
                Console.WriteLine("Signature installed successfully.");
                return 0;
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error: {ex.Message}");
                return 1;
            }
        }

        // --- Arg parsing ---

        private static Dictionary<string, string> ParseArgs(string[] args)
        {
            var result = new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);
            for (int i = 0; i < args.Length - 1; i++)
            {
                string flag = args[i];
                string value = args[i + 1];
                // Skip if next token is itself a flag
                if (value.StartsWith('-')) continue;

                switch (flag)
                {
                    case "--name": case "-n": result["name"] = value; i++; break;
                    case "--title": case "-t": result["title"] = value; i++; break;
                    case "--email": case "-e": result["email"] = value; i++; break;
                    case "--phone": case "-p": result["phone"] = value; i++; break;
                    case "--template-source": case "-s": result["template-source"] = value; i++; break;
                }
            }
            return result;
        }

        private static bool ContainsFlag(string[] args, params string[] flags)
        {
            foreach (string arg in args)
                foreach (string f in flags)
                    if (arg == f) return true;
            return false;
        }

        // --- Input collection + validation ---

        /// <summary>
        /// Mirrors getUserInput(): collects all fields, prompting for any that are missing.
        /// Validates each field; throws ArgumentException on invalid input.
        /// </summary>
        private static SignatureData GetUserInput(Dictionary<string, string> flags)
        {
            string name = GetOrPrompt(flags.GetValueOrDefault("name"), "Enter your full name: ");
            ThrowIfInvalid(Validation.ValidateName(name));

            // Title is intentionally permissive — no validation, same as Go.
            string title = GetOrPrompt(flags.GetValueOrDefault("title"),
                "Enter your job title (optional, press Enter to skip): ");

            string email = GetOrPrompt(flags.GetValueOrDefault("email"), "Enter your email address: ");
            ThrowIfInvalid(Validation.ValidateEmail(email));

            string phone = GetOrPrompt(flags.GetValueOrDefault("phone"), "Enter your phone number: ");
            ThrowIfInvalid(Validation.ValidatePhoneNumber(phone));

            var (phoneDisplay, phoneLink) = Validation.FormatPhoneNumber(phone);

            return new SignatureData
            {
                Name = name,
                Title = title,
                Email = email,
                PhoneDisplay = phoneDisplay,
                PhoneLink = phoneLink,
            };
        }

        private static void RunInstallation(Dictionary<string, string> flags, SignatureData data)
        {
            string templateBase = Paths.GetTemplateBase();
            var installer = new Installer(templateBase);
            installer.LoadConfig();

            // --template-source overrides config.yaml, same as the Go CLI flag handling.
            if (flags.TryGetValue("template-source", out string? source)
                && !string.IsNullOrEmpty(source)
                && installer.Config != null)
            {
                installer.Config.TemplateSource = source;
            }

            installer.InstallAsync(data).GetAwaiter().GetResult();
        }

        // --- Helpers ---

        /// <summary>
        /// Mirrors getOrPrompt(): returns the value if non-empty, otherwise prompts the user.
        /// </summary>
        private static string GetOrPrompt(string? value, string prompt)
        {
            if (!string.IsNullOrWhiteSpace(value))
                return value;

            Console.Write(prompt);
            return Console.ReadLine() ?? string.Empty;
        }

        private static void ThrowIfInvalid(ValidationError? err)
        {
            if (err != null)
                throw new ArgumentException(err.Message);
        }

        private static void PrintHelp()
        {
            Console.WriteLine("Outlook Signature Installer");
            Console.WriteLine();
            Console.WriteLine("Usage: OutlookSignatureInstaller [options]");
            Console.WriteLine();
            Console.WriteLine("Options:");
            Console.WriteLine("  --gui, -g                  Launch GUI mode (default when no args)");
            Console.WriteLine("  --name, -n <name>          Your full name");
            Console.WriteLine("  --title, -t <title>        Your job title (optional)");
            Console.WriteLine("  --email, -e <email>        Your email address");
            Console.WriteLine("  --phone, -p <phone>        Your phone number");
            Console.WriteLine("  --template-source, -s      Template source: 'local' or 'web'");
            Console.WriteLine("  --help, -h                 Show this help");
        }
    }
}
