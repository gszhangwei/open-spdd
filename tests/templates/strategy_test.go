package templates_test

import (
	"sort"
	"testing"

	"github.com/gszhangwei/open-spdd/internal/detector"
	"github.com/gszhangwei/open-spdd/internal/templates"
)

func TestStrategyFor_GitHubCopilotReturnsCopilotStrategy(t *testing.T) {
	mgr := templates.NewEmbeddedTemplateManager()
	got := templates.StrategyFor(detector.GitHubCopilot, mgr)
	if _, ok := got.(*templates.CopilotInstructionFileStrategy); !ok {
		t.Errorf("StrategyFor(GitHubCopilot) = %T, want *templates.CopilotInstructionFileStrategy", got)
	}
}

func TestStrategyFor_CodexReturnsCodexStrategy(t *testing.T) {
	mgr := templates.NewEmbeddedTemplateManager()
	got := templates.StrategyFor(detector.Codex, mgr)
	if _, ok := got.(*templates.CodexSkillStrategy); !ok {
		t.Errorf("StrategyFor(Codex) = %T, want *templates.CodexSkillStrategy", got)
	}
}

func TestStrategyFor_AntigravityReturnsAntigravityStrategy(t *testing.T) {
	mgr := templates.NewEmbeddedTemplateManager()
	got := templates.StrategyFor(detector.Antigravity, mgr)
	if _, ok := got.(*templates.AntigravitySkillStrategy); !ok {
		t.Errorf("StrategyFor(Antigravity) = %T, want *templates.AntigravitySkillStrategy", got)
	}
}

func TestStrategyFor_FlatMarkdownToolsReturnFlatStrategy(t *testing.T) {
	mgr := templates.NewEmbeddedTemplateManager()
	tools := []detector.AIToolType{
		detector.Cursor,
		detector.ClaudeCode,
		detector.OpenCode,
	}
	for _, tool := range tools {
		t.Run(string(tool), func(t *testing.T) {
			got := templates.StrategyFor(tool, mgr)
			if _, ok := got.(*templates.FlatMarkdownStrategy); !ok {
				t.Errorf("StrategyFor(%v) = %T, want *templates.FlatMarkdownStrategy", tool, got)
			}
		})
	}
}

func TestStrategyFor_UnknownToolFallsBackToFlatStrategy(t *testing.T) {
	mgr := templates.NewEmbeddedTemplateManager()
	got := templates.StrategyFor(detector.Unknown, mgr)
	if _, ok := got.(*templates.FlatMarkdownStrategy); !ok {
		t.Errorf("StrategyFor(Unknown) = %T, want *templates.FlatMarkdownStrategy", got)
	}
}

func TestStrategyRegistry_ContainsExpectedTools(t *testing.T) {
	registered := templates.RegisteredToolsForTest()
	got := make([]string, 0, len(registered))
	for _, tool := range registered {
		got = append(got, string(tool))
	}
	sort.Strings(got)

	want := []string{string(detector.Antigravity), string(detector.Codex), string(detector.GitHubCopilot)}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("RegisteredToolsForTest() = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("RegisteredToolsForTest()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
