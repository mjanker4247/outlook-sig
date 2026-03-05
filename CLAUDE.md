# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

Uses [Task](https://taskfile.dev) as the build tool (`brew install go-task` on macOS).

```bash
task                        # Build for current platform + copy config/templates
task build                  # Compile binary only
task cross-build-windows    # Cross-compile for Windows from macOS/Linux
task test                   # Run all tests: go test ./...
task lint                   # Run golangci-lint
task fmt                    # Format with go fmt + goimports
task check                  # fmt + lint + test
task clean                  # Remove build/ directory
task vendor                 # Sync vendor/ with go.mod
```

Run a single test:
```bash
go test ./pkg/common/... -run TestValidateName
go test ./pkg/signature/... -run TestInstall
```

## Architecture

Single binary with dual-mode operation (CLI and GUI via [Fyne](https://fyne.io)).

**Entry point:** `cmd/signature-installer/main.go` → `pkg/cli` → either `pkg/gui` or inline CLI flow.

**Package responsibilities:**
- `pkg/cli` — CLI flag parsing (urfave/cli), prompts for missing fields, calls `pkg/signature.Installer`
- `pkg/gui` — Fyne GUI form; validates input inline and calls `pkg/signature.Installer`
- `pkg/signature` — Core logic: loads `config.yaml`, optionally downloads templates from a web server, renders HTML (Go `html/template`) and text templates, writes to the platform-specific Outlook signatures directory
- `pkg/common` — Shared utilities: `GetTemplateBase()` (resolves `<exe_dir>/templates/`), input validation (name, email, phone via nyaruka/phonenumbers), phone formatting, shared constants (`BufferSizeLimit`, `FileSizeLimit`)

**Configuration flow:**
1. `config.yaml` lives next to the binary in the build directory (copied by `task copy-config`)
2. `Installer.LoadConfig()` reads it; `template_source` is either `"local"` or `"web"`
3. When `"web"`, `DownloadWebTemplates()` fetches `<base_url><template_name>.htm` and `<base_url><template_name>.txt`
4. `Installer.Install(data)` renders templates and writes to the Outlook signatures directory

**Template placeholders** (HTML uses `html/template`, text uses regex replacement):
- `{{ .Name }}`, `{{ .Title }}`, `{{ .Email }}`, `{{ .PhoneDisplay }}`, `{{ .PhoneLink }}`

**Outlook signature directories:**
- macOS: `~/Library/Group Containers/UBF8T346G9.Office/Outlook/Outlook 15 Profiles/Main Profile/Signatures/`
- Windows: `%APPDATA%\Microsoft\Signatures\`

**Filesystem abstraction:** `pkg/signature` uses `afero.Fs`, making it testable via `afero.NewMemMapFs()` in tests.

## Code Standards

- Go naming: functions camelCase, constants SCREAMING_SNAKE_CASE, types/interfaces PascalCase
- Replace `panic` with error returns unless truly unrecoverable
- Tests colocate with source files (`_test.go` in same package)
- Phone number parsing defaults to `"DE"` country code as fallback
