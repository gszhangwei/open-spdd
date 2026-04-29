package ui_test

import (
	"testing"

	"github.com/gszhangwei/open-spdd/internal/ui"
)

func TestStylesAreDefined(t *testing.T) {
	if ui.SuccessStyle == nil {
		t.Error("SuccessStyle should not be nil")
	}
	if ui.ErrorStyle == nil {
		t.Error("ErrorStyle should not be nil")
	}
	if ui.WarningStyle == nil {
		t.Error("WarningStyle should not be nil")
	}
	if ui.InfoStyle == nil {
		t.Error("InfoStyle should not be nil")
	}
	if ui.HeaderStyle == nil {
		t.Error("HeaderStyle should not be nil")
	}
}
