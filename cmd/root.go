package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"open-spdd/internal/detector"
	"open-spdd/internal/templates"
	"open-spdd/internal/ui"
)

var (
	det             detector.DetectorService
	uiRenderer      ui.UIRenderer
	templateManager templates.TemplateManager
	detectedResult  detector.DetectResult
	toolFlag        string
)

var rootCmd = &cobra.Command{
	Use:   "openspdd",
	Short: "AI Coding Assistant Command Template Manager",
	Long: `SPDD (Structured Prompt-Driven Development) CLI tool for managing
AI coding assistant command templates.

Supports Cursor, Claude Code, and Antigravity environments.
Auto-detects your current environment and manages command templates.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		det = detector.NewDefaultDetector()
		uiRenderer = ui.NewCharmUIRenderer()
		templateManager = templates.NewEmbeddedTemplateManager()

		workingDir, _ := os.Getwd()

		if toolFlag != "" {
			toolType := parseToolFlag(toolFlag)
			detectedResult = detector.DetectResult{
				ToolType:   toolType,
				ConfigPath: det.GetConfigDirPath(toolType, workingDir),
				IsValid:    toolType != detector.Unknown,
				Message:    "tool manually specified: " + toolType.String(),
			}
		} else {
			detectedResult = det.Detect(workingDir)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&toolFlag, "tool", "t", "", "Manually specify tool type (cursor, claude-code, antigravity)")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func parseToolFlag(flag string) detector.AIToolType {
	switch strings.ToLower(flag) {
	case "cursor":
		return detector.Cursor
	case "claude-code", "claude":
		return detector.ClaudeCode
	case "antigravity":
		return detector.Antigravity
	default:
		return detector.Unknown
	}
}
