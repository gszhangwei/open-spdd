package cmd_test

import (
	"testing"

	"github.com/gszhangwei/open-spdd/cmd"
)

func TestFormatCount(t *testing.T) {
	tests := []struct {
		count int
		label string
		want  string
	}{
		{0, "failed", "0 failed"},
		{1, "succeeded", "1 succeeded"},
		{5, "succeeded", "5 succeeded"},
		{120, "failed", "120 failed"},
	}

	for _, tt := range tests {
		if got := cmd.FormatCount(tt.count, tt.label); got != tt.want {
			t.Errorf("FormatCount(%d, %q) = %q, want %q", tt.count, tt.label, got, tt.want)
		}
	}
}
