# Outlook Signature Installer

A tool to install and manage email signatures in Microsoft Outlook for macOS and Windows. It supports both command-line and graphical user interface modes.

## Features

- Install email signatures with custom templates
- Support for both HTML and plain text signatures
- Automatic phone number formatting
- Cross-platform support (macOS and Windows)
- Command-line interface (CLI) and Graphical User Interface (GUI)
- Easy template customization
- Images embedded in HTML, no need for external links
- Configuration-based template selection

## Prerequisites

- Go 1.24.2
- Microsoft Outlook installed
- Git (for cloning the repository)
- Task (for building)

## Installation

1. Install prerequisites using the provided scripts:

On Windows (PowerShell):
```powershell
.\install-prerequisites.ps1
```

On macOS:
```bash
chmod +x install-prerequisites.sh
./install-prerequisites.sh
```

2. Clone the repository:
```bash
git clone https://git.ululuu.de/jankerm/outlook-signature.git
cd outlook-signature
```

2. Install Task (build tool):
```bash
# macOS
brew install go-task

# Windows
scoop install task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin
```

## Building

The project uses Task for building. Here are the available commands:

```bash
# Build for current platform
task build

# Build for Windows (cross-compilation)
task cross-build-windows

# Build and sign Windows executable
task sign-windows
```

The built binaries will be available in the `build/` directory:
- macOS: `build/SignatureInstaller`
- Windows: `build/SignatureInstaller.exe`

## Configuration

The application uses a configuration file (`config.yaml`) located in the build root directory to specify which template to use and where to source templates from. The configuration file should contain:

```yaml
template_name: "Standard"
template_source: "local"  # "local" or "web"
web_templates:
  base_url: "https://team.emea.tuv.group/sites/002458/TRLP%20KM%20Dokumente/Outlook%20Signatur/"
  template_files:
    - "Standard.htm"
    - "Standard.txt"
```

The `template_name` field should match the base filename of your template (without the `.htm` or `.txt` extension).

### Template Sources

The application supports two template sources:

#### Local Templates (Default)
Set `template_source: "local"` to use templates stored locally in the `templates/` directory.

#### Web Templates
Set `template_source: "web"` to download templates from a web server. This requires additional configuration:

- `base_url`: The base URL where templates are hosted
- `template_files`: List of template files to download

When using web templates, the application will:
1. Download the specified template files from the web server
2. Store them locally in the `templates/` directory
3. Use the downloaded templates for signature generation

**Note**: Web templates are downloaded each time the application runs, ensuring you always have the latest versions.

## Usage

### GUI Mode

The application can be run in GUI mode in two ways:

1. Without any arguments:
```bash
SignatureInstaller
```

2. Using the --gui flag:
```bash
SignatureInstaller --gui
```

The GUI provides a form to enter:
- Your full name
- Email address
- Phone number

The template to use is automatically determined from the configuration file.

### Command Line Mode

```bash
SignatureInstaller[.exe] [options]
```

#### Command Line Options

- `--name`, `-n`: Your name (can include profession/title on separate lines)
- `--email`, `-e`: Your email address
- `--phone`, `-p`: Your phone number
- `--template-source`, `-s`: Template source: 'local' or 'web' (overrides config.yaml)
- `--gui`, `-g`: Launch in GUI mode

Example:
```bash
SignatureInstaller --name "John Doe" --email "john.doe@example.com" --phone "+49123456789"
```

Using multiline name (with profession/title):
```bash
SignatureInstaller --name "John Doe\nSoftware Engineer" --email "john.doe@example.com" --phone "+49123456789"
```

Using web templates:
```bash
SignatureInstaller --name "John Doe" --email "john.doe@example.com" --phone "+49123456789" --template-source web
```

## Templates

The program uses HTML and text templates located in the `templates/` directory. You can customize these templates to match your organization's branding.

### Template Placeholders

Templates may include these placeholders:
- `{{ .Name }}`: Your full name (supports multiline: name, profession, title)
- `{{ .Email }}`: Your email address
- `{{ .PhoneLink }}`: Phone number in E.164 format (for links)
- `{{ .PhoneDisplay }}`: Formatted phone number (for display)

### Template Files

- `Standard.htm`: HTML signature template
- `Standard.txt`: Plain text signature template
- `Standard_files/`: Directory for images and other assets

## Signature Locations

The signatures are installed in the following locations:

### Windows
```
%APPDATA%\Microsoft\Signatures\
```

### macOS
```
~/Library/Group Containers/UBF8T346G9.Office/Outlook/Outlook 15 Profiles/Main Profile/Signatures/
```

## Development

### VS Code Integration

The project includes VS Code configuration for easy development:

1. Open the project in VS Code
2. Use the Command Palette (Cmd/Ctrl + Shift + P)
3. Type "Tasks: Run Task" to see available build tasks:
   - Build
   - Build Windows
   - Sign Windows

### Debugging

1. Open the project in VS Code
2. Press F5 or use the Run menu
3. The debugger will launch with the default configuration

## License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.

## Summary

I've successfully removed the template scanner functionality and implemented a configuration-based approach. Here are the key changes made:

### 1. **Removed Template Scanning**
- Removed `GetAvailableTemplates()` function from `pkg/common/paths.go`
- Removed template selection dropdown from the GUI
- Removed `--template` flag from the CLI

### 2. **Added Configuration System**
- Created a `Config` struct in `pkg/signature/signature.go` with `TemplateName` field
- Added `LoadConfig()` method to load configuration from `config.yaml` in the build root
- Modified the `Install()` method to use the configured template instead of a parameter

### 3. **Updated User Interfaces**
- **CLI**: Removed template parameter, now automatically uses configured template
- **GUI**: Removed template selection, simplified form with only name, email, and phone fields

### 4. **Configuration File**
- Created `config.yaml` that specifies which template to use
- The file is placed in the build root directory alongside the executable

### 5. **Build Process Updates**
- Updated `Taskfile.yml` to copy the configuration file to the build directory
- Added `copy-config` task for both platforms

### 6. **Documentation Updates**
- Updated `README.md` to reflect the new configuration-based approach
- Removed references to template selection and scanning

The application now uses a single, preconfigured template specified in the configuration file, eliminating the need for template scanning while maintaining the same functionality for users. The template name is automatically loaded from the configuration when installing signatures.