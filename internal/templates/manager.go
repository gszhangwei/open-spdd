package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gszhangwei/open-spdd/internal"
	"github.com/gszhangwei/open-spdd/internal/detector"
)

// TemplateManager defines the interface for template operations.
//
// Per-tool generation behavior (e.g., Copilot's instruction-file merge,
// Codex's skill-bundle output) is intentionally NOT exposed on this interface.
// Use templates.StrategyFor(tool, mgr).GenerateAll(workingDir, force) instead.
type TemplateManager interface {
	ListCore() ([]TemplateMeta, error)
	ListForTool(tool detector.AIToolType) ([]TemplateMeta, error)
	ListOptional() ([]TemplateMeta, error)
	ListAvailable(tool detector.AIToolType) ([]TemplateMeta, error)
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

// loadTemplatesFromDir loads and parses templates from a specific embedded directory path.
func (m *EmbeddedTemplateManager) loadTemplatesFromDir(dir string) ([]TemplateMeta, error) {
	entries, err := fs.ReadDir(embeddedTemplates, dir)
	if err != nil {
		return []TemplateMeta{}, nil
	}

	var templates []TemplateMeta
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := fs.ReadFile(embeddedTemplates, dir+"/"+entry.Name())
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

// ListCore returns all core templates that should be installed by default.
func (m *EmbeddedTemplateManager) ListCore() ([]TemplateMeta, error) {
	return m.loadTemplatesFromDir("data/core")
}

// ListForTool returns tool-specific templates for the given AI tool type.
func (m *EmbeddedTemplateManager) ListForTool(tool detector.AIToolType) ([]TemplateMeta, error) {
	if tool == detector.Unknown {
		return []TemplateMeta{}, nil
	}

	dirName := tool.GetToolDirName()
	if dirName == "" {
		return []TemplateMeta{}, nil
	}

	return m.loadTemplatesFromDir("data/tools/" + dirName)
}

// ListOptional returns all optional templates available for manual selection.
func (m *EmbeddedTemplateManager) ListOptional() ([]TemplateMeta, error) {
	return m.loadTemplatesFromDir("data/optional")
}

// ListAvailable returns all templates available for the current environment.
func (m *EmbeddedTemplateManager) ListAvailable(tool detector.AIToolType) ([]TemplateMeta, error) {
	coreTemplates, err := m.ListCore()
	if err != nil {
		return nil, err
	}

	toolTemplates, err := m.ListForTool(tool)
	if err != nil {
		return nil, err
	}

	templates := append(coreTemplates, toolTemplates...)

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

// ListAll returns ALL templates across all categories (for admin/debug use).
func (m *EmbeddedTemplateManager) ListAll() ([]TemplateMeta, error) {
	coreTemplates, err := m.ListCore()
	if err != nil {
		return nil, err
	}

	templates := coreTemplates

	knownTools := []detector.AIToolType{
		detector.Cursor,
		detector.ClaudeCode,
		detector.Antigravity,
		detector.GitHubCopilot,
		detector.OpenCode,
		detector.Codex,
	}

	for _, tool := range knownTools {
		toolTemplates, err := m.ListForTool(tool)
		if err != nil {
			continue
		}
		templates = append(templates, toolTemplates...)
	}

	optionalTemplates, err := m.ListOptional()
	if err == nil {
		templates = append(templates, optionalTemplates...)
	}

	seen := make(map[string]bool)
	var unique []TemplateMeta
	for _, t := range templates {
		if !seen[t.ID] {
			seen[t.ID] = true
			unique = append(unique, t)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].Name < unique[j].Name
	})

	return unique, nil
}

// GetByName returns a template by its name (case-insensitive).
// A single leading "/" is trimmed for compatibility with slash-prefixed lookups
// (e.g. "/spdd-analysis" matches name/id "spdd-analysis").
func (m *EmbeddedTemplateManager) GetByName(name string) (TemplateMeta, error) {
	templates, err := m.ListAll()
	if err != nil {
		return TemplateMeta{}, err
	}

	nameLower := strings.ToLower(name)
	normalized := strings.TrimPrefix(nameLower, "/")
	for _, t := range templates {
		tName := strings.ToLower(t.Name)
		tID := strings.ToLower(t.ID)
		if tName == nameLower || tID == nameLower || tName == normalized || tID == normalized {
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

	return m.generateWithContent(targetPath, template.Content, req.Force)
}

func (m *EmbeddedTemplateManager) generateWithContent(targetPath, content string, force bool) GenerateResult {
	if _, err := os.Stat(targetPath); err == nil && !force {
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

	if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
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
