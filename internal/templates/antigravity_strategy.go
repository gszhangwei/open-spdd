package templates

import (
	"os"
	"path/filepath"

	"github.com/gszhangwei/open-spdd/internal/detector"
)

type AntigravitySkillStrategy struct {
	manager *EmbeddedTemplateManager
}

func init() {
	RegisterStrategy(detector.Antigravity, func(mgr *EmbeddedTemplateManager) GenerationStrategy {
		return &AntigravitySkillStrategy{manager: mgr}
	})
}

func (s *AntigravitySkillStrategy) GenerateAll(workingDir string, force bool) []GenerateResult {
	tmpls, err := s.manager.ListAvailable(detector.Antigravity)
	if err != nil {
		return []GenerateResult{{
			Success: false,
			Message: "failed to list templates: " + err.Error(),
			Error:   err,
		}}
	}

	results := make([]GenerateResult, 0, len(tmpls))
	for _, tmpl := range tmpls {
		results = append(results, s.GenerateOne(workingDir, tmpl, force)...)
	}
	return results
}

func (s *AntigravitySkillStrategy) GenerateOne(workingDir string, tmpl TemplateMeta, force bool) []GenerateResult {
	baseDir := filepath.Join(workingDir, ".agents", "skills")
	skillDir := filepath.Join(baseDir, tmpl.ID)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return []GenerateResult{{
			Success:  false,
			FilePath: skillDir,
			Message:  "failed to create skill directory: " + err.Error(),
			Error:    err,
		}}
	}

	skillPath := filepath.Join(skillDir, "SKILL.md")
	return []GenerateResult{writeSkillFile(skillPath, buildSkillMd(tmpl), force)}
}
