package templates_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gszhangwei/open-spdd/internal/detector"
	"github.com/gszhangwei/open-spdd/internal/templates"
)

func TestFlatMarkdownStrategy_OpenCodeGenerateOne_RemovesFrontmatterName(t *testing.T) {
	tempDir := t.TempDir()
	mgr := templates.NewEmbeddedTemplateManager()

	tmpl, err := mgr.GetByName("/spdd-analysis")
	if err != nil {
		t.Fatalf("GetByName(/spdd-analysis) error = %v", err)
	}

	results := templates.StrategyFor(detector.OpenCode, mgr).GenerateOne(tempDir, tmpl, false)
	if len(results) != 1 {
		t.Fatalf("GenerateOne() returned %d results, want 1", len(results))
	}
	if !results[0].Success {
		t.Fatalf("GenerateOne() failed: %s", results[0].Message)
	}

	path := filepath.Join(tempDir, ".opencode", "commands", tmpl.ID+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	content := string(raw)

	if strings.Contains(content, "\nname:") || strings.HasPrefix(content, "name:") {
		t.Fatalf("generated OpenCode template contains forbidden frontmatter name field")
	}
	if !strings.Contains(content, "id: "+tmpl.ID) {
		t.Fatalf("generated OpenCode template missing id line")
	}
	if !strings.Contains(content, "Analyze a business requirement document") {
		t.Fatalf("generated OpenCode template body was not preserved")
	}
}

func TestFlatMarkdownStrategy_OpenCodeGenerateAll_RemovesFrontmatterName(t *testing.T) {
	tempDir := t.TempDir()
	mgr := templates.NewEmbeddedTemplateManager()

	results := templates.StrategyFor(detector.OpenCode, mgr).GenerateAll(tempDir, false)
	if len(results) == 0 {
		t.Fatal("GenerateAll() returned no results")
	}

	for _, result := range results {
		if !result.Success {
			t.Fatalf("GenerateAll() failed for %s: %s", result.FilePath, result.Message)
		}
		raw, err := os.ReadFile(result.FilePath)
		if err != nil {
			t.Fatalf("read generated file %s: %v", result.FilePath, err)
		}
		content := string(raw)
		if strings.Contains(content, "\nname:") || strings.HasPrefix(content, "name:") {
			t.Fatalf("generated OpenCode template %s contains forbidden frontmatter name field", result.FilePath)
		}
	}
}

func TestFlatMarkdownStrategy_NonOpenCode_PreservesFrontmatterName(t *testing.T) {
	tempDir := t.TempDir()
	mgr := templates.NewEmbeddedTemplateManager()

	tmpl, err := mgr.GetByName("/spdd-analysis")
	if err != nil {
		t.Fatalf("GetByName(/spdd-analysis) error = %v", err)
	}

	results := templates.StrategyFor(detector.Cursor, mgr).GenerateOne(tempDir, tmpl, false)
	if len(results) != 1 {
		t.Fatalf("GenerateOne() returned %d results, want 1", len(results))
	}
	if !results[0].Success {
		t.Fatalf("GenerateOne() failed: %s", results[0].Message)
	}

	path := filepath.Join(tempDir, ".cursor", "commands", tmpl.ID+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	content := string(raw)

	if !strings.Contains(content, "name: spdd-analysis") {
		t.Fatalf("non-OpenCode generation should preserve name field")
	}
}

func TestFlatMarkdownStrategy_OpenCodeForceRegeneration_IsIdempotent(t *testing.T) {
	tempDir := t.TempDir()
	mgr := templates.NewEmbeddedTemplateManager()

	tmpl, err := mgr.GetByName("/spdd-analysis")
	if err != nil {
		t.Fatalf("GetByName(/spdd-analysis) error = %v", err)
	}

	first := templates.StrategyFor(detector.OpenCode, mgr).GenerateOne(tempDir, tmpl, false)
	if len(first) != 1 || !first[0].Success {
		t.Fatalf("first GenerateOne() failed: %+v", first)
	}

	second := templates.StrategyFor(detector.OpenCode, mgr).GenerateOne(tempDir, tmpl, true)
	if len(second) != 1 || !second[0].Success {
		t.Fatalf("second GenerateOne() with force failed: %+v", second)
	}

	path := filepath.Join(tempDir, ".opencode", "commands", tmpl.ID+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	content := string(raw)

	if strings.Count(content, "---") < 2 {
		t.Fatalf("frontmatter delimiters malformed after force regeneration")
	}
	if strings.Contains(content, "\nname:") || strings.HasPrefix(content, "name:") {
		t.Fatalf("name field reappeared after force regeneration")
	}
}
