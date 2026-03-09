package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"open-spdd/internal/templates"
)

var (
	forceFlag  bool
	allFlag    bool
	outputFlag string
)

var generateCmd = &cobra.Command{
	Use:     "generate [template-name]",
	Aliases: []string{"gen", "g"},
	Short:   "Generate command template file",
	Long: `Generate a command template file to the detected AI tool's config directory.
If no template name is specified, an interactive selection will be shown.`,
	Run: func(cmd *cobra.Command, args []string) {
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
			generateSingleTemplate(args[0], targetDir)
			return
		}

		generateInteractively(targetDir)
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Overwrite existing files")
	generateCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Generate all available templates")
	generateCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Custom output directory (overrides detection)")
	rootCmd.AddCommand(generateCmd)
}

func determineTargetDir() string {
	if outputFlag != "" {
		return outputFlag
	}

	if detectedResult.IsValid && detectedResult.ConfigPath != "" {
		return detectedResult.ConfigPath
	}

	workingDir, _ := os.Getwd()
	if workingDir != "" {
		return workingDir
	}

	return ""
}

func generateAllTemplates(targetDir string) {
	tmpls, err := templateManager.ListAll()
	if err != nil {
		uiRenderer.RenderError("Failed to list templates: " + err.Error())
		return
	}

	if len(tmpls) == 0 {
		uiRenderer.RenderWarning("No templates available")
		return
	}

	var successCount, failCount int
	for _, t := range tmpls {
		result := generateTemplate(t, targetDir)
		if result.Success {
			successCount++
			uiRenderer.RenderSuccess("Generated: " + result.FilePath)
		} else {
			failCount++
			uiRenderer.RenderError("Failed: " + t.Name + " - " + result.Message)
		}
	}

	uiRenderer.RenderSuccess("Generation complete: " + formatCount(successCount, "succeeded") + ", " + formatCount(failCount, "failed"))
}

func generateSingleTemplate(name, targetDir string) {
	tmpl, err := templateManager.GetByName(name)
	if err != nil {
		uiRenderer.RenderError("Template not found: " + name)
		return
	}

	result := generateTemplate(tmpl, targetDir)
	if result.Success {
		uiRenderer.RenderSuccess("Generated: " + result.FilePath)
	} else {
		uiRenderer.RenderError(result.Message)
	}
}

func generateInteractively(targetDir string) {
	tmpls, err := templateManager.ListAll()
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

	result := generateTemplate(selected, targetDir)
	if result.Success {
		uiRenderer.RenderSuccess("Generated: " + result.FilePath)
	} else {
		uiRenderer.RenderError(result.Message)
	}
}

func generateTemplate(tmpl templates.TemplateMeta, targetDir string) templates.GenerateResult {
	filename := tmpl.ID + ".md"
	targetPath := filepath.Join(targetDir, filename)

	req := templates.GenerateRequest{
		TemplateName: tmpl.Name,
		TargetPath:   targetPath,
		Force:        forceFlag,
	}

	return templateManager.Generate(req)
}

func formatCount(count int, label string) string {
	if count == 1 {
		return "1 " + label[:len(label)-2]
	}
	return formatInt(count) + " " + label
}

func formatInt(n int) string {
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}
