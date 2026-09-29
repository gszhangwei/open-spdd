package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/gszhangwei/open-spdd/internal/detector"
	"github.com/gszhangwei/open-spdd/internal/templates"
)

var (
	forceFlag         bool
	allFlag           bool
	outputFlag        string
	allowImplicitFlag bool
)

var generateCmd = &cobra.Command{
	Use:     "generate [template-name]",
	Aliases: []string{"gen", "g"},
	Short:   "Generate command template file",
	Long: `Generate a command template file to the detected AI tool's config directory.
If no template name is specified, an interactive selection will be shown.`,
	Run: func(cmd *cobra.Command, args []string) {
		if outputFlag == "" && toolFlag == "" {
			tool := selectToolInteractively()
			if tool == detector.Unknown {
				uiRenderer.RenderError("No tool selected. Use --output or --tool flag.")
				return
			}
			workingDir, _ := os.Getwd()
			detectedResult = detector.DetectResult{
				ToolType:   tool,
				ConfigPath: det.GetConfigDirPath(tool, workingDir),
				IsValid:    true,
				Message:    "tool manually selected: " + tool.String(),
			}
		}

		targetDir := determineTargetDir()
		if targetDir == "" {
			uiRenderer.RenderError("Could not determine target directory. Use --output or --tool flag.")
			return
		}

		if allFlag {
			generateAllTemplates(targetDir)
			return
		}

		if len(args) > 0 {
			generateSingleTemplate(args[0])
			return
		}

		generateInteractively()
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Overwrite existing files")
	generateCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Generate all available templates")
	generateCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Custom output directory (overrides detection)")
	generateCmd.Flags().BoolVar(&allowImplicitFlag, "allow-implicit", false, "Allow implicit invocation of generated Codex skills (Codex only)")
	rootCmd.AddCommand(generateCmd)
}

func determineTargetDir() string {
	if outputFlag != "" {
		return outputFlag
	}

	if detectedResult.IsValid && detectedResult.ConfigPath != "" {
		return detectedResult.ConfigPath
	}

	return ""
}

func generateAllTemplates(_ string) {
	workingDir, _ := os.Getwd()
	mgr, ok := templateManager.(*templates.EmbeddedTemplateManager)
	if !ok {
		uiRenderer.RenderError("template manager is not the expected concrete type")
		return
	}

	templates.AllowImplicitInvocation = allowImplicitFlag
	strategy := templates.StrategyFor(detectedResult.ToolType, mgr)
	results := strategy.GenerateAll(workingDir, forceFlag)

	if len(results) == 0 {
		uiRenderer.RenderWarning("No templates available")
		return
	}

	var successCount, failCount int
	for _, result := range results {
		if result.Success {
			successCount++
			uiRenderer.RenderSuccess("Generated: " + result.FilePath)
		} else {
			failCount++
			uiRenderer.RenderError("Failed: " + result.Message)
		}
	}

	uiRenderer.RenderSuccess("Generation complete: " + FormatCount(successCount, "succeeded") + ", " + FormatCount(failCount, "failed"))
}

func generateSingleTemplate(name string) {
	tmpl, err := templateManager.GetByName(name)
	if err != nil {
		uiRenderer.RenderError("Template not found: " + name)
		return
	}
	generateOneViaStrategy(tmpl)
}

func generateInteractively() {
	tmpls, err := templateManager.ListAvailable(detectedResult.ToolType)
	if err != nil {
		uiRenderer.RenderError("Failed to list templates: " + err.Error())
		return
	}

	if len(tmpls) == 0 {
		uiRenderer.RenderWarning("No templates available")
		return
	}

	selected, err := uiRenderer.SelectTemplate(tmpls)
	if err != nil {
		uiRenderer.RenderError("Selection cancelled: " + err.Error())
		return
	}
	generateOneViaStrategy(selected)
}

// generateOneViaStrategy dispatches a single-template generation through
// the registered GenerationStrategy for the active tool. This keeps Codex
// skills (and any future archetype) consistent between --all and
// per-template invocations.
func generateOneViaStrategy(tmpl templates.TemplateMeta) {
	workingDir, _ := os.Getwd()
	mgr, ok := templateManager.(*templates.EmbeddedTemplateManager)
	if !ok {
		uiRenderer.RenderError("template manager is not the expected concrete type")
		return
	}

	templates.AllowImplicitInvocation = allowImplicitFlag
	results := templates.StrategyFor(detectedResult.ToolType, mgr).GenerateOne(workingDir, tmpl, forceFlag)
	for _, r := range results {
		if r.Success {
			uiRenderer.RenderSuccess("Generated: " + r.FilePath)
		} else {
			uiRenderer.RenderError(r.Message)
		}
	}
}

// FormatCount renders "<n> <label>" for the generation summary. The labels are
// past-tense verbs ("succeeded", "failed"), so they need no pluralization.
func FormatCount(count int, label string) string {
	return strconv.Itoa(count) + " " + label
}
