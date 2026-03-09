using Xunit;
using OutlookSignatureInstaller.Common;

namespace OutlookSignatureInstaller.Tests.Common
{
    /// <summary>
    /// Mirrors pkg/common/validation_test.go: unit tests for all validation functions.
    /// xunit [Theory] + [InlineData] replaces Go's table-driven tests.
    /// </summary>
    public class ValidationTests
    {
        // ─── ValidateName ────────────────────────────────────────────────────────────

        [Theory]
        [InlineData("John Doe")]
        [InlineData("O'Brien")]
        [InlineData("Anne-Marie")]
        [InlineData("Dr. Smith")]
        [InlineData("Müller")]       // Unicode letters
        [InlineData("Jean-Luc")]
        public void ValidateName_ValidNames_ReturnsNull(string name) =>
            Assert.Null(Validation.ValidateName(name));

        [Theory]
        [InlineData("J")]           // too short
        [InlineData("")]            // empty
        [InlineData("-John")]       // starts with hyphen
        [InlineData("John-")]       // ends with hyphen
        [InlineData("'John")]       // starts with apostrophe
        [InlineData("John'")]       // ends with apostrophe
        [InlineData("John..Doe")]   // consecutive punctuation
        [InlineData("John--Doe")]   // consecutive hyphens
        [InlineData("John123")]     // digits not allowed
        [InlineData("John@Doe")]    // special chars not allowed
        public void ValidateName_InvalidNames_ReturnsError(string name) =>
            Assert.NotNull(Validation.ValidateName(name));

        // ─── ValidateEmail ───────────────────────────────────────────────────────────

        [Theory]
        [InlineData("user@example.com")]
        [InlineData("user.name+tag@example.co.uk")]
        [InlineData("first.last@subdomain.example.org")]
        public void ValidateEmail_ValidEmails_ReturnsNull(string email) =>
            Assert.Null(Validation.ValidateEmail(email));

        [Theory]
        [InlineData("not-an-email")]
        [InlineData("@missing-local.com")]
        [InlineData("missing-at-sign")]
        [InlineData("")]
        public void ValidateEmail_InvalidEmails_ReturnsError(string email) =>
            Assert.NotNull(Validation.ValidateEmail(email));

        // ─── ValidatePhoneNumber ─────────────────────────────────────────────────────

        [Theory]
        [InlineData("+49 211 123456")]       // German international
        [InlineData("+1 800 555 0100")]      // US toll-free
        [InlineData("+44 20 7946 0958")]     // UK London
        public void ValidatePhoneNumber_ValidNumbers_ReturnsNull(string phone) =>
            Assert.Null(Validation.ValidatePhoneNumber(phone));

        [Theory]
        [InlineData("123")]
        [InlineData("not-a-phone")]
        [InlineData("")]
        public void ValidatePhoneNumber_InvalidNumbers_ReturnsError(string phone) =>
            Assert.NotNull(Validation.ValidatePhoneNumber(phone));

        // ─── FormatPhoneNumber ───────────────────────────────────────────────────────

        [Fact]
        public void FormatPhoneNumber_GermanNumber_ReturnsInternationalAndE164()
        {
            var (display, link) = Validation.FormatPhoneNumber("+49 211 123456");
            Assert.Contains("+49", display);   // INTERNATIONAL format
            Assert.StartsWith("+49", link);    // E164 format
            Assert.DoesNotContain(" ", link);  // E164 has no spaces
        }

        // ─── ValidateTitle ───────────────────────────────────────────────────────────

        [Theory]
        [InlineData("")]
        [InlineData("Senior Engineer, Dept. A")]
        [InlineData("任意のタイトル")]
        public void ValidateTitle_AnyValue_ReturnsNull(string title) =>
            Assert.Null(Validation.ValidateTitle(title));

        // ─── ValidateSignatureName ───────────────────────────────────────────────────

        [Theory]
        [InlineData("Standard")]
        [InlineData("CyberSecurityDays")]
        [InlineData("My Signature")]
        public void ValidateSignatureName_ValidNames_ReturnsNull(string name) =>
            Assert.Null(Validation.ValidateSignatureName(name));

        [Theory]
        [InlineData("../evil")]
        [InlineData("bad\\name")]
        [InlineData("bad/name")]
        [InlineData("bad:name")]
        [InlineData("")]
        public void ValidateSignatureName_InvalidNames_ReturnsError(string name) =>
            Assert.NotNull(Validation.ValidateSignatureName(name));

        // ─── ValidateURL ─────────────────────────────────────────────────────────────

        [Theory]
        [InlineData("http://example.com/templates/")]
        [InlineData("https://secure.example.com/sigs/")]
        [InlineData("")]       // empty is allowed (optional field)
        public void ValidateURL_ValidURLs_ReturnsNull(string url) =>
            Assert.Null(Validation.ValidateURL(url));

        [Theory]
        [InlineData("ftp://invalid.com")]
        [InlineData("not-a-url")]
        [InlineData("file:///local/path")]
        public void ValidateURL_InvalidURLs_ReturnsError(string url) =>
            Assert.NotNull(Validation.ValidateURL(url));
    }
}
