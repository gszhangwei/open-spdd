package detector

import (
	"os"
	"path/filepath"
)

// DetectorService defines the interface for detecting AI coding environments.
type DetectorService interface {
	Detect(workingDir string) DetectResult
	GetConfigDirPath(tool AIToolType, workingDir string) string
}

// DefaultDetector implements DetectorService for detecting AI tool environments.
type DefaultDetector struct{}

// NewDefaultDetector creates a new DefaultDetector instance.
func NewDefaultDetector() *DefaultDetector {
	return &DefaultDetector{}
}

// Detect scans the working directory for AI tool signature files.
func (d *DefaultDetector) Detect(workingDir string) DetectResult {
	if workingDir == "" {
		var err error
		workingDir, err = os.Getwd()
		if err != nil {
			return DetectResult{
				ToolType: Unknown,
				IsValid:  false,
				Message:  "failed to get current working directory",
			}
		}
	}

	toolTypes := []AIToolType{Cursor, ClaudeCode, Antigravity}

	for _, tool := range toolTypes {
		signatures := tool.GetSignatureFiles()
		for _, sig := range signatures {
			sigPath := filepath.Join(workingDir, sig)
			if _, err := os.Stat(sigPath); err == nil {
				configPath := d.GetConfigDirPath(tool, workingDir)
				return DetectResult{
					ToolType:   tool,
					ConfigPath: configPath,
					IsValid:    true,
					Message:    tool.String() + " environment detected",
				}
			}
		}
	}

	return DetectResult{
		ToolType: Unknown,
		IsValid:  false,
		Message:  "no AI coding tool environment detected",
	}
}

// GetConfigDirPath returns the absolute path to the config directory for a given tool.
func (d *DefaultDetector) GetConfigDirPath(tool AIToolType, workingDir string) string {
	configDir := tool.GetConfigDir()
	if configDir == "" {
		return ""
	}
	return filepath.Join(workingDir, configDir)
}
