package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gszhangwei/open-spdd/internal/templates"
)

var (
	categoryFlag string
	quietFlag    bool
	optionalFlag bool
	listAllFlag  bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available command templates",
	Long: `List command templates that can be generated.

By default, shows templates available for the detected environment (core + tool-specific).
Use --optional to show optional templates, or --all to show all templates across categories.`,
	Run: func(cmd *cobra.Command, args []string) {
		var tmpls []templates.TemplateMeta
		var err error

		if listAllFlag {
			tmpls, err = templateManager.ListAll()
		} else if optionalFlag {
			tmpls, err = templateManager.ListOptional()
		} else {
			tmpls, err = templateManager.ListAvailable(detectedResult.ToolType)
		}

		if err != nil {
			uiRenderer.RenderError("Failed to list templates: " + err.Error())
			return
		}

		if categoryFlag != "" {
			var filtered []templates.TemplateMeta
			for _, t := range tmpls {
				if t.Category == categoryFlag {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				uiRenderer.RenderWarning("No templates found in category: " + categoryFlag)
				return
			}
			tmpls = filtered
		}

		if quietFlag {
			for _, t := range tmpls {
				fmt.Println(t.ID)
			}
			return
		}

		var rows [][]string
		for _, t := range tmpls {
			rows = append(rows, []string{t.Name, t.Category, t.Description})
		}

		if len(rows) == 0 {
			uiRenderer.RenderWarning("No templates available")
			return
		}

		uiRenderer.RenderTable([]string{"Name", "Category", "Description"}, rows)
	},
}

func init() {
	listCmd.Flags().StringVarP(&categoryFlag, "category", "c", "", "Filter by category")
	listCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Output only template names (for piping)")
	listCmd.Flags().BoolVar(&optionalFlag, "optional", false, "List optional templates")
	listCmd.Flags().BoolVar(&listAllFlag, "all", false, "List all templates across all categories")
	rootCmd.AddCommand(listCmd)
}
