# Assets Directory

This directory contains assets required for building the application.

## Required Files

### Icon.png
- Required for macOS application bundle
- Recommended size: 1024x1024 pixels
- Format: PNG with alpha channel
- Place your application icon here named exactly as `Icon.png`

## Usage

The build system will automatically use these assets when building the application:

- macOS: Uses `Icon.png` for the application bundle icon
- Windows: Can be configured to use `Icon.png` for the application icon (requires conversion to .ico) 