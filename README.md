# OpenSPDD

> **Structured Prompt-Driven Development** — Transform AI coding prompts into executable design contracts

[Chinese](README.zh-CN.md) | [Design Philosophy](docs/design-philosophy.md)

OpenSPDD is a methodology and cross-platform CLI tool for the AI coding era. It upgrades AI coding prompts from "disposable inputs" to "executable design contracts" with bidirectional synchronization between design and implementation.

## Why OpenSPDD?

Existing AI coding tools generate plan documents, but these documents have fundamental limitations:

| Problem          | Typical Plan Documents                | REASONS Canvas                                                        |
| ---------------- | ------------------------------------- | --------------------------------------------------------------------- |
| **Nature**       | Task list                             | Design contract                                                       |
| **Constraints**  | None — AI improvises freely           | Explicit — Norms define "how", Safeguards define "what not to do"     |
| **Detail Level** | High-level: _"Create BillingService"_ | Precise: _method signatures, parameters, error handling, DI patterns_ |
| **Traceability** | None — docs don't update with code    | Yes — `/spdd-sync` enables reverse sync                               |
| **Validation**   | Vague — _"done when complete"_        | Explicit — exact error messages, HTTP status codes in Safeguards      |
| **Dependencies** | Implicit — AI infers                  | Explicit — Operations define strict execution order                   |

**The core insight**: Plans are "suggestions", REASONS Canvas is a "contract".

## The REASONS Canvas Framework

REASONS Canvas is a 7-dimensional structured design framework:

```
┌─────────────────────────────────────────────────────────────────────┐
│                        REASONS Canvas                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  R - Requirements    The "why" — business goals and scope            │
│  E - Entities        Domain model (Mermaid class diagrams)           │
│  A - Approach        Solution strategy and trade-offs                │
│  S - Structure       Architecture, inheritance, dependencies         │
│  O - Operations      Precise implementation tasks in order           │
│  N - Norms           Coding standards and patterns                   │
│  S - Safeguards      Constraints and guardrails                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Why 7 dimensions?**

- **R+E+A** = Design decisions ("Why" and "What")
- **S+O** = Implementation path ("How")
- **N+S** = Quality guardrails

All three are essential: without N+S, AI improvises; without S+O, AI restructures arbitrarily; without R+E+A, AI lacks context.

## Core Workflow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SPDD Workflow                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Business Requirement                                                │
│      │                                                               │
│      ▼                                                               │
│  /spdd-analysis ──────→ Strategic analysis (concepts, risks)         │
│      │                                                               │
│      ▼                                                               │
│  /spdd-reasons-canvas ─→ REASONS Canvas design document              │
│      │                                                               │
│      ▼                                                               │
│  /spdd-generate ───────→ AI generates code per contract              │
│      │                                                               │
│      ▼                                                               │
│  Code Review / Refactoring                                           │
│      │                                                               │
│      ▼                                                               │
│  /spdd-sync ───────────→ Reverse sync changes back to design         │
│      │                                                               │
│      ▼                                                               │
│  Design stays in sync ─→ Next iteration builds on accurate design    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Key principle**: _"When reality diverges, fix the prompt first — then update the code."_

## Features

- **Cross-platform**: Supports Cursor, Claude Code, GitHub Copilot, and Antigravity
- **Auto-detection**: Automatically detects your AI coding environment
- **Single Binary**: All templates embedded via Go's embed directive
- **Bidirectional Sync**: Keep design documents and code in sync
- **Interactive UI**: Modern terminal UI for command selection

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

### Go Install

```bash
go install github.com/gszhangwei/open-spdd/cmd/openspdd@latest
```

The binary is installed to `$(go env GOPATH)/bin/openspdd` (typically `~/go/bin/openspdd`). Make sure that directory is on your `$PATH`:

```bash
# zsh
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc

# bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc
```

The first time `openspdd` runs, it will also detect this and print a one-time hint with the exact command for your shell.

### One-shot installer script

If you cloned the repo, the script `scripts/install.sh` runs `go install` and prints PATH instructions automatically:

```bash
./scripts/install.sh           # installs @latest
./scripts/install.sh v1.2.3    # installs a specific tag
```

### Download Binary

Download from [GitHub Releases](https://github.com/gszhangwei/open-spdd/releases).

### Uninstall

`openspdd uninstall` detects how the binary was installed (Homebrew or `go install`) and runs the matching cleanup. The plan is printed before any change, and confirmation is required by default.

```bash
# Preview without changing anything
openspdd uninstall --dry-run

# Interactive (default): prints the plan, asks for confirmation
openspdd uninstall

# Non-interactive (e.g., for scripts): skip the confirmation prompt
openspdd uninstall --yes
```

For a Homebrew install, this is equivalent to running `brew uninstall gszhangwei/tools/openspdd` plus a small first-run-marker cleanup. For a `go install` install, it removes the binary at the resolved path. If the install method cannot be classified (e.g., a manually-copied binary), `uninstall` refuses to act and prints the resolved path so you can remove it manually.

> **Scope**: Only the openspdd binary and openspdd's own first-run marker are removed. Generated SPDD command templates inside your projects (`.cursor/commands/spdd-*.md`, `.claude/commands/spdd-*.md`, etc.) are user files and are left untouched. The Homebrew tap `gszhangwei/tools` is also left in place.

## Quick Start

```bash
# Print the installed version
openspdd -v

# Navigate to your project
cd your-project

# Initialize (auto-detects AI tool)
openspdd init

# Generate SPDD commands
openspdd generate --all
```

Then in your AI coding tool, follow the complete SPDD workflow:

```bash
# Step 1: Strategic analysis (recommended for complex features)
/spdd-analysis @requirements/user-registration.md

# Step 2: Generate REASONS Canvas from analysis
/spdd-reasons-canvas @spdd/analysis/xxx.md

# Step 3: Generate code from REASONS Canvas
/spdd-generate @spdd/prompt/xxx.md

# Step 4: After code review/refactoring, sync changes back
/spdd-sync @spdd/prompt/xxx.md
```

For simpler features, you can skip Step 1 and provide requirements directly:

```bash
/spdd-reasons-canvas Implement user registration with email verification
```

## Usage

### Initialize Environment

```bash
# Auto-detect and initialize
openspdd init

# Specify tool manually
openspdd --tool cursor init
```

### List Commands

```bash
# List available commands (core + tool-specific)
openspdd list

# List optional commands
openspdd list --optional

# List all commands
openspdd list --all

# Filter by category
openspdd list -c Development
```

### Generate Commands

```bash
# Generate all default commands
openspdd generate --all

# Interactive selection
openspdd generate

# Generate specific command
openspdd generate spdd-generate

# Force overwrite
openspdd generate --force spdd-generate
```

### Global Flags

```bash
openspdd --tool cursor <command>
openspdd --tool claude-code <command>
openspdd --tool antigravity <command>
openspdd --tool github-copilot <command>
```

## Supported Environments

| Tool           | Detection                                                     | Config Directory           |
| -------------- | ------------------------------------------------------------- | -------------------------- |
| Cursor         | `.cursor/`, `.cursorrules`                                    | `.cursor/commands/`        |
| Claude Code    | `.claude/`, `CLAUDE.md`                                       | `.claude/commands/`        |
| Antigravity    | `.antigravity/`                                               | `.antigravity/commands/`   |
| GitHub Copilot | `.github/copilot-instructions.md`, `.github/copilot-prompts/` | `.github/copilot-prompts/` |

### GitHub Copilot File Structure

```
.github/
├── copilot-instructions.md     # Main instruction file (auto-merged with markers)
└── copilot-prompts/
    ├── spdd-analysis.md
    ├── spdd-reasons-canvas.md
    ├── spdd-generate.md
    ├── spdd-prompt-update.md
    └── spdd-sync.md
```

## Available Commands

### Core Commands

| Command               | Description                                       |
| --------------------- | ------------------------------------------------- |
| `spdd-analysis`       | Strategic analysis of requirements                |
| `spdd-reasons-canvas` | Generate REASONS-Canvas structured prompts        |
| `spdd-generate`       | Generate code from structured SPDD prompt files   |
| `spdd-prompt-update`  | Update existing SPDD prompt with new requirements |
| `spdd-sync`           | Sync code changes back to SPDD prompt files       |

### Tool-Specific Commands

| Tool           | Command                | Description                       |
| -------------- | ---------------------- | --------------------------------- |
| GitHub Copilot | `copilot-instructions` | Main instruction file for Copilot |

### Optional Commands (Beta)

The following commands are available as beta — not installed by default, but can be installed manually:

| Command            | Description                                                                           |
| ------------------ | ------------------------------------------------------------------------------------- |
| `spdd-story`       | Decompose feature requirements into INVEST-compliant stories with acceptance criteria |
| `spdd-code-review` | Review code against REASONS-Canvas, detecting intent drift and violations             |
| `spdd-api-test`    | Generate self-contained shell scripts with cURL commands for API testing              |

```bash
# List all optional commands
openspdd list --optional

# Install a specific optional command
openspdd generate spdd-story
openspdd generate spdd-code-review
openspdd generate spdd-api-test
```

## Plan vs REASONS Canvas: An Example

**Scenario**: Implement user registration

**Typical Plan**:

```
1. Create UserRegistrationController
2. Create UserRegistrationService
3. Create UserRegistrationRequest DTO
4. Implement email validation
5. Save user to database
```

**REASONS Canvas (Operations excerpt)**:

```markdown
### Create UserRegistrationService - `UserRegistrationServiceImpl`

1. **Responsibility**: Handle user registration business logic
2. **Package**: `com.example.user.service.impl`
3. **Implements**: `UserRegistrationService` interface
4. **Dependencies** (constructor injection):
   - `UserRepository userRepository`
   - `EmailValidator emailValidator`
   - `PasswordEncoder passwordEncoder`
5. **Methods**:
   - `register(UserRegistrationRequest request): UserRegistrationResponse`
     - **Input Validation**: Call `emailValidator.validate(request.getEmail())`
     - **Business Logic**:
       1. Check if email exists via `userRepository.existsByEmail()`
       2. If exists, throw `EmailAlreadyExistsException` with message "Email already registered"
       3. Encode password via `passwordEncoder.encode()`
       4. Create User entity with status `PENDING_VERIFICATION`
       5. Save via `userRepository.save()`
     - **Exception Handling**: Let exceptions propagate to GlobalExceptionHandler
6. **Annotations**: `@Service`, `@Transactional`
```

**The difference**: Plan says "what to do", REASONS Canvas specifies "exactly how to do it".

## When to Use OpenSPDD

| Scenario                       | Recommendation     | Reason                                                        |
| ------------------------------ | ------------------ | ------------------------------------------------------------- |
| Enterprise feature development | Highly recommended | Design-implementation traceability, long-term maintainability |
| Team collaboration             | Highly recommended | Unified AI coding standards, reduced style conflicts          |
| Complex refactoring            | Recommended        | Strict Operations order prevents dependency chaos             |
| Cross-tool workflows           | Recommended        | Same REASONS Canvas works across different AI tools           |
| Quick prototypes               | Consider           | May be overhead, but valuable if maintenance is needed        |
| One-off scripts                | Not recommended    | ROI too low                                                   |

## Building from Source

```bash
git clone https://github.com/gszhangwei/open-spdd.git
cd open-spdd
go build -o openspdd ./cmd/openspdd
go install ./cmd/openspdd
```

## Testing

```bash
# Run all tests
go test ./tests/...

# Run with verbose output
go test ./tests/... -v

# Run specific module tests
go test ./tests/detector/...
go test ./tests/templates/...
```

## License

[MIT License](LICENSE)
