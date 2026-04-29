package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gszhangwei/open-spdd/internal/detector"
	"github.com/gszhangwei/open-spdd/internal/templates"
	"github.com/gszhangwei/open-spdd/internal/ui"
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

Supports Cursor, Claude Code, Antigravity, and GitHub Copilot environments.
Auto-detects your current environment and manages command templates.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		det = detector.NewDefaultDetector()
		uiRenderer = ui.NewCharmUIRenderer()
		templateManager = templates.NewEmbeddedTemplateManager()

		workingDir, _ := os.Getwd()

		if toolFlag != "" {
			toolType := ParseToolFlag(toolFlag)
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
	rootCmd.PersistentFlags().StringVarP(&toolFlag, "tool", "t", "", "Manually specify tool type (cursor, claude-code, antigravity, github-copilot)")
}

func SetVersion(v string) {
	v = strings.TrimSpace(v)
	if v == "" {
		v = "dev"
	}
	rootCmd.Version = v
	rootCmd.SetVersionTemplate("openspdd {{.Version}}\n")

	// Cobra registers the `--version` flag lazily when Version is non-empty
	// and the command is executed. To bind a `-v` shorthand we touch the
	// flag explicitly: declare it ourselves if absent, otherwise patch the
	// existing definition.
	if f := rootCmd.Flags().Lookup("version"); f != nil {
		f.Shorthand = "v"
		f.Usage = "Print openspdd version and exit"
	} else {
		rootCmd.Flags().BoolP("version", "v", false, "Print openspdd version and exit")
	}
}

// RootCommand exposes the root *cobra.Command for testing. Production code
// should use Execute() instead.
func RootCommand() *cobra.Command {
	return rootCmd
}

// Execute runs the root command.
func Execute() {
	maybePrintPathHint()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ParseToolFlag converts a tool flag string to AIToolType.
func ParseToolFlag(flag string) detector.AIToolType {
	switch strings.ToLower(flag) {
	case "cursor":
		return detector.Cursor
	case "claude-code", "claude":
		return detector.ClaudeCode
	case "antigravity":
		return detector.Antigravity
	case "github-copilot", "copilot", "gh-copilot":
		return detector.GitHubCopilot
	default:
		return detector.Unknown
	}
}
