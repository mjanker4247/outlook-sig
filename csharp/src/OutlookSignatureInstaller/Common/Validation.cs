using PhoneNumbers;
using System;
using System.Linq;
using System.Net.Mail;

namespace OutlookSignatureInstaller.Common
{
    /// <summary>
    /// Mirrors pkg/common/validation.go: validates user inputs for name, email, phone,
    /// title, signature name, and URL.
    /// </summary>
    public static class Validation
    {
        /// <summary>Mirrors MinNameLength constant.</summary>
        public const int MinNameLength = 2;

        /// <summary>Mirrors BufferSizeLimit: 5 MB max for reading template files.</summary>
        public const long BufferSizeLimit = 5L * 1024 * 1024;

        /// <summary>Mirrors FileSizeLimit: 50 MB max for web template downloads.</summary>
        public const long FileSizeLimit = 50L * 1024 * 1024;

        private static readonly char[] AllowedNamePunctuation = { ' ', '.', '-', '\'' };
        private static readonly char[] NamePunctuation = { '.', '-', '\'' };

        /// <summary>
        /// Mirrors ValidateName(): validates a person's full name.
        /// Rules:
        ///   - At least MinNameLength characters
        ///   - Only Unicode letters, spaces, dots, hyphens, apostrophes
        ///   - No consecutive punctuation characters
        ///   - Cannot start or end with a hyphen or apostrophe (dots are OK)
        /// </summary>
        public static ValidationError? ValidateName(string name)
        {
            if (string.IsNullOrEmpty(name) || name.Length < MinNameLength)
                return new ValidationError("name", $"Name must be at least {MinNameLength} characters long.");

            foreach (char c in name)
            {
                if (!char.IsLetter(c) && !AllowedNamePunctuation.Contains(c))
                    return new ValidationError("name",
                        "Name may only contain letters, spaces, dots, hyphens, and apostrophes.");
            }

            // No consecutive punctuation (e.g. ".." or "-." or "'-")
            for (int i = 0; i < name.Length - 1; i++)
            {
                if (NamePunctuation.Contains(name[i]) && NamePunctuation.Contains(name[i + 1]))
                    return new ValidationError("name",
                        "Name must not contain consecutive punctuation characters.");
            }

            char first = name[0];
            char last = name[^1];
            if (first == '-' || first == '\'')
                return new ValidationError("name", "Name must not start with a hyphen or apostrophe.");
            if (last == '-' || last == '\'')
                return new ValidationError("name", "Name must not end with a hyphen or apostrophe.");

            return null;
        }

        /// <summary>
        /// Mirrors ValidateEmail(): parses with MailAddress (equivalent to Go's net/mail.ParseAddress).
        /// </summary>
        public static ValidationError? ValidateEmail(string email)
        {
            if (string.IsNullOrWhiteSpace(email))
                return new ValidationError("email", "Email address is required.");

            try
            {
                _ = new MailAddress(email);
                return null;
            }
            catch
            {
                return new ValidationError("email", "Invalid email address.");
            }
        }

        /// <summary>
        /// Mirrors ValidatePhoneNumber(): uses libphonenumber-csharp (Google's libphonenumber port).
        /// Defaults to "DE" country code as fallback, same as the Go implementation.
        /// </summary>
        public static ValidationError? ValidatePhoneNumber(string phone)
        {
            if (string.IsNullOrWhiteSpace(phone))
                return new ValidationError("phone", "Phone number is required.");

            try
            {
                var util = PhoneNumberUtil.GetInstance();
                var parsed = util.Parse(phone, "DE");
                if (!util.IsValidNumber(parsed))
                    return new ValidationError("phone", "Invalid phone number.");
                return null;
            }
            catch
            {
                return new ValidationError("phone", "Invalid phone number.");
            }
        }

        /// <summary>
        /// Mirrors FormatPhoneNumber(): returns (display, link) formats.
        /// - display: INTERNATIONAL format, e.g. "+49 211 123456"
        /// - link:    E164 format for tel: links, e.g. "+49211123456"
        /// </summary>
        public static (string Display, string Link) FormatPhoneNumber(string phone, string countryCode = "DE")
        {
            var util = PhoneNumberUtil.GetInstance();
            var parsed = util.Parse(phone, countryCode);
            string display = util.Format(parsed, PhoneNumberFormat.INTERNATIONAL);
            string link = util.Format(parsed, PhoneNumberFormat.E164);
            return (display, link);
        }

        /// <summary>
        /// Mirrors ValidateTitle(): intentionally permissive — no validation performed.
        /// </summary>
        public static ValidationError? ValidateTitle(string title) => null;

        /// <summary>
        /// Mirrors ValidateSignatureName(): rejects names with filesystem-special characters.
        /// Forbidden: / \ : * ? " &lt; &gt; |
        /// </summary>
        public static ValidationError? ValidateSignatureName(string name)
        {
            if (string.IsNullOrWhiteSpace(name))
                return new ValidationError("signature_name", "Signature name is required.");

            char[] forbidden = { '/', '\\', ':', '*', '?', '"', '<', '>', '|' };
            if (name.IndexOfAny(forbidden) >= 0)
                return new ValidationError("signature_name",
                    $"Signature name must not contain any of: {string.Join(" ", forbidden)}");

            return null;
        }

        /// <summary>
        /// Mirrors ValidateURL(): allows only HTTP and HTTPS URLs. Empty string is allowed
        /// (treated as optional, matching Go behaviour where empty base_url is valid for local mode).
        /// </summary>
        public static ValidationError? ValidateURL(string rawUrl)
        {
            if (string.IsNullOrWhiteSpace(rawUrl))
                return null;

            if (!Uri.TryCreate(rawUrl, UriKind.Absolute, out var uri) ||
                (uri.Scheme != "http" && uri.Scheme != "https"))
            {
                return new ValidationError("url", "URL must be a valid HTTP or HTTPS address.");
            }

            return null;
        }
    }

    /// <summary>
    /// Mirrors the ValidationError type in Go: carries the field name and a human-readable message.
    /// C# change: implemented as a class rather than a struct since it is nullable (returned as nil/null).
    /// </summary>
    public class ValidationError
    {
        public string Field { get; }
        public string Message { get; }

        public ValidationError(string field, string message)
        {
            Field = field;
            Message = message;
        }

        public override string ToString() => $"{Field}: {Message}";
    }
}
