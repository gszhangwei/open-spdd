package cmd_test

import (
	"testing"

	"github.com/gszhangwei/open-spdd/cmd"
	"github.com/gszhangwei/open-spdd/internal/detector"
)

func TestShouldPromptForTool(t *testing.T) {
	detected := detector.DetectResult{ToolType: detector.ClaudeCode, IsValid: true}
	none := detector.DetectResult{ToolType: detector.Unknown, IsValid: false}

	tests := []struct {
		name       string
		outputFlag string
		toolFlag   string
		detected   detector.DetectResult
		want       bool
	}{
		{"nothing given, nothing detected", "", "", none, true},
		{"auto-detected tool wins", "", "", detected, false},
		{"--tool given", "", "cursor", none, false},
		{"--output given", "out", "", none, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmd.ShouldPromptForTool(tt.outputFlag, tt.toolFlag, tt.detected); got != tt.want {
				t.Errorf("ShouldPromptForTool(%q, %q, valid=%v) = %v, want %v", tt.outputFlag, tt.toolFlag, tt.detected.IsValid, got, tt.want)
			}
		})
	}
}
