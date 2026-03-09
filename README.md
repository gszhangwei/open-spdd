# OpenSPDD

AI Coding Assistant Command Template Manager - A CLI tool for managing command templates across Cursor, Claude Code, and Antigravity environments.

## Features

- **Auto-detection**: Automatically detects your AI coding environment (Cursor, Claude Code, Antigravity)
- **Template Management**: Embedded templates distributed via a single binary
- **Interactive UI**: Modern terminal UI for template selection
- **Cross-platform**: Supports macOS, Linux, and Windows

## Installation

### Homebrew (macOS/Linux)

```bash
brew install wwdzhang/tap/openspdd
```

### Go Install

```bash
go install github.com/wwdzhang/open-spdd@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/wwdzhang/open-spdd/releases).

## Usage

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
# Interactive selection
openspdd generate

# Generate specific template
openspdd generate spdd-generate

# Generate all templates
openspdd generate --all

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
```

## Supported Environments

| Tool | Detection | Config Directory |
|------|-----------|------------------|
| Cursor | `.cursor/`, `.cursorrules` | `.cursor/commands/` |
| Claude Code | `.claude/`, `CLAUDE.md` | `.claude/commands/` |
| Antigravity | `.antigravity/` | `.antigravity/commands/` |

## Available Templates

| Template | Description |
|----------|-------------|
| `spdd-generate` | Generate code from structured SPDD prompt files |
| `spdd-sync` | Sync code changes back to SPDD prompt files |
| `spdd-reasons-canvas` | Generate REASONS-Canvas structured prompts |

## Building from Source

```bash
# Clone the repository
git clone https://github.com/wwdzhang/open-spdd.git
cd open-spdd

# Build
go build -o openspdd .

# Install to GOPATH/bin
go install .
```

## License

[MIT License](LICENSE)
