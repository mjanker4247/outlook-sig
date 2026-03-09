# Outlook Signature Installer

A tool to install and manage email signatures in Microsoft Outlook for macOS, Windows, and Linux. Supports both command-line and graphical user interface modes.

## Features

- Install HTML and plain text Outlook signatures from templates
- GUI and CLI modes in a single binary
- Template placeholders: name (multiline with title), email, phone (display + E.164 link)
- Local or web-hosted template sources
- Cross-platform: macOS, Windows, Linux

## Prerequisites

### Required (all platforms)

| Tool | Purpose |
|---|---|
| Go 1.24+ | Build toolchain |
| Git | Source control |
| Task | Build automation (`go-task`) |
| golangci-lint | Linting |
| goimports | Code formatting |

### Optional

| Tool | Purpose | Platform |
|---|---|---|
| Docker + fyne-cross | Cross-compilation | Linux, macOS |
| osslsigncode | Cross-sign Windows `.exe` from Linux/macOS | Linux, macOS |
| signtool | Sign Windows `.exe` natively | Windows (via Windows SDK) |
| Xcode Command Line Tools | Sign and notarize macOS binary | macOS |
| gnupg | GPG-sign Linux binary | Linux |

## Installation

### Windows (PowerShell + Scoop)

Run the provided script — it installs Scoop if needed and sets up all required tools:

```powershell
.\install-prerequisites.ps1
```

Or install manually:

```powershell
# Install Scoop
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
Invoke-RestMethod get.scoop.sh | Invoke-Expression

scoop bucket add main extras versions
scoop install git go@1.24.2 task golangci-lint
scoop install extras/docker-desktop        # for cross-compilation
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/fyne-io/fyne-cross@latest
```

> **Signing:** `signtool.exe` is included with the
> [Windows SDK](https://developer.microsoft.com/windows/downloads/windows-sdk/),
> which ships with Visual Studio Build Tools.

### macOS (Homebrew)

Run the provided script — it installs Homebrew if needed and sets up all required tools:

```bash
chmod +x install-prerequisites.sh
./install-prerequisites.sh
```

Or install manually:

```bash
brew install git go@1.24.2 go-task golangci-lint osslsigncode
brew install --cask docker             # for cross-compilation
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/fyne-io/fyne-cross@latest
```

> **Signing:** `codesign` and `xcrun notarytool` are bundled with Xcode Command Line Tools
> (`xcode-select --install`). Distribution signing requires a paid
> [Apple Developer account](https://developer.apple.com/programs/).

### Linux (apt)

```bash
sudo apt update
sudo apt install git docker.io gnupg osslsigncode

# Install Go 1.24+ (if not available via apt, download from https://go.dev/dl)
sudo apt install golang-go

# Add your user to the docker group (re-login after)
sudo usermod -aG docker $USER

# Install Task
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

# Install Go-based tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/fyne-io/fyne-cross@latest
```

> **Note:** Go 1.24 may not be available via apt on older distros.
> Install from [go.dev/dl](https://go.dev/dl) if `apt install golang-go` gives an older version.

## Building

```bash
# Build for current platform (binary + assets → build/)
task

# Build binary only
task build
```

### Cross-compilation (requires Docker + fyne-cross)

The app uses [Fyne](https://fyne.io) (CGO + OpenGL/GLFW), so cross-compilation requires
`fyne-cross`, which provides Docker images with the correct toolchains (MingW for Windows,
osxcross for macOS).

```bash
task fyne-cross-install    # install fyne-cross (once)

task cross-windows         # → build/SignatureInstaller.exe  (Linux/macOS)
task cross-linux           # → build/SignatureInstaller
task cross-darwin          # → build/SignatureInstaller       (Linux only)
task cross-all             # all platforms at once
```

## Code Signing

### Windows

```bash
# On Windows — uses signtool with certificate from the Windows certificate store:
task sign-windows

# On Linux/macOS — uses osslsigncode with a PFX/PKCS12 file:
export SIGNING_CERT_PFX=/path/to/cert.pfx
export SIGNING_CERT_PASSWORD=yourpassword
task sign-windows          # → build/SignatureInstaller-signed.exe
```

### macOS

```bash
export MACOS_SIGNING_IDENTITY="Developer ID Application: Your Name (TEAMID)"
task sign-macos

# Notarize for distribution outside the App Store:
export APPLE_ID=you@example.com
export APPLE_TEAM_ID=XXXXXXXXXX
export APPLE_APP_PASSWORD=xxxx-xxxx-xxxx-xxxx   # app-specific password from appleid.apple.com
task notarize-macos
```

### Linux (GPG)

```bash
export GPG_KEY_ID=your@email.com
task sign-linux            # → build/SignatureInstaller.asc
```

## Configuration

`config.yaml` must reside next to the binary. It is copied automatically by the build tasks.

```yaml
template_name: "Standard"
template_source: "local"   # "local" or "web"
base_url: ""               # required when template_source is "web"
                           # e.g. "http://intranet-server/templates/"
```

When `template_source: "web"`, the app downloads `<base_url><template_name>.htm` and
`<base_url><template_name>.txt` at startup.

### Per-user config

Each user's settings are persisted separately and survive application updates:

- macOS: `~/Library/Application Support/OutlookSignatureInstaller/config.yaml`
- Windows: `%APPDATA%\OutlookSignatureInstaller\config.yaml`

## Usage

### GUI mode

Double-click the binary, or run:

```bash
./SignatureInstaller --gui
```

### CLI mode

```bash
SignatureInstaller --name "Jane Doe" --email "jane@example.com" --phone "+49123456789"

# With title on a second line
SignatureInstaller --name "Jane Doe\nSoftware Engineer" --email "jane@example.com" --phone "+49123456789"

# Override template source at runtime
SignatureInstaller --name "Jane Doe" --email "jane@example.com" --phone "+49123456789" --template-source web
```

**Options:**

| Flag | Short | Description |
|---|---|---|
| `--name` | `-n` | Full name; use `\n` to put title on a second line |
| `--email` | `-e` | Email address |
| `--phone` | `-p` | Phone number (parsed with DE as default country code) |
| `--template-source` | `-s` | Override template source: `local` or `web` |
| `--gui` | `-g` | Launch in GUI mode |

## Templates

HTML templates use Go's `html/template` syntax; text templates use regex substitution.
Place template files in the `templates/` directory next to the binary.

**Available placeholders:**

| Placeholder | Description |
|---|---|
| `{{ .Name }}` | Full name (supports multiline) |
| `{{ .Title }}` | Job title (second line of name, if provided) |
| `{{ .Email }}` | Email address |
| `{{ .PhoneDisplay }}` | Phone number formatted for display |
| `{{ .PhoneLink }}` | Phone in E.164 format (for `tel:` links) |

Default template files: `Standard.htm`, `Standard.txt`, `Standard_files/` (images).

## Signature locations

| Platform | Path |
|---|---|
| Windows | `%APPDATA%\Microsoft\Signatures\` |
| macOS | `~/Library/Group Containers/UBF8T346G9.Office/Outlook/Outlook 15 Profiles/Main Profile/Signatures/` |

## Development

```bash
task test            # run all tests
task test-verbose    # run with -v -race
task lint            # golangci-lint
task fmt             # go fmt + goimports
task check           # fmt + lint + test
```

Run a single test package:

```bash
go test ./pkg/common/... -run TestValidateName
go test ./pkg/signature/... -run TestInstall
```

## License

GNU General Public License v3.0 — see [LICENSE](LICENSE) for details.
