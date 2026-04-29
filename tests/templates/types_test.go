package templates_test

import (
	"testing"

	"github.com/gszhangwei/open-spdd/internal/templates"
)

func TestParseFrontmatter_ValidFrontmatter(t *testing.T) {
	content := `---
name: Test Template
id: test-template
category: Development
description: A test template for unit testing
---

This is the template content.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "Test Template" {
		t.Errorf("ParseFrontmatter() Name = %v, want %v", meta.Name, "Test Template")
	}
	if meta.ID != "test-template" {
		t.Errorf("ParseFrontmatter() ID = %v, want %v", meta.ID, "test-template")
	}
	if meta.Category != "Development" {
		t.Errorf("ParseFrontmatter() Category = %v, want %v", meta.Category, "Development")
	}
	if meta.Description != "A test template for unit testing" {
		t.Errorf("ParseFrontmatter() Description = %v, want %v", meta.Description, "A test template for unit testing")
	}
	if meta.Content != content {
		t.Error("ParseFrontmatter() Content should contain the full original content")
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	content := `This is a template without frontmatter.

Just plain content.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "" {
		t.Errorf("ParseFrontmatter() Name = %v, want empty", meta.Name)
	}
	if meta.ID != "" {
		t.Errorf("ParseFrontmatter() ID = %v, want empty", meta.ID)
	}
	if meta.Category != "" {
		t.Errorf("ParseFrontmatter() Category = %v, want empty", meta.Category)
	}
	if meta.Description != "" {
		t.Errorf("ParseFrontmatter() Description = %v, want empty", meta.Description)
	}
	if meta.Content != content {
		t.Error("ParseFrontmatter() Content should still contain the full content")
	}
}

func TestParseFrontmatter_PartialFrontmatter(t *testing.T) {
	content := `---
name: Partial Template
description: Only name and description
---

Content here.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "Partial Template" {
		t.Errorf("ParseFrontmatter() Name = %v, want %v", meta.Name, "Partial Template")
	}
	if meta.ID != "" {
		t.Errorf("ParseFrontmatter() ID = %v, want empty", meta.ID)
	}
	if meta.Category != "" {
		t.Errorf("ParseFrontmatter() Category = %v, want empty", meta.Category)
	}
	if meta.Description != "Only name and description" {
		t.Errorf("ParseFrontmatter() Description = %v, want %v", meta.Description, "Only name and description")
	}
}

func TestParseFrontmatter_IncompleteFrontmatter(t *testing.T) {
	content := `---
name: Incomplete
This line has no colon separator
---

Content.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "Incomplete" {
		t.Errorf("ParseFrontmatter() Name = %v, want %v", meta.Name, "Incomplete")
	}
}

func TestParseFrontmatter_OnlyOpeningDelimiter(t *testing.T) {
	content := `---
name: Broken
This never closes properly.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "" {
		t.Errorf("ParseFrontmatter() Name = %v, want empty for broken frontmatter", meta.Name)
	}
	if meta.Content != content {
		t.Error("ParseFrontmatter() Content should still be set")
	}
}

func TestParseFrontmatter_EmptyFrontmatter(t *testing.T) {
	content := `---
---

Content only.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "" {
		t.Errorf("ParseFrontmatter() Name = %v, want empty", meta.Name)
	}
	if meta.Content != content {
		t.Error("ParseFrontmatter() Content should be set")
	}
}

func TestParseFrontmatter_WhitespaceHandling(t *testing.T) {
	content := `---
name:   Whitespace Test   
  id: spaced-id
category:    Category   
---

Content.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "Whitespace Test" {
		t.Errorf("ParseFrontmatter() Name = %v, want %v", meta.Name, "Whitespace Test")
	}
	if meta.ID != "spaced-id" {
		t.Errorf("ParseFrontmatter() ID = %v, want %v", meta.ID, "spaced-id")
	}
	if meta.Category != "Category" {
		t.Errorf("ParseFrontmatter() Category = %v, want %v", meta.Category, "Category")
	}
}

func TestParseFrontmatter_ColonInValue(t *testing.T) {
	content := `---
name: Template: With Colon
description: A description: with a colon inside
---

Content.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "Template: With Colon" {
		t.Errorf("ParseFrontmatter() Name = %v, want %v", meta.Name, "Template: With Colon")
	}
	if meta.Description != "A description: with a colon inside" {
		t.Errorf("ParseFrontmatter() Description = %v, want %v", meta.Description, "A description: with a colon inside")
	}
}

func TestParseFrontmatter_EmptyContent(t *testing.T) {
	content := ""
	meta := templates.ParseFrontmatter(content)

	if meta.Content != "" {
		t.Errorf("ParseFrontmatter() Content = %v, want empty", meta.Content)
	}
}

func TestParseFrontmatter_RealWorldTemplate(t *testing.T) {
	content := `---
name: /spdd-generate
id: spdd-generate
category: Development
description: Generate code from a structured SPDD prompt file following the REASONS Canvas methodology
---

Generate implementation code from a structured SPDD (Structured Prompt-Driven Development) prompt file.

**Input**: The argument after ` + "`/spdd-generate`" + ` is the path to the structured prompt file.`

	meta := templates.ParseFrontmatter(content)

	if meta.Name != "/spdd-generate" {
		t.Errorf("ParseFrontmatter() Name = %v, want %v", meta.Name, "/spdd-generate")
	}
	if meta.ID != "spdd-generate" {
		t.Errorf("ParseFrontmatter() ID = %v, want %v", meta.ID, "spdd-generate")
	}
	if meta.Category != "Development" {
		t.Errorf("ParseFrontmatter() Category = %v, want %v", meta.Category, "Development")
	}
}

func TestGenerateRequest_Fields(t *testing.T) {
	req := templates.GenerateRequest{
		TemplateName: "test-template",
		TargetPath:   "/path/to/output.md",
		Force:        true,
	}

	if req.TemplateName != "test-template" {
		t.Errorf("GenerateRequest.TemplateName = %v, want %v", req.TemplateName, "test-template")
	}
	if req.TargetPath != "/path/to/output.md" {
		t.Errorf("GenerateRequest.TargetPath = %v, want %v", req.TargetPath, "/path/to/output.md")
	}
	if !req.Force {
		t.Error("GenerateRequest.Force = false, want true")
	}
}

func TestGenerateResult_Success(t *testing.T) {
	result := templates.GenerateResult{
		Success:  true,
		FilePath: "/path/to/file.md",
		Message:  "template generated successfully",
		Error:    nil,
	}

	if !result.Success {
		t.Error("GenerateResult.Success = false, want true")
	}
	if result.FilePath != "/path/to/file.md" {
		t.Errorf("GenerateResult.FilePath = %v, want %v", result.FilePath, "/path/to/file.md")
	}
	if result.Error != nil {
		t.Errorf("GenerateResult.Error = %v, want nil", result.Error)
	}
}

func TestGenerateResult_Failure(t *testing.T) {
	err := &testError{message: "test error"}
	result := templates.GenerateResult{
		Success:  false,
		FilePath: "/path/to/file.md",
		Message:  "failed to generate",
		Error:    err,
	}

	if result.Success {
		t.Error("GenerateResult.Success = true, want false")
	}
	if result.Error == nil {
		t.Error("GenerateResult.Error = nil, want error")
	}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

func TestTemplateMeta_Fields(t *testing.T) {
	meta := templates.TemplateMeta{
		Name:        "Test Template",
		ID:          "test",
		Category:    "Testing",
		Description: "A test template",
		Content:     "Template content here",
		Tags:        []string{"test", "unit"},
	}

	if meta.Name != "Test Template" {
		t.Errorf("TemplateMeta.Name = %v, want %v", meta.Name, "Test Template")
	}
	if meta.ID != "test" {
		t.Errorf("TemplateMeta.ID = %v, want %v", meta.ID, "test")
	}
	if len(meta.Tags) != 2 {
		t.Errorf("TemplateMeta.Tags length = %v, want 2", len(meta.Tags))
	}
}
