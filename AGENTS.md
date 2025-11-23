# Outlook-Signature-Installer

A tool written in golang to install and manage email signatures in Microsoft Outlook for macOS and Windows with a command-line and graphical user interface.

## Development Setup


## Tech Stack

- **Framework**: fyne
- **Language**: golang
- **Styling**: [z.B. Tailwind CSS v4]
- **Testing**: go test
- `gofmt`, `goimports` – for formatting and imports.
- `golangci-lint` – includes `govet`, `golint`, `staticcheck`, etc.
- `go vet` – built-in static analyzer.
- `staticcheck` – powerful linting tool.
- `task` - build tool
- Replace `panic` with error returns unless truly unrecoverable.

## Project Structure

```bash
src/
├─ go.mod
├─ Taskfile.yml
├─ README.md
├─ build/            # The build directory with the binaries
├─ cmd/              # main package
├─ pkg/              # Utilities functions and services
└─ templates/        # HTML and text templates
```

## Code Standards

### General Rules

- Write tests for new features

## Constants Over Magic Numbers

- Replace hard-coded values with named constants
- Use descriptive constant names that explain the value's purpose
- Keep constants at the top of the file or in a dedicated constants file

## Meaningful Names

- Variables, functions, and classes should reveal their purpose
- Names should explain why something exists and how it's used
- Avoid abbreviations unless they're universally understood

## Smart Comments

- Don't comment on what the code does - make the code self-documenting
- Use comments to explain why something is done a certain way
- Document APIs, complex algorithms, and non-obvious side effects

## Single Responsibility

- Each function should do exactly one thing
- Functions should be small and focused
- If a function needs a comment to explain what it does, it should be split

## DRY (Don't Repeat Yourself)

- Extract repeated code into reusable functions
- Share common logic through proper abstraction
- Maintain single sources of truth

## Clean Structure

- Keep related code together
- Organize code in a logical hierarchy
- Use consistent file and folder naming conventions

## Encapsulation

- Hide implementation details
- Expose clear interfaces
- Move nested conditionals into well-named functions

## Code Quality Maintenance

- Refactor continuously
- Fix technical debt early
- Leave code cleaner than you found it

## Testing

- Write tests before fixing bugs
- Keep tests readable and maintainable
- Test edge cases and error conditions

## Version Control

- Write clear commit messages
- Make small, focused commits
- Use meaningful branch names

### Naming Conventions

- Functions: camelCase (formatDate.go)
- Constants: SCREAMING_SNAKE_CASE
- Types/Interfaces: PascalCase with suffix (UserType)

### File Organization

- Colocate tests with source files
- Group related components in folders

## Testing Guidelines

- Write tests alongside implementation
- Focus on user behavior, not implementation
- Maintain &gt; 80% coverage for critical paths

## Common Pitfalls to Avoid

- DON'T: Create new files unless necessary
- DON'T: Use console.log in production code
- DON'T: Ignore errors
- DON'T: Skip tests for "simple" features
- DO: Check existing components before creating new ones
- DO: Follow established patterns in the codebase
- DO: Keep functions small and focused

## Deployment

- Main branch deploys to production
- PR previews on Vercel
