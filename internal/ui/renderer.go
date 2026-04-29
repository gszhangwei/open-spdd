package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/gszhangwei/open-spdd/internal/templates"
)

// UIRenderer defines the interface for terminal output rendering.
type UIRenderer interface {
	RenderSuccess(msg string)
	RenderError(msg string)
	RenderWarning(msg string)
	RenderTable(headers []string, rows [][]string)
	SelectTemplate(templates []templates.TemplateMeta) (templates.TemplateMeta, error)
	Confirm(prompt string) bool
}

// CharmUIRenderer implements UIRenderer using Charm libraries.
type CharmUIRenderer struct{}

// NewCharmUIRenderer creates a new CharmUIRenderer instance.
func NewCharmUIRenderer() *CharmUIRenderer {
	return &CharmUIRenderer{}
}

// RenderSuccess prints a success message in green.
func (r *CharmUIRenderer) RenderSuccess(msg string) {
	SuccessStyle.Println("✓ " + msg)
}

// RenderError prints an error message in red.
func (r *CharmUIRenderer) RenderError(msg string) {
	ErrorStyle.Println("✗ " + msg)
}

// RenderWarning prints a warning message in yellow.
func (r *CharmUIRenderer) RenderWarning(msg string) {
	WarningStyle.Println("⚠ " + msg)
}

// RenderTable prints a formatted table with headers and rows.
func (r *CharmUIRenderer) RenderTable(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var headerLine strings.Builder
	var separatorLine strings.Builder
	for i, h := range headers {
		if i > 0 {
			headerLine.WriteString("  ")
			separatorLine.WriteString("  ")
		}
		headerLine.WriteString(PadRight(h, colWidths[i]))
		separatorLine.WriteString(strings.Repeat("-", colWidths[i]))
	}

	HeaderStyle.Println(headerLine.String())
	fmt.Println(separatorLine.String())

	for _, row := range rows {
		var rowLine strings.Builder
		for i, cell := range row {
			if i > 0 {
				rowLine.WriteString("  ")
			}
			if i < len(colWidths) {
				rowLine.WriteString(PadRight(cell, colWidths[i]))
			} else {
				rowLine.WriteString(cell)
			}
		}
		fmt.Println(rowLine.String())
	}
}

// SelectTemplate displays an interactive selection form for templates.
func (r *CharmUIRenderer) SelectTemplate(tmpls []templates.TemplateMeta) (templates.TemplateMeta, error) {
	if len(tmpls) == 0 {
		return templates.TemplateMeta{}, fmt.Errorf("no templates available")
	}

	options := make([]huh.Option[int], len(tmpls))
	for i, t := range tmpls {
		label := t.Name
		if t.Description != "" {
			label += " - " + t.Description
		}
		options[i] = huh.NewOption(label, i)
	}

	var selected int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a template").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return templates.TemplateMeta{}, err
	}

	return tmpls[selected], nil
}

// Confirm displays a confirmation prompt and returns the user's response.
func (r *CharmUIRenderer) Confirm(prompt string) bool {
	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(prompt).
				Affirmative("Yes").
				Negative("No").
				Value(&confirmed),
		),
	)

	if err := form.Run(); err != nil {
		return false
	}

	return confirmed
}

// PadRight pads a string with spaces to the specified width.
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
