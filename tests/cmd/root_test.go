package cmd_test

import (
	"strings"
	"testing"

	"github.com/gszhangwei/open-spdd/cmd"
	"github.com/gszhangwei/open-spdd/internal/detector"
)

func TestParseToolFlag_Cursor(t *testing.T) {
	tests := []struct {
		input string
		want  detector.AIToolType
	}{
		{"cursor", detector.Cursor},
		{"Cursor", detector.Cursor},
		{"CURSOR", detector.Cursor},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := cmd.ParseToolFlag(tt.input); got != tt.want {
				t.Errorf("ParseToolFlag(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseToolFlag_ClaudeCode(t *testing.T) {
	tests := []struct {
		input string
		want  detector.AIToolType
	}{
		{"claude-code", detector.ClaudeCode},
		{"Claude-Code", detector.ClaudeCode},
		{"CLAUDE-CODE", detector.ClaudeCode},
		{"claude", detector.ClaudeCode},
		{"Claude", detector.ClaudeCode},
		{"CLAUDE", detector.ClaudeCode},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := cmd.ParseToolFlag(tt.input); got != tt.want {
				t.Errorf("ParseToolFlag(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseToolFlag_Antigravity(t *testing.T) {
	tests := []struct {
		input string
		want  detector.AIToolType
	}{
		{"antigravity", detector.Antigravity},
		{"Antigravity", detector.Antigravity},
		{"ANTIGRAVITY", detector.Antigravity},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := cmd.ParseToolFlag(tt.input); got != tt.want {
				t.Errorf("ParseToolFlag(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseToolFlag_GitHubCopilot(t *testing.T) {
	tests := []struct {
		input string
		want  detector.AIToolType
	}{
		{"github-copilot", detector.GitHubCopilot},
		{"GitHub-Copilot", detector.GitHubCopilot},
		{"GITHUB-COPILOT", detector.GitHubCopilot},
		{"copilot", detector.GitHubCopilot},
		{"Copilot", detector.GitHubCopilot},
		{"COPILOT", detector.GitHubCopilot},
		{"gh-copilot", detector.GitHubCopilot},
		{"GH-Copilot", detector.GitHubCopilot},
		{"GH-COPILOT", detector.GitHubCopilot},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := cmd.ParseToolFlag(tt.input); got != tt.want {
				t.Errorf("ParseToolFlag(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseToolFlag_OpenCode(t *testing.T) {
	tests := []struct {
		input string
		want  detector.AIToolType
	}{
		{"opencode", detector.OpenCode},
		{"OpenCode", detector.OpenCode},
		{"OPENCODE", detector.OpenCode},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := cmd.ParseToolFlag(tt.input); got != tt.want {
				t.Errorf("ParseToolFlag(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseToolFlag_Codex(t *testing.T) {
	tests := []struct {
		input string
		want  detector.AIToolType
	}{
		{"codex", detector.Codex},
		{"Codex", detector.Codex},
		{"CODEX", detector.Codex},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := cmd.ParseToolFlag(tt.input); got != tt.want {
				t.Errorf("ParseToolFlag(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseToolFlag_Unknown(t *testing.T) {
	tests := []struct {
		input string
		want  detector.AIToolType
	}{
		{"unknown", detector.Unknown},
		{"invalid", detector.Unknown},
		{"", detector.Unknown},
		{"vscode", detector.Unknown},
		{"jetbrains", detector.Unknown},
		{"sublime", detector.Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := cmd.ParseToolFlag(tt.input); got != tt.want {
				t.Errorf("ParseToolFlag(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseToolFlag_AllValidToolsReturnValidType(t *testing.T) {
	validTools := []string{
		"cursor",
		"claude-code", "claude",
		"antigravity",
		"github-copilot", "copilot", "gh-copilot",
		"opencode",
		"codex",
	}

	for _, tool := range validTools {
		t.Run(tool, func(t *testing.T) {
			result := cmd.ParseToolFlag(tool)
			if result == detector.Unknown {
				t.Errorf("ParseToolFlag(%q) returned Unknown, expected valid type", tool)
			}
		})
	}
}

func TestValidateToolFlag(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"", false},
		{"cursor", false},
		{"Claude-Code", false},
		{"copilot", false},
		{"codex", false},
		{"bogus", true},
		{"vscode", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := cmd.ValidateToolFlag(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateToolFlag(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "accepted: cursor, claude-code") {
				t.Errorf("ValidateToolFlag(%q) error %q should list the accepted tools", tt.input, err)
			}
		})
	}
}
