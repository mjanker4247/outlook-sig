# Change History

## 2024-12-19 - Migrated from JSON to YAML Configuration

**Summary**: Replaced `config.json` with `config.yaml` throughout the codebase for better configuration management.

**Changes**:
- Updated configuration file format from JSON to YAML
- Modified Go code to use YAML parsing instead of JSON
- Updated build system (Taskfile.yml) to handle YAML config
- Updated documentation and help text references
- Fixed template copying in build process

**Files affected**: `pkg/signature/signature.go`, `pkg/cli/cli.go`, `Taskfile.yml`, `README.md`, `config.yaml`, `config.json` (removed)

## 2024-12-19 - Added Web Template Support

**Summary**: Added functionality to download templates from web servers, with configurable template sources and CLI override options.

**Changes**:
- Added web template configuration support in config.yaml
- Implemented HTTP template downloading functionality
- Added `--template-source` CLI flag to override configuration
- Enhanced configuration structure with template source options
- Updated documentation for web template usage

**Files affected**: `pkg/signature/signature.go`, `pkg/cli/cli.go`, `config.yaml`, `README.md`

## 2024-12-19 - Created Test Server for Web Templates

**Summary**: Added a simple local HTTP server for testing web template download functionality.

**Changes**:
- Created standalone test server in `test-server/` directory
- Added test configuration file for local testing
- Created automated test script for web template functionality
- Added test server to .gitignore (excluded from builds)

**Files affected**: `test-server/main.go`, `test-server/config-test.yaml`, `test-server/README.md`, `test-server/test-web-templates.sh`, `.gitignore`

## 2024-12-19 - Simplified Multiline Name Implementation

**Summary**: Streamlined the multiline name functionality by removing the separate multiline flag and consolidating into a single name field.

**Changes**:
- Removed `--multiline-name` flag, consolidated into single `--name` flag
- Enhanced validation to automatically filter out empty/non-visible lines
- Updated template processing to clean up names by removing empty lines
- Simplified CLI interface while maintaining full functionality

**Files affected**: `pkg/cli/cli.go`, `pkg/common/validation.go`, `pkg/signature/signature.go`, `README.md`
