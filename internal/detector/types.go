package detector

// AIToolType represents the type of AI coding assistant tool.
type AIToolType string

const (
	Cursor      AIToolType = "cursor"
	ClaudeCode  AIToolType = "claude-code"
	Antigravity AIToolType = "antigravity"
	Unknown     AIToolType = "unknown"
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
	default:
		return nil
	}
}

// DetectResult holds the result of environment detection.
type DetectResult struct {
	ToolType   AIToolType
	ConfigPath string
	IsValid    bool
	Message    string
}
