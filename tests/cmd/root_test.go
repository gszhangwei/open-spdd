package cmd_test

import (
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
