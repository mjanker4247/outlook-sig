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
- Make

## Installation

1. Clone the repository:
```bash
git clone https://git.ululuu.de/jankerm/outlook-signature.git
cd outlook-signature
```

## Building

The project includes a Makefile for easy building on different platforms:

### Build for all platforms
```bash
make build-all
```

The built binaries will be available in the `build/` directory:
- macOS: `build/SignatureInstaller`
- Windows: `build/SignatureInstaller.exe`

## Usage

After building, you can run the program with:

```bash
SignatureInstaller[.exe] [options]
```

### Command Line Options

- `--name`: Your full name
- `--email`: Your email address
- `--phone`: Your phone number
- `--template`: Name of the template to use (default: "OutlookSignature")

Example:
```bash
SignatureInstaller[.exe] --name "John Doe" --email "john.doe@example.com" --phone "+49123456789"
```

## Templates

The program uses HTML and text templates located in the `templates/` directory. You can customize these templates to match your organization's branding.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Sign

osslsigncode sign \
  -pkcs12 code-sign-certificate.pfx \
  -askpass \
  -n "Program Name" \
  -i https://www.program-website.com \
  -in program.exe \
  -out program-signed.exe

