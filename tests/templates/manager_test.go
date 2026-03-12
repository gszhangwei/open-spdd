package templates_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"open-spdd/internal"
	"open-spdd/internal/detector"
	"open-spdd/internal/templates"
)

func TestNewEmbeddedTemplateManager(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	if manager == nil {
		t.Error("NewEmbeddedTemplateManager() returned nil")
	}
}

func TestEmbeddedTemplateManager_ListCore(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListCore()

	if err != nil {
		t.Fatalf("ListCore() error = %v", err)
	}

	if len(tmpls) == 0 {
		t.Error("ListCore() should return at least one template")
	}

	expectedIDs := []string{"spdd-generate", "spdd-sync", "spdd-reasons-canvas"}
	foundIDs := make(map[string]bool)
	for _, tmpl := range tmpls {
		foundIDs[tmpl.ID] = true
	}

	for _, id := range expectedIDs {
		if !foundIDs[id] {
			t.Errorf("ListCore() missing expected template: %s", id)
		}
	}
}

func TestEmbeddedTemplateManager_ListForTool_Copilot(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListForTool(detector.GitHubCopilot)

	if err != nil {
		t.Fatalf("ListForTool(GitHubCopilot) error = %v", err)
	}

	if len(tmpls) == 0 {
		t.Error("ListForTool(GitHubCopilot) should return copilot-specific templates")
	}

	foundCopilotInstructions := false
	for _, tmpl := range tmpls {
		if tmpl.ID == "copilot-instructions" {
			foundCopilotInstructions = true
			break
		}
	}

	if !foundCopilotInstructions {
		t.Error("ListForTool(GitHubCopilot) should include copilot-instructions template")
	}
}

func TestEmbeddedTemplateManager_ListForTool_Cursor(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListForTool(detector.Cursor)

	if err != nil {
		t.Fatalf("ListForTool(Cursor) error = %v", err)
	}

	for _, tmpl := range tmpls {
		if tmpl.ID == "copilot-instructions" {
			t.Error("ListForTool(Cursor) should NOT include copilot-instructions template")
		}
	}
}

func TestEmbeddedTemplateManager_ListForTool_ClaudeCode(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListForTool(detector.ClaudeCode)

	if err != nil {
		t.Fatalf("ListForTool(ClaudeCode) error = %v", err)
	}

	for _, tmpl := range tmpls {
		if tmpl.ID == "copilot-instructions" {
			t.Error("ListForTool(ClaudeCode) should NOT include copilot-instructions template")
		}
	}
}

func TestEmbeddedTemplateManager_ListForTool_Unknown(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListForTool(detector.Unknown)

	if err != nil {
		t.Fatalf("ListForTool(Unknown) error = %v", err)
	}

	if len(tmpls) != 0 {
		t.Errorf("ListForTool(Unknown) should return empty slice, got %d templates", len(tmpls))
	}
}

func TestEmbeddedTemplateManager_ListOptional(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListOptional()

	if err != nil {
		t.Fatalf("ListOptional() error = %v", err)
	}

	_ = tmpls
}

func TestEmbeddedTemplateManager_ListAvailable_Copilot(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListAvailable(detector.GitHubCopilot)

	if err != nil {
		t.Fatalf("ListAvailable(GitHubCopilot) error = %v", err)
	}

	foundCore := false
	foundCopilot := false
	for _, tmpl := range tmpls {
		if tmpl.ID == "spdd-generate" {
			foundCore = true
		}
		if tmpl.ID == "copilot-instructions" {
			foundCopilot = true
		}
	}

	if !foundCore {
		t.Error("ListAvailable(GitHubCopilot) should include core templates")
	}
	if !foundCopilot {
		t.Error("ListAvailable(GitHubCopilot) should include copilot-specific templates")
	}
}

func TestEmbeddedTemplateManager_ListAvailable_Cursor(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListAvailable(detector.Cursor)

	if err != nil {
		t.Fatalf("ListAvailable(Cursor) error = %v", err)
	}

	foundCore := false
	for _, tmpl := range tmpls {
		if tmpl.ID == "spdd-generate" {
			foundCore = true
		}
		if tmpl.ID == "copilot-instructions" {
			t.Error("ListAvailable(Cursor) should NOT include copilot-instructions")
		}
	}

	if !foundCore {
		t.Error("ListAvailable(Cursor) should include core templates")
	}
}

func TestEmbeddedTemplateManager_ListAvailable_TemplateIsolation(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()

	cursorTemplates, _ := manager.ListAvailable(detector.Cursor)
	copilotTemplates, _ := manager.ListAvailable(detector.GitHubCopilot)

	if len(copilotTemplates) <= len(cursorTemplates) {
		t.Error("ListAvailable(GitHubCopilot) should have more templates than ListAvailable(Cursor)")
	}
}

func TestEmbeddedTemplateManager_ListAll(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListAll()

	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	coreTemplates, _ := manager.ListCore()
	if len(tmpls) < len(coreTemplates) {
		t.Errorf("ListAll() should include at least all core templates, got %d, expected at least %d",
			len(tmpls), len(coreTemplates))
	}

	foundCopilot := false
	for _, tmpl := range tmpls {
		if tmpl.ID == "copilot-instructions" {
			foundCopilot = true
			break
		}
	}
	if !foundCopilot {
		t.Error("ListAll() should include copilot-instructions template")
	}
}

func TestEmbeddedTemplateManager_ListAll_NoDuplicates(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListAll()

	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	seen := make(map[string]bool)
	for _, tmpl := range tmpls {
		if seen[tmpl.ID] {
			t.Errorf("ListAll() returned duplicate template ID: %s", tmpl.ID)
		}
		seen[tmpl.ID] = true
	}
}

func TestEmbeddedTemplateManager_ListAll_Sorted(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpls, err := manager.ListAll()

	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	for i := 1; i < len(tmpls); i++ {
		if tmpls[i-1].Name > tmpls[i].Name {
			t.Errorf("ListAll() not sorted: %s > %s", tmpls[i-1].Name, tmpls[i].Name)
		}
	}
}

func TestEmbeddedTemplateManager_GetByName_Found(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	tmpl, err := manager.GetByName("spdd-generate")

	if err != nil {
		t.Fatalf("GetByName('spdd-generate') error = %v", err)
	}

	if tmpl.ID != "spdd-generate" {
		t.Errorf("GetByName() ID = %v, want spdd-generate", tmpl.ID)
	}
}

func TestEmbeddedTemplateManager_GetByName_CaseInsensitive(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()

	tests := []string{"SPDD-GENERATE", "Spdd-Generate", "spdd-generate"}
	for _, name := range tests {
		tmpl, err := manager.GetByName(name)
		if err != nil {
			t.Errorf("GetByName('%s') error = %v", name, err)
			continue
		}
		if tmpl.ID != "spdd-generate" {
			t.Errorf("GetByName('%s') ID = %v, want spdd-generate", name, tmpl.ID)
		}
	}
}

func TestEmbeddedTemplateManager_GetByName_NotFound(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	_, err := manager.GetByName("nonexistent-template")

	if err == nil {
		t.Error("GetByName('nonexistent-template') should return error")
	}
	if err != internal.ErrTemplateNotFound {
		t.Errorf("GetByName() error = %v, want ErrTemplateNotFound", err)
	}
}

func TestEmbeddedTemplateManager_Generate_Success(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "test-output.md")

	manager := templates.NewEmbeddedTemplateManager()
	result := manager.Generate(templates.GenerateRequest{
		TemplateName: "spdd-generate",
		TargetPath:   targetPath,
		Force:        false,
	})

	if !result.Success {
		t.Errorf("Generate() Success = false, Message = %s", result.Message)
	}
	if result.FilePath != targetPath {
		t.Errorf("Generate() FilePath = %v, want %v", result.FilePath, targetPath)
	}

	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Errorf("Generated file not readable: %v", err)
	}
	if len(content) == 0 {
		t.Error("Generated file is empty")
	}
}

func TestEmbeddedTemplateManager_Generate_CreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "nested", "dir", "test-output.md")

	manager := templates.NewEmbeddedTemplateManager()
	result := manager.Generate(templates.GenerateRequest{
		TemplateName: "spdd-generate",
		TargetPath:   targetPath,
		Force:        false,
	})

	if !result.Success {
		t.Errorf("Generate() Success = false, Message = %s", result.Message)
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Error("Generate() did not create file in nested directory")
	}
}

func TestEmbeddedTemplateManager_Generate_FileExists_NoForce(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "existing.md")

	if err := os.WriteFile(targetPath, []byte("existing content"), 0644); err != nil {
		t.Fatal(err)
	}

	manager := templates.NewEmbeddedTemplateManager()
	result := manager.Generate(templates.GenerateRequest{
		TemplateName: "spdd-generate",
		TargetPath:   targetPath,
		Force:        false,
	})

	if result.Success {
		t.Error("Generate() should fail when file exists without force flag")
	}
	if result.Error != internal.ErrFileExists {
		t.Errorf("Generate() error = %v, want ErrFileExists", result.Error)
	}

	content, _ := os.ReadFile(targetPath)
	if string(content) != "existing content" {
		t.Error("Generate() modified existing file without force flag")
	}
}

func TestEmbeddedTemplateManager_Generate_FileExists_WithForce(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "existing.md")

	if err := os.WriteFile(targetPath, []byte("existing content"), 0644); err != nil {
		t.Fatal(err)
	}

	manager := templates.NewEmbeddedTemplateManager()
	result := manager.Generate(templates.GenerateRequest{
		TemplateName: "spdd-generate",
		TargetPath:   targetPath,
		Force:        true,
	})

	if !result.Success {
		t.Errorf("Generate() with force should succeed, Message = %s", result.Message)
	}

	content, _ := os.ReadFile(targetPath)
	if string(content) == "existing content" {
		t.Error("Generate() with force did not overwrite file")
	}
}

func TestEmbeddedTemplateManager_Generate_TemplateNotFound(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "test-output.md")

	manager := templates.NewEmbeddedTemplateManager()
	result := manager.Generate(templates.GenerateRequest{
		TemplateName: "nonexistent-template",
		TargetPath:   targetPath,
		Force:        false,
	})

	if result.Success {
		t.Error("Generate() should fail for nonexistent template")
	}
}

func TestEmbeddedTemplateManager_Generate_EmptyTargetPath(t *testing.T) {
	manager := templates.NewEmbeddedTemplateManager()
	result := manager.Generate(templates.GenerateRequest{
		TemplateName: "spdd-generate",
		TargetPath:   "",
		Force:        false,
	})

	if result.Success {
		t.Error("Generate() should fail with empty target path")
	}
}

func TestEmbeddedTemplateManager_GenerateForCopilot_Success(t *testing.T) {
	tempDir := t.TempDir()

	manager := templates.NewEmbeddedTemplateManager()
	results := manager.GenerateForCopilot(tempDir, false)

	if len(results) == 0 {
		t.Fatal("GenerateForCopilot() should return at least one result")
	}

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	if successCount == 0 {
		t.Error("GenerateForCopilot() all results failed")
	}

	githubDir := filepath.Join(tempDir, ".github")
	if _, err := os.Stat(githubDir); os.IsNotExist(err) {
		t.Error("GenerateForCopilot() did not create .github directory")
	}

	instructionFile := filepath.Join(githubDir, "copilot-instructions.md")
	if _, err := os.Stat(instructionFile); os.IsNotExist(err) {
		t.Error("GenerateForCopilot() did not create copilot-instructions.md")
	}

	promptsDir := filepath.Join(githubDir, "copilot-prompts")
	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		t.Error("GenerateForCopilot() did not create copilot-prompts directory")
	}
}

func TestEmbeddedTemplateManager_GenerateForCopilot_InstructionFileMarkers(t *testing.T) {
	tempDir := t.TempDir()

	manager := templates.NewEmbeddedTemplateManager()
	manager.GenerateForCopilot(tempDir, false)

	instructionFile := filepath.Join(tempDir, ".github", "copilot-instructions.md")
	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, templates.SPDDMarkerStart) {
		t.Error("GenerateForCopilot() instruction file missing start marker")
	}
	if !strings.Contains(contentStr, templates.SPDDMarkerEnd) {
		t.Error("GenerateForCopilot() instruction file missing end marker")
	}
}

func TestEmbeddedTemplateManager_GenerateForCopilot_MarkerBasedMerge(t *testing.T) {
	tempDir := t.TempDir()

	githubDir := filepath.Join(tempDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		t.Fatal(err)
	}

	instructionFile := filepath.Join(githubDir, "copilot-instructions.md")
	existingContent := `# My Custom Instructions

Some custom content before SPDD.

<!-- openspdd:start -->
old spdd content
<!-- openspdd:end -->

Some custom content after SPDD.
`
	if err := os.WriteFile(instructionFile, []byte(existingContent), 0644); err != nil {
		t.Fatal(err)
	}

	manager := templates.NewEmbeddedTemplateManager()
	results := manager.GenerateForCopilot(tempDir, false)

	var instructionResult *templates.GenerateResult
	for i, r := range results {
		if strings.Contains(r.FilePath, "copilot-instructions.md") {
			instructionResult = &results[i]
			break
		}
	}

	if instructionResult == nil {
		t.Fatal("GenerateForCopilot() did not return result for instruction file")
	}

	if !instructionResult.Success {
		t.Errorf("GenerateForCopilot() instruction file update failed: %s", instructionResult.Message)
	}

	updatedContent, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(updatedContent)

	if !strings.Contains(contentStr, "# My Custom Instructions") {
		t.Error("GenerateForCopilot() lost custom content before markers")
	}
	if !strings.Contains(contentStr, "Some custom content before SPDD.") {
		t.Error("GenerateForCopilot() lost custom content before markers")
	}

	if !strings.Contains(contentStr, "Some custom content after SPDD.") {
		t.Error("GenerateForCopilot() lost custom content after markers")
	}

	if strings.Contains(contentStr, "old spdd content") {
		t.Error("GenerateForCopilot() did not replace old SPDD content")
	}
}

func TestEmbeddedTemplateManager_GenerateForCopilot_FileExistsNoMarkers_NoForce(t *testing.T) {
	tempDir := t.TempDir()

	githubDir := filepath.Join(tempDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		t.Fatal(err)
	}

	instructionFile := filepath.Join(githubDir, "copilot-instructions.md")
	existingContent := `# My Custom Instructions

This file has no SPDD markers.
`
	if err := os.WriteFile(instructionFile, []byte(existingContent), 0644); err != nil {
		t.Fatal(err)
	}

	manager := templates.NewEmbeddedTemplateManager()
	results := manager.GenerateForCopilot(tempDir, false)

	var instructionResult *templates.GenerateResult
	for i, r := range results {
		if strings.Contains(r.FilePath, "copilot-instructions.md") {
			instructionResult = &results[i]
			break
		}
	}

	if instructionResult == nil {
		t.Fatal("GenerateForCopilot() did not return result for instruction file")
	}

	if instructionResult.Success {
		t.Error("GenerateForCopilot() should fail when file exists without markers and no force flag")
	}

	if instructionResult.Error != internal.ErrExistingFileNoMarkers {
		t.Errorf("GenerateForCopilot() error = %v, want ErrExistingFileNoMarkers", instructionResult.Error)
	}

	content, _ := os.ReadFile(instructionFile)
	if string(content) != existingContent {
		t.Error("GenerateForCopilot() modified file without markers despite no force flag")
	}
}

func TestEmbeddedTemplateManager_GenerateForCopilot_FileExistsNoMarkers_WithForce(t *testing.T) {
	tempDir := t.TempDir()

	githubDir := filepath.Join(tempDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		t.Fatal(err)
	}

	instructionFile := filepath.Join(githubDir, "copilot-instructions.md")
	existingContent := `# My Custom Instructions

This file has no SPDD markers.
`
	if err := os.WriteFile(instructionFile, []byte(existingContent), 0644); err != nil {
		t.Fatal(err)
	}

	manager := templates.NewEmbeddedTemplateManager()
	results := manager.GenerateForCopilot(tempDir, true)

	var instructionResult *templates.GenerateResult
	for i, r := range results {
		if strings.Contains(r.FilePath, "copilot-instructions.md") {
			instructionResult = &results[i]
			break
		}
	}

	if instructionResult == nil {
		t.Fatal("GenerateForCopilot() did not return result for instruction file")
	}

	if !instructionResult.Success {
		t.Errorf("GenerateForCopilot() with force should succeed: %s", instructionResult.Message)
	}

	content, _ := os.ReadFile(instructionFile)
	if string(content) == existingContent {
		t.Error("GenerateForCopilot() with force did not overwrite file")
	}
}

func TestEmbeddedTemplateManager_GenerateForCopilot_TemplatesGenerated(t *testing.T) {
	tempDir := t.TempDir()

	manager := templates.NewEmbeddedTemplateManager()
	manager.GenerateForCopilot(tempDir, false)

	promptsDir := filepath.Join(tempDir, ".github", "copilot-prompts")

	expectedTemplates := []string{"spdd-generate.md", "spdd-sync.md", "spdd-reasons-canvas.md"}
	for _, tmpl := range expectedTemplates {
		path := filepath.Join(promptsDir, tmpl)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("GenerateForCopilot() did not generate %s", tmpl)
		}
	}
}

func TestSPDDMarkerConstants(t *testing.T) {
	if templates.SPDDMarkerStart != "<!-- openspdd:start -->" {
		t.Errorf("SPDDMarkerStart = %v, want '<!-- openspdd:start -->'", templates.SPDDMarkerStart)
	}
	if templates.SPDDMarkerEnd != "<!-- openspdd:end -->" {
		t.Errorf("SPDDMarkerEnd = %v, want '<!-- openspdd:end -->'", templates.SPDDMarkerEnd)
	}
}

func TestTemplateManager_Interface(t *testing.T) {
	var _ templates.TemplateManager = (*templates.EmbeddedTemplateManager)(nil)
}
