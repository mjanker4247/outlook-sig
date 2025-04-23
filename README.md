# Outlook Signature Installer

A command-line tool to install and manage email signatures in Microsoft Outlook for macOS and Windows.

## Features

- Install email signatures with custom templates
- Support for both HTML and plain text signatures
- Automatic phone number formatting
- Cross-platform support (macOS and Windows)
- Easy-to-use command-line interface

## Prerequisites

- Go 1.16 or later
- Microsoft Outlook installed
- Git (for cloning the repository)

## Installation

1. Clone the repository:
```bash
git clone https://git.ululuu.de/jankerm/outlook-signature.git
cd outlook-signature
```

2. Install dependencies:
```bash
make deps
```

## Building

The project includes a Makefile for easy building on different platforms:

### Build for macOS
```bash
make build-mac
```

### Build for Windows
```bash
make build-win
```

### Build for all platforms
```bash
make build-all
```

The built binaries will be available in the `build/` directory:
- macOS: `build/mac/signature-installer`
- Windows: `build/win/signature-installer.exe`

## Usage

After building, you can run the program with:

```bash
./build/mac/signature-installer [options]
```

or on Windows:
```bash
.\build\win\signature-installer.exe [options]
```

### Command Line Options

- `--name`: Your full name
- `--email`: Your email address
- `--phone`: Your phone number
- `--template`: Name of the template to use (default: "OutlookSignature")

Example:
```bash
./signature-installer --name "John Doe" --email "john.doe@example.com" --phone "+49123456789"
```

## Templates

The program uses HTML and text templates located in the `templates/` directory. You can customize these templates to match your organization's branding.

## License

This project is licensed under the MIT License - see the LICENSE file for details.