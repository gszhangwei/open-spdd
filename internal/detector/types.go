package detector

// AIToolType represents the type of AI coding assistant tool.
type AIToolType string

const (
	Cursor        AIToolType = "cursor"
	ClaudeCode    AIToolType = "claude-code"
	Antigravity   AIToolType = "antigravity"
	GitHubCopilot AIToolType = "github-copilot"
	Unknown       AIToolType = "unknown"
)

// String returns the human-readable name of the tool type.
func (t AIToolType) String() string {
	switch t {
	case Cursor:
		return "Cursor"
	case ClaudeCode:
		return "Claude Code"
	case Antigravity:
		return "Antigravity"
	case GitHubCopilot:
		return "GitHub Copilot"
	default:
		return "Unknown"
	}
}

// GetConfigDir returns the config directory name for each tool type.
func (t AIToolType) GetConfigDir() string {
	switch t {
	case Cursor:
		return ".cursor/commands"
	case ClaudeCode:
		return ".claude/commands"
	case Antigravity:
		return ".antigravity/commands"
	case GitHubCopilot:
		return ".github/copilot-prompts"
	default:
		return ""
	}
}

// GetSignatureFiles returns the list of signature files/directories to detect.
func (t AIToolType) GetSignatureFiles() []string {
	switch t {
	case Cursor:
		return []string{".cursor", ".cursorrules"}
	case ClaudeCode:
		return []string{".claude", "CLAUDE.md"}
	case Antigravity:
		return []string{".antigravity"}
	case GitHubCopilot:
		return []string{".github"}
	default:
		return nil
	}
}

// GetInstructionFile returns the path to the instruction file for tools that use one.
func (t AIToolType) GetInstructionFile() string {
	switch t {
	case GitHubCopilot:
		return ".github/copilot-instructions.md"
	default:
		return ""
	}
}

// HasInstructionFile returns whether this tool type uses a separate instruction file.
func (t AIToolType) HasInstructionFile() bool {
	return t == GitHubCopilot
}

// GetToolDirName returns the directory name under tools/ for this tool type.
func (t AIToolType) GetToolDirName() string {
	switch t {
	case Cursor:
		return "cursor"
	case ClaudeCode:
		return "claude-code"
	case Antigravity:
		return "antigravity"
	case GitHubCopilot:
		return "copilot"
	default:
		return ""
	}
}

// DetectResult holds the result of environment detection.
type DetectResult struct {
	ToolType   AIToolType
	ConfigPath string
	IsValid    bool
	Message    string
}
