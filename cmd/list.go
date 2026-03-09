package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	categoryFlag string
	quietFlag    bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available command templates",
	Long:    `List all available command templates that can be generated.`,
	Run: func(cmd *cobra.Command, args []string) {
		tmpls, err := templateManager.ListAll()
		if err != nil {
			uiRenderer.RenderError("Failed to list templates: " + err.Error())
			return
		}

		if categoryFlag != "" {
			var filtered []interface{}
			for _, t := range tmpls {
				if t.Category == categoryFlag {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				uiRenderer.RenderWarning("No templates found in category: " + categoryFlag)
				return
			}
		}

		if quietFlag {
			for _, t := range tmpls {
				if categoryFlag != "" && t.Category != categoryFlag {
					continue
				}
				fmt.Println(t.ID)
			}
			return
		}

		var rows [][]string
		for _, t := range tmpls {
			if categoryFlag != "" && t.Category != categoryFlag {
				continue
			}
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
	rootCmd.AddCommand(listCmd)
}
