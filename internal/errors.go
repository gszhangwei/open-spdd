package internal

import "errors"

var (
	ErrToolNotDetected  = errors.New("no AI coding tool environment detected")
	ErrFileExists       = errors.New("file already exists (use --force to overwrite)")
	ErrTemplateNotFound = errors.New("template not found")
)
