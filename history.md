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
