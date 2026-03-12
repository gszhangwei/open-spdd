package templates

import "embed"

//go:embed all:data
var embeddedTemplates embed.FS
