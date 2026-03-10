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
	GenerateForCopilot(targetDir string, force bool) []GenerateResult
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

const (
	SPDDMarkerStart = "<!-- openspdd:start -->"
	SPDDMarkerEnd   = "<!-- openspdd:end -->"
)

// GenerateForCopilot generates the complete Copilot file structure with marker-based merge support.
func (m *EmbeddedTemplateManager) GenerateForCopilot(targetDir string, force bool) []GenerateResult {
	var results []GenerateResult

	githubDir := filepath.Join(targetDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		results = append(results, GenerateResult{
			Success: false,
			Message: "failed to create .github directory: " + err.Error(),
			Error:   err,
		})
		return results
	}

	instructionContent, err := fs.ReadFile(embeddedTemplates, "data/copilot-instructions.md")
	if err != nil {
		results = append(results, GenerateResult{
			Success: false,
			Message: "failed to read copilot-instructions template: " + err.Error(),
			Error:   err,
		})
		return results
	}

	instructionPath := filepath.Join(githubDir, "copilot-instructions.md")
	markedContent := SPDDMarkerStart + "\n" + string(instructionContent) + "\n" + SPDDMarkerEnd
	result := m.writeInstructionFile(instructionPath, markedContent, force)
	results = append(results, result)

	promptsDir := filepath.Join(githubDir, "copilot-prompts")
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		results = append(results, GenerateResult{
			Success: false,
			Message: "failed to create copilot-prompts directory: " + err.Error(),
			Error:   err,
		})
		return results
	}

	entries, err := fs.ReadDir(embeddedTemplates, "data")
	if err != nil {
		results = append(results, GenerateResult{
			Success: false,
			Message: "failed to read templates: " + err.Error(),
			Error:   err,
		})
		return results
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if entry.Name() == "copilot-instructions.md" {
			continue
		}

		content, err := fs.ReadFile(embeddedTemplates, "data/"+entry.Name())
		if err != nil {
			results = append(results, GenerateResult{
				Success:  false,
				FilePath: entry.Name(),
				Message:  "failed to read template: " + err.Error(),
				Error:    err,
			})
			continue
		}

		targetPath := filepath.Join(promptsDir, entry.Name())
		if _, err := os.Stat(targetPath); err == nil && !force {
			results = append(results, GenerateResult{
				Success:  false,
				FilePath: targetPath,
				Message:  "file already exists (use --force to overwrite)",
				Error:    internal.ErrFileExists,
			})
			continue
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			results = append(results, GenerateResult{
				Success:  false,
				FilePath: targetPath,
				Message:  "failed to write file: " + err.Error(),
				Error:    err,
			})
			continue
		}

		results = append(results, GenerateResult{
			Success:  true,
			FilePath: targetPath,
			Message:  "template generated successfully",
		})
	}

	return results
}

func (m *EmbeddedTemplateManager) writeInstructionFile(path, markedContent string, force bool) GenerateResult {
	existingContent, err := os.ReadFile(path)

	if err != nil {
		if err := os.WriteFile(path, []byte(markedContent), 0644); err != nil {
			return GenerateResult{
				Success:  false,
				FilePath: path,
				Message:  "failed to write file: " + err.Error(),
				Error:    err,
			}
		}
		return GenerateResult{
			Success:  true,
			FilePath: path,
			Message:  "instruction file created successfully",
		}
	}

	content := string(existingContent)
	hasStartMarker := strings.Contains(content, SPDDMarkerStart)
	hasEndMarker := strings.Contains(content, SPDDMarkerEnd)

	if hasStartMarker && hasEndMarker {
		updatedContent := m.replaceMarkedSection(content, markedContent)
		if err := os.WriteFile(path, []byte(updatedContent), 0644); err != nil {
			return GenerateResult{
				Success:  false,
				FilePath: path,
				Message:  "failed to update file: " + err.Error(),
				Error:    err,
			}
		}
		return GenerateResult{
			Success:  true,
			FilePath: path,
			Message:  "instruction file updated (marked section replaced)",
		}
	}

	if force {
		if err := os.WriteFile(path, []byte(markedContent), 0644); err != nil {
			return GenerateResult{
				Success:  false,
				FilePath: path,
				Message:  "failed to overwrite file: " + err.Error(),
				Error:    err,
			}
		}
		return GenerateResult{
			Success:  true,
			FilePath: path,
			Message:  "instruction file overwritten (--force)",
		}
	}

	return GenerateResult{
		Success:  false,
		FilePath: path,
		Message:  "file exists without SPDD markers. Use --force to overwrite, or manually add markers",
		Error:    internal.ErrExistingFileNoMarkers,
	}
}

func (m *EmbeddedTemplateManager) replaceMarkedSection(content, markedContent string) string {
	startIdx := strings.Index(content, SPDDMarkerStart)
	endIdx := strings.Index(content, SPDDMarkerEnd) + len(SPDDMarkerEnd)

	return content[:startIdx] + markedContent + content[endIdx:]
}
