package detector_test

import (
	"testing"

	"open-spdd/internal/detector"
)

func TestAIToolType_String(t *testing.T) {
	tests := []struct {
		name     string
		toolType detector.AIToolType
		want     string
	}{
		{
			name:     "Cursor returns human-readable name",
			toolType: detector.Cursor,
			want:     "Cursor",
		},
		{
			name:     "ClaudeCode returns human-readable name",
			toolType: detector.ClaudeCode,
			want:     "Claude Code",
		},
		{
			name:     "Antigravity returns human-readable name",
			toolType: detector.Antigravity,
			want:     "Antigravity",
		},
		{
			name:     "GitHubCopilot returns human-readable name",
			toolType: detector.GitHubCopilot,
			want:     "GitHub Copilot",
		},
		{
			name:     "Unknown returns Unknown",
			toolType: detector.Unknown,
			want:     "Unknown",
		},
		{
			name:     "Invalid type returns Unknown",
			toolType: detector.AIToolType("invalid"),
			want:     "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.toolType.String(); got != tt.want {
				t.Errorf("AIToolType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIToolType_GetConfigDir(t *testing.T) {
	tests := []struct {
		name     string
		toolType detector.AIToolType
		want     string
	}{
		{
			name:     "Cursor config directory",
			toolType: detector.Cursor,
			want:     ".cursor/commands",
		},
		{
			name:     "ClaudeCode config directory",
			toolType: detector.ClaudeCode,
			want:     ".claude/commands",
		},
		{
			name:     "Antigravity config directory",
			toolType: detector.Antigravity,
			want:     ".antigravity/commands",
		},
		{
			name:     "GitHubCopilot config directory",
			toolType: detector.GitHubCopilot,
			want:     ".github/copilot-prompts",
		},
		{
			name:     "Unknown returns empty string",
			toolType: detector.Unknown,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.toolType.GetConfigDir(); got != tt.want {
				t.Errorf("AIToolType.GetConfigDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIToolType_GetSignatureFiles(t *testing.T) {
	tests := []struct {
		name     string
		toolType detector.AIToolType
		want     []string
	}{
		{
			name:     "Cursor signature files",
			toolType: detector.Cursor,
			want:     []string{".cursor", ".cursorrules"},
		},
		{
			name:     "ClaudeCode signature files",
			toolType: detector.ClaudeCode,
			want:     []string{".claude", "CLAUDE.md"},
		},
		{
			name:     "Antigravity signature files",
			toolType: detector.Antigravity,
			want:     []string{".antigravity"},
		},
		{
			name:     "GitHubCopilot signature files",
			toolType: detector.GitHubCopilot,
			want:     []string{".github/copilot-instructions.md", ".github/copilot-prompts"},
		},
		{
			name:     "Unknown returns nil",
			toolType: detector.Unknown,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.toolType.GetSignatureFiles()
			if len(got) != len(tt.want) {
				t.Errorf("AIToolType.GetSignatureFiles() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("AIToolType.GetSignatureFiles()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestAIToolType_GetInstructionFile(t *testing.T) {
	tests := []struct {
		name     string
		toolType detector.AIToolType
		want     string
	}{
		{
			name:     "GitHubCopilot has instruction file",
			toolType: detector.GitHubCopilot,
			want:     ".github/copilot-instructions.md",
		},
		{
			name:     "Cursor has no instruction file",
			toolType: detector.Cursor,
			want:     "",
		},
		{
			name:     "ClaudeCode has no instruction file",
			toolType: detector.ClaudeCode,
			want:     "",
		},
		{
			name:     "Antigravity has no instruction file",
			toolType: detector.Antigravity,
			want:     "",
		},
		{
			name:     "Unknown has no instruction file",
			toolType: detector.Unknown,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.toolType.GetInstructionFile(); got != tt.want {
				t.Errorf("AIToolType.GetInstructionFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIToolType_HasInstructionFile(t *testing.T) {
	tests := []struct {
		name     string
		toolType detector.AIToolType
		want     bool
	}{
		{
			name:     "GitHubCopilot has instruction file",
			toolType: detector.GitHubCopilot,
			want:     true,
		},
		{
			name:     "Cursor has no instruction file",
			toolType: detector.Cursor,
			want:     false,
		},
		{
			name:     "ClaudeCode has no instruction file",
			toolType: detector.ClaudeCode,
			want:     false,
		},
		{
			name:     "Antigravity has no instruction file",
			toolType: detector.Antigravity,
			want:     false,
		},
		{
			name:     "Unknown has no instruction file",
			toolType: detector.Unknown,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.toolType.HasInstructionFile(); got != tt.want {
				t.Errorf("AIToolType.HasInstructionFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIToolType_GetToolDirName(t *testing.T) {
	tests := []struct {
		name     string
		toolType detector.AIToolType
		want     string
	}{
		{
			name:     "Cursor tool directory name",
			toolType: detector.Cursor,
			want:     "cursor",
		},
		{
			name:     "ClaudeCode tool directory name",
			toolType: detector.ClaudeCode,
			want:     "claude-code",
		},
		{
			name:     "Antigravity tool directory name",
			toolType: detector.Antigravity,
			want:     "antigravity",
		},
		{
			name:     "GitHubCopilot tool directory name",
			toolType: detector.GitHubCopilot,
			want:     "copilot",
		},
		{
			name:     "Unknown returns empty string",
			toolType: detector.Unknown,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.toolType.GetToolDirName(); got != tt.want {
				t.Errorf("AIToolType.GetToolDirName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectResult_Fields(t *testing.T) {
	result := detector.DetectResult{
		ToolType:   detector.Cursor,
		ConfigPath: "/path/to/.cursor/commands",
		IsValid:    true,
		Message:    "Cursor environment detected",
	}

	if result.ToolType != detector.Cursor {
		t.Errorf("DetectResult.ToolType = %v, want %v", result.ToolType, detector.Cursor)
	}
	if result.ConfigPath != "/path/to/.cursor/commands" {
		t.Errorf("DetectResult.ConfigPath = %v, want %v", result.ConfigPath, "/path/to/.cursor/commands")
	}
	if !result.IsValid {
		t.Error("DetectResult.IsValid = false, want true")
	}
	if result.Message != "Cursor environment detected" {
		t.Errorf("DetectResult.Message = %v, want %v", result.Message, "Cursor environment detected")
	}
}
