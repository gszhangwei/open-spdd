package detector_test

import (
	"os"
	"path/filepath"
	"testing"

	"open-spdd/internal/detector"
)

func TestNewDefaultDetector(t *testing.T) {
	det := detector.NewDefaultDetector()
	if det == nil {
		t.Error("NewDefaultDetector() returned nil")
	}
}

func TestDefaultDetector_Detect_CursorEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	cursorDir := filepath.Join(tempDir, ".cursor")
	if err := os.MkdirAll(cursorDir, 0755); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if !result.IsValid {
		t.Error("Detect() should return valid result for Cursor environment")
	}
	if result.ToolType != detector.Cursor {
		t.Errorf("Detect() ToolType = %v, want %v", result.ToolType, detector.Cursor)
	}
	if result.ConfigPath != filepath.Join(tempDir, ".cursor/commands") {
		t.Errorf("Detect() ConfigPath = %v, want %v", result.ConfigPath, filepath.Join(tempDir, ".cursor/commands"))
	}
}

func TestDefaultDetector_Detect_CursorRulesFile(t *testing.T) {
	tempDir := t.TempDir()
	cursorRules := filepath.Join(tempDir, ".cursorrules")
	if err := os.WriteFile(cursorRules, []byte("rules"), 0644); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if !result.IsValid {
		t.Error("Detect() should return valid result for .cursorrules file")
	}
	if result.ToolType != detector.Cursor {
		t.Errorf("Detect() ToolType = %v, want %v", result.ToolType, detector.Cursor)
	}
}

func TestDefaultDetector_Detect_ClaudeCodeEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	claudeDir := filepath.Join(tempDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if !result.IsValid {
		t.Error("Detect() should return valid result for Claude Code environment")
	}
	if result.ToolType != detector.ClaudeCode {
		t.Errorf("Detect() ToolType = %v, want %v", result.ToolType, detector.ClaudeCode)
	}
	if result.ConfigPath != filepath.Join(tempDir, ".claude/commands") {
		t.Errorf("Detect() ConfigPath = %v, want %v", result.ConfigPath, filepath.Join(tempDir, ".claude/commands"))
	}
}

func TestDefaultDetector_Detect_ClaudeMdFile(t *testing.T) {
	tempDir := t.TempDir()
	claudeMd := filepath.Join(tempDir, "CLAUDE.md")
	if err := os.WriteFile(claudeMd, []byte("# Claude Config"), 0644); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if !result.IsValid {
		t.Error("Detect() should return valid result for CLAUDE.md file")
	}
	if result.ToolType != detector.ClaudeCode {
		t.Errorf("Detect() ToolType = %v, want %v", result.ToolType, detector.ClaudeCode)
	}
}

func TestDefaultDetector_Detect_AntigravityEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	antigravityDir := filepath.Join(tempDir, ".antigravity")
	if err := os.MkdirAll(antigravityDir, 0755); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if !result.IsValid {
		t.Error("Detect() should return valid result for Antigravity environment")
	}
	if result.ToolType != detector.Antigravity {
		t.Errorf("Detect() ToolType = %v, want %v", result.ToolType, detector.Antigravity)
	}
	if result.ConfigPath != filepath.Join(tempDir, ".antigravity/commands") {
		t.Errorf("Detect() ConfigPath = %v, want %v", result.ConfigPath, filepath.Join(tempDir, ".antigravity/commands"))
	}
}

func TestDefaultDetector_Detect_GitHubCopilotEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	githubDir := filepath.Join(tempDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if !result.IsValid {
		t.Error("Detect() should return valid result for GitHub Copilot environment")
	}
	if result.ToolType != detector.GitHubCopilot {
		t.Errorf("Detect() ToolType = %v, want %v", result.ToolType, detector.GitHubCopilot)
	}
	if result.ConfigPath != filepath.Join(tempDir, ".github/copilot-prompts") {
		t.Errorf("Detect() ConfigPath = %v, want %v", result.ConfigPath, filepath.Join(tempDir, ".github/copilot-prompts"))
	}
}

func TestDefaultDetector_Detect_NoEnvironment(t *testing.T) {
	tempDir := t.TempDir()

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if result.IsValid {
		t.Error("Detect() should return invalid result when no environment detected")
	}
	if result.ToolType != detector.Unknown {
		t.Errorf("Detect() ToolType = %v, want %v", result.ToolType, detector.Unknown)
	}
	if result.Message != "no AI coding tool environment detected" {
		t.Errorf("Detect() Message = %v, want %v", result.Message, "no AI coding tool environment detected")
	}
}

func TestDefaultDetector_Detect_Priority_CursorOverClaudeCode(t *testing.T) {
	tempDir := t.TempDir()

	cursorDir := filepath.Join(tempDir, ".cursor")
	claudeDir := filepath.Join(tempDir, ".claude")
	if err := os.MkdirAll(cursorDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if result.ToolType != detector.Cursor {
		t.Errorf("Detect() should prioritize Cursor over Claude Code, got %v", result.ToolType)
	}
}

func TestDefaultDetector_Detect_Priority_ClaudeCodeOverAntigravity(t *testing.T) {
	tempDir := t.TempDir()

	claudeDir := filepath.Join(tempDir, ".claude")
	antigravityDir := filepath.Join(tempDir, ".antigravity")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(antigravityDir, 0755); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if result.ToolType != detector.ClaudeCode {
		t.Errorf("Detect() should prioritize Claude Code over Antigravity, got %v", result.ToolType)
	}
}

func TestDefaultDetector_Detect_Priority_AntigravityOverCopilot(t *testing.T) {
	tempDir := t.TempDir()

	antigravityDir := filepath.Join(tempDir, ".antigravity")
	githubDir := filepath.Join(tempDir, ".github")
	if err := os.MkdirAll(antigravityDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if result.ToolType != detector.Antigravity {
		t.Errorf("Detect() should prioritize Antigravity over GitHub Copilot, got %v", result.ToolType)
	}
}

func TestDefaultDetector_Detect_Priority_AllTools(t *testing.T) {
	tempDir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(tempDir, ".cursor"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tempDir, ".claude"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tempDir, ".antigravity"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tempDir, ".github"), 0755); err != nil {
		t.Fatal(err)
	}

	det := detector.NewDefaultDetector()
	result := det.Detect(tempDir)

	if result.ToolType != detector.Cursor {
		t.Errorf("Detect() should prioritize Cursor when all tools present, got %v", result.ToolType)
	}
}

func TestDefaultDetector_Detect_EmptyWorkingDir(t *testing.T) {
	det := detector.NewDefaultDetector()
	result := det.Detect("")

	if result.ToolType == detector.Unknown && result.Message == "failed to get current working directory" {
		t.Skip("Could not get current working directory")
	}
}

func TestDefaultDetector_GetConfigDirPath(t *testing.T) {
	det := detector.NewDefaultDetector()

	tests := []struct {
		name       string
		tool       detector.AIToolType
		workingDir string
		want       string
	}{
		{
			name:       "Cursor config path",
			tool:       detector.Cursor,
			workingDir: "/project",
			want:       "/project/.cursor/commands",
		},
		{
			name:       "ClaudeCode config path",
			tool:       detector.ClaudeCode,
			workingDir: "/project",
			want:       "/project/.claude/commands",
		},
		{
			name:       "Antigravity config path",
			tool:       detector.Antigravity,
			workingDir: "/project",
			want:       "/project/.antigravity/commands",
		},
		{
			name:       "GitHubCopilot config path",
			tool:       detector.GitHubCopilot,
			workingDir: "/project",
			want:       "/project/.github/copilot-prompts",
		},
		{
			name:       "Unknown tool returns empty",
			tool:       detector.Unknown,
			workingDir: "/project",
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := det.GetConfigDirPath(tt.tool, tt.workingDir)
			if got != tt.want {
				t.Errorf("GetConfigDirPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectorService_Interface(t *testing.T) {
	var _ detector.DetectorService = (*detector.DefaultDetector)(nil)
}
