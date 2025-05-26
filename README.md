# Outlook Signature Installer

A tool to install and manage email signatures in Microsoft Outlook for macOS and Windows. It supports both command-line and graphical user interface modes.

## Features

- Install email signatures with custom templates
- Support for both HTML and plain text signatures
- Automatic phone number formatting
- Cross-platform support (macOS and Windows)
- Command-line interface (CLI) and Graphical User Interface (GUI)
- Easy template customization
- Automatic image handling

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

# Clean build directory
task clean

# Run the application
task run

# Build and sign Windows executable
task sign-windows
```

The built binaries will be available in the `build/` directory:
- macOS: `build/SignatureInstaller`
- Windows: `build/SignatureInstaller.exe`

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
- Template name (defaults to "OutlookSignature")

### Command Line Mode

```bash
SignatureInstaller[.exe] [options]
```

#### Command Line Options

- `--name`, `-n`: Your full name
- `--email`, `-e`: Your email address
- `--phone`, `-p`: Your phone number
- `--template`, `-t`: Name of the template to use (default: "OutlookSignature")
- `--gui`, `-g`: Launch in GUI mode

Example:
```bash
SignatureInstaller --name "John Doe" --email "john.doe@example.com" --phone "+49123456789"
```

## Templates

The program uses HTML and text templates located in the `templates/` directory. You can customize these templates to match your organization's branding.

### Template Placeholders

Templates may include these placeholders:
- `{{ .Name }}`: Your full name
- `{{ .Email }}`: Your email address
- `{{ .PhoneLink }}`: Phone number in E.164 format (for links)
- `{{ .PhoneDisplay }}`: Formatted phone number (for display)

### Template Files

- `OutlookSignature.htm`: HTML signature template
- `OutlookSignature.txt`: Plain text signature template
- `OutlookSignature_files/`: Directory for images and other assets

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
   - Clean
   - Run
   - Build Windows
   - Sign Windows
   - Build All

### Debugging

1. Open the project in VS Code
2. Press F5 or use the Run menu
3. The debugger will launch with the default configuration

## License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.