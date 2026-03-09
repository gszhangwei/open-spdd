package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"open-spdd/internal"
)

// TemplateManager defines the interface for template operations.
type TemplateManager interface {
	ListAll() ([]TemplateMeta, error)
	GetByName(name string) (TemplateMeta, error)
	Generate(req GenerateRequest) GenerateResult
}

// EmbeddedTemplateManager implements TemplateManager using embedded templates.
type EmbeddedTemplateManager struct{}

// NewEmbeddedTemplateManager creates a new EmbeddedTemplateManager instance.
func NewEmbeddedTemplateManager() *EmbeddedTemplateManager {
	return &EmbeddedTemplateManager{}
}

// ListAll returns all available templates sorted by name.
func (m *EmbeddedTemplateManager) ListAll() ([]TemplateMeta, error) {
	entries, err := fs.ReadDir(embeddedTemplates, "data")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded templates: %w", err)
	}

	var templates []TemplateMeta
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := fs.ReadFile(embeddedTemplates, "data/"+entry.Name())
		if err != nil {
			continue
		}

		meta := ParseFrontmatter(string(content))
		if meta.ID == "" {
			meta.ID = strings.TrimSuffix(entry.Name(), ".md")
		}
		templates = append(templates, meta)
	}

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

// GetByName returns a template by its name (case-insensitive).
func (m *EmbeddedTemplateManager) GetByName(name string) (TemplateMeta, error) {
	templates, err := m.ListAll()
	if err != nil {
		return TemplateMeta{}, err
	}

	nameLower := strings.ToLower(name)
	for _, t := range templates {
		if strings.ToLower(t.Name) == nameLower || strings.ToLower(t.ID) == nameLower {
			return t, nil
		}
	}

	return TemplateMeta{}, internal.ErrTemplateNotFound
}

// Generate creates a template file at the specified target path.
func (m *EmbeddedTemplateManager) Generate(req GenerateRequest) GenerateResult {
	template, err := m.GetByName(req.TemplateName)
	if err != nil {
		return GenerateResult{
			Success: false,
			Message: "template not found: " + req.TemplateName,
			Error:   err,
		}
	}

	targetPath := req.TargetPath
	if targetPath == "" {
		return GenerateResult{
			Success: false,
			Message: "target path is required",
			Error:   fmt.Errorf("target path is required"),
		}
	}

	if _, err := os.Stat(targetPath); err == nil && !req.Force {
		return GenerateResult{
			Success:  false,
			FilePath: targetPath,
			Message:  "file already exists (use --force to overwrite)",
			Error:    internal.ErrFileExists,
		}
	}

	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return GenerateResult{
			Success: false,
			Message: "failed to create directory: " + targetDir,
			Error:   fmt.Errorf("failed to create directory: %w", err),
		}
	}

	if err := os.WriteFile(targetPath, []byte(template.Content), 0644); err != nil {
		return GenerateResult{
			Success: false,
			Message: "failed to write file: " + targetPath,
			Error:   fmt.Errorf("failed to write file: %w", err),
		}
	}

	return GenerateResult{
		Success:  true,
		FilePath: targetPath,
		Message:  "template generated successfully",
	}
}
