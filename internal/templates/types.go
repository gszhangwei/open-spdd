package templates

import (
	"strings"
)

// TemplateMeta holds metadata parsed from template YAML frontmatter.
type TemplateMeta struct {
	Name        string
	ID          string
	Category    string
	Description string
	Content     string
	Tags        []string
}

// ParseFrontmatter extracts metadata from YAML frontmatter in template content.
func ParseFrontmatter(content string) TemplateMeta {
	meta := TemplateMeta{
		Content: content,
	}

	if !strings.HasPrefix(content, "---") {
		return meta
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return meta
	}

	frontmatter := parts[1]
	lines := strings.Split(frontmatter, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])

		switch key {
		case "name":
			meta.Name = value
		case "id":
			meta.ID = value
		case "category":
			meta.Category = value
		case "description":
			meta.Description = value
		}
	}

	return meta
}

// GenerateRequest holds the request parameters for template generation.
type GenerateRequest struct {
	TemplateName string
	TargetPath   string
	Force        bool
}

// GenerateResult holds the result of a template generation operation.
type GenerateResult struct {
	Success  bool
	FilePath string
	Message  string
	Error    error
}
