# OpenSPDD

AI Coding Assistant Command Template Manager - A CLI tool for managing command templates across Cursor, Claude Code, Antigravity, and GitHub Copilot environments.

## Features

- **Auto-detection**: Automatically detects your AI coding environment (Cursor, Claude Code, Antigravity, GitHub Copilot)
- **Template Management**: Embedded templates distributed via a single binary
- **Interactive UI**: Modern terminal UI for template selection
- **Cross-platform**: Supports macOS, Linux, and Windows

## Installation

### Homebrew (macOS/Linux)

```bash
brew install gszhangwei/tools/openspdd
```

Or:

```bash
brew tap gszhangwei/tools
brew install openspdd
```

Upgrade to latest version:

```bash
brew upgrade openspdd
```

### Go Install

```bash
go install github.com/gszhangwei/open-spdd@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/gszhangwei/open-spdd/releases).

## Usage

First, navigate to your project directory:

```bash
cd /path/to/your/project
```

### Initialize Environment

```bash
# Auto-detect and initialize
openspdd init

# Specify tool manually
openspdd --tool cursor init
```

### List Templates

```bash
# List all templates
openspdd list

# Filter by category
openspdd list -c Development

# Quiet mode (for scripting)
openspdd list -q
```

### Generate Templates

```bash
# Generate all default templates
openspdd generate --all

# Interactive selection
openspdd generate

# Generate specific template
openspdd generate spdd-generate

# Force overwrite existing files
openspdd generate --force spdd-generate

# Custom output directory
openspdd generate --output ./my-commands spdd-generate
```

### Global Flags

```bash
# Manually specify AI tool
openspdd --tool cursor <command>
openspdd --tool claude-code <command>
openspdd --tool antigravity <command>
openspdd --tool github-copilot <command>
```

## Supported Environments

| Tool           | Detection                                              | Config Directory           |
| -------------- | ------------------------------------------------------ | -------------------------- |
| Cursor         | `.cursor/`, `.cursorrules`                             | `.cursor/commands/`        |
| Claude Code    | `.claude/`, `CLAUDE.md`                                | `.claude/commands/`        |
| Antigravity    | `.antigravity/`                                        | `.antigravity/commands/`   |
| GitHub Copilot | `.github/copilot-instructions.md`, `.github/copilot-prompts/` | `.github/copilot-prompts/` |

### GitHub Copilot File Structure

For GitHub Copilot, OpenSPDD generates a different file structure:

```
.github/
├── copilot-instructions.md     # Main instruction file (auto-merged with markers)
└── copilot-prompts/
    ├── spdd-reasons-canvas.md  # REASONS-Canvas workflow
    ├── spdd-generate.md        # Code generation workflow
    └── spdd-sync.md            # Sync workflow
```

The `copilot-instructions.md` file uses marker-based merging (`<!-- openspdd:start -->` and `<!-- openspdd:end -->`) to preserve any custom content you add outside the marked section.

## Available Templates

| Template              | Description                                     |
| --------------------- | ----------------------------------------------- |
| `spdd-generate`       | Generate code from structured SPDD prompt files |
| `spdd-sync`           | Sync code changes back to SPDD prompt files     |
| `spdd-reasons-canvas` | Generate REASONS-Canvas structured prompts      |

## Building from Source

```bash
# Clone the repository
git clone https://github.com/gszhangwei/open-spdd.git
cd open-spdd

# Build
go build -o openspdd .

# Install to GOPATH/bin
go install .
```

## Testing

Tests are organized in the `tests/` directory, structured by module:

```
tests/
├── cmd/           # CLI command tests
├── detector/      # Environment detection tests
├── templates/     # Template management tests
├── ui/            # UI renderer tests
└── internal/      # Error constants tests
```

### Running Tests

```bash
# Run all tests
go test ./tests/...

# Run tests with verbose output
go test ./tests/... -v

# Run specific module tests
go test ./tests/detector/...
go test ./tests/templates/...
```

## License

[MIT License](LICENSE)
