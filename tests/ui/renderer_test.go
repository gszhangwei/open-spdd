package ui_test

import (
	"testing"

	"github.com/gszhangwei/open-spdd/internal/templates"
	"github.com/gszhangwei/open-spdd/internal/ui"
)

func TestNewCharmUIRenderer(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	if renderer == nil {
		t.Error("NewCharmUIRenderer() returned nil")
	}
}

func TestUIRenderer_Interface(t *testing.T) {
	var _ ui.UIRenderer = (*ui.CharmUIRenderer)(nil)
}

func TestCharmUIRenderer_RenderSuccess(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	renderer.RenderSuccess("test success message")
}

func TestCharmUIRenderer_RenderError(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	renderer.RenderError("test error message")
}

func TestCharmUIRenderer_RenderWarning(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	renderer.RenderWarning("test warning message")
}

func TestCharmUIRenderer_RenderTable_Empty(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	renderer.RenderTable([]string{}, [][]string{})
}

func TestCharmUIRenderer_RenderTable_HeadersOnly(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	renderer.RenderTable([]string{"Name", "Category", "Description"}, [][]string{})
}

func TestCharmUIRenderer_RenderTable_WithRows(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	headers := []string{"Name", "Category", "Description"}
	rows := [][]string{
		{"Template 1", "Development", "A test template"},
		{"Template 2", "Testing", "Another template"},
	}
	renderer.RenderTable(headers, rows)
}

func TestCharmUIRenderer_RenderTable_LongContent(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	headers := []string{"A", "B"}
	rows := [][]string{
		{"Short", "This is a very long description that should be handled properly"},
	}
	renderer.RenderTable(headers, rows)
}

func TestCharmUIRenderer_SelectTemplate_EmptyList(t *testing.T) {
	renderer := ui.NewCharmUIRenderer()
	_, err := renderer.SelectTemplate([]templates.TemplateMeta{})
	if err == nil {
		t.Error("SelectTemplate() with empty list should return error")
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  string
	}{
		{
			name:  "short string",
			input: "abc",
			width: 5,
			want:  "abc  ",
		},
		{
			name:  "exact width",
			input: "abcde",
			width: 5,
			want:  "abcde",
		},
		{
			name:  "longer than width",
			input: "abcdef",
			width: 5,
			want:  "abcdef",
		},
		{
			name:  "empty string",
			input: "",
			width: 3,
			want:  "   ",
		},
		{
			name:  "zero width",
			input: "abc",
			width: 0,
			want:  "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ui.PadRight(tt.input, tt.width)
			if got != tt.want {
				t.Errorf("PadRight(%q, %d) = %q, want %q", tt.input, tt.width, got, tt.want)
			}
		})
	}
}
