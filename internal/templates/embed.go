package templates

import "embed"

//go:embed data/*.md
var embeddedTemplates embed.FS
