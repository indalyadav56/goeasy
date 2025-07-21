package goembed

import (
	"embed"
)

//go:embed templates/*.tmpl templates/*.tmpl
var TemplateFS embed.FS
