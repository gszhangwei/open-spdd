package cmd

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/gszhangwei/open-spdd/internal/detector"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize command template directory for detected AI tool",
	Long: `Initialize the command template directory for the detected AI coding tool.
Creates the necessary directory structure for storing command templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		workingDir, _ := os.Getwd()

		if !detectedResult.IsValid {
			tool := selectToolInteractively()
			if tool == detector.Unknown {
				uiRenderer.RenderError("No tool selected")
				return
			}
			detectedResult = detector.DetectResult{
				ToolType:   tool,
				ConfigPath: det.GetConfigDirPath(tool, workingDir),
				IsValid:    true,
				Message:    "tool manually selected: " + tool.String(),
			}
		}

		configPath := detectedResult.ConfigPath
		if configPath == "" {
			uiRenderer.RenderError("Could not determine config directory")
			return
		}

		if _, err := os.Stat(configPath); err == nil {
			uiRenderer.RenderWarning("Directory already exists: " + configPath)
			if !uiRenderer.Confirm("Do you want to continue anyway?") {
				uiRenderer.RenderWarning("Initialization cancelled")
				return
			}
		}

		if err := os.MkdirAll(configPath, 0755); err != nil {
			uiRenderer.RenderError("Failed to create directory: " + err.Error())
			return
		}

		if detectedResult.ToolType == detector.GitHubCopilot {
			promptsDir := filepath.Join(workingDir, ".github", "copilot-prompts")
			if err := os.MkdirAll(promptsDir, 0755); err != nil {
				uiRenderer.RenderError("Failed to create copilot-prompts directory: " + err.Error())
				return
			}
		}

		uiRenderer.RenderSuccess("Initialized " + detectedResult.ToolType.String() + " command directory at: " + configPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func selectToolInteractively() detector.AIToolType {
	options := []huh.Option[string]{
		huh.NewOption("Cursor", "cursor"),
		huh.NewOption("Claude Code", "claude-code"),
		huh.NewOption("Antigravity", "antigravity"),
		huh.NewOption("GitHub Copilot", "github-copilot"),
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select your AI coding tool").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return detector.Unknown
	}

	return ParseToolFlag(selected)
}
