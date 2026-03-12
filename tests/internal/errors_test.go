package internal_test

import (
	"errors"
	"testing"

	"open-spdd/internal"
)

func TestErrToolNotDetected(t *testing.T) {
	if internal.ErrToolNotDetected == nil {
		t.Error("ErrToolNotDetected should not be nil")
	}

	expectedMessage := "no AI coding tool environment detected"
	if internal.ErrToolNotDetected.Error() != expectedMessage {
		t.Errorf("ErrToolNotDetected.Error() = %v, want %v", internal.ErrToolNotDetected.Error(), expectedMessage)
	}
}

func TestErrFileExists(t *testing.T) {
	if internal.ErrFileExists == nil {
		t.Error("ErrFileExists should not be nil")
	}

	expectedMessage := "file already exists (use --force to overwrite)"
	if internal.ErrFileExists.Error() != expectedMessage {
		t.Errorf("ErrFileExists.Error() = %v, want %v", internal.ErrFileExists.Error(), expectedMessage)
	}
}

func TestErrTemplateNotFound(t *testing.T) {
	if internal.ErrTemplateNotFound == nil {
		t.Error("ErrTemplateNotFound should not be nil")
	}

	expectedMessage := "template not found"
	if internal.ErrTemplateNotFound.Error() != expectedMessage {
		t.Errorf("ErrTemplateNotFound.Error() = %v, want %v", internal.ErrTemplateNotFound.Error(), expectedMessage)
	}
}

func TestErrExistingFileNoMarkers(t *testing.T) {
	if internal.ErrExistingFileNoMarkers == nil {
		t.Error("ErrExistingFileNoMarkers should not be nil")
	}

	expectedMessage := "file exists without SPDD markers"
	if internal.ErrExistingFileNoMarkers.Error() != expectedMessage {
		t.Errorf("ErrExistingFileNoMarkers.Error() = %v, want %v", internal.ErrExistingFileNoMarkers.Error(), expectedMessage)
	}
}

func TestErrorsAreDistinct(t *testing.T) {
	allErrors := []error{
		internal.ErrToolNotDetected,
		internal.ErrFileExists,
		internal.ErrTemplateNotFound,
		internal.ErrExistingFileNoMarkers,
	}

	for i, err1 := range allErrors {
		for j, err2 := range allErrors {
			if i != j && errors.Is(err1, err2) {
				t.Errorf("Errors %d and %d should be distinct: %v vs %v", i, j, err1, err2)
			}
		}
	}
}

func TestErrorsCanBeWrapped(t *testing.T) {
	wrapped := internal.ErrToolNotDetected

	if !errors.Is(wrapped, internal.ErrToolNotDetected) {
		t.Error("errors.Is should match wrapped error")
	}
}
