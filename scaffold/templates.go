package scaffolder

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// // go:embed templates/*.tmpl
var templateFS embed.FS

type TemplateRenderer struct {
	funcMap template.FuncMap
}

func NewTemplateRenderer() *TemplateRenderer {
	funcMap := template.FuncMap{
		"ToLower": func(s string) string {
			return strings.ToLower(s)
		},
		"ToUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"ToCamelCase": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
		"ToPascalCase": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
	}

	return &TemplateRenderer{
		funcMap: funcMap,
	}
}

func (tr *TemplateRenderer) RenderToFile(templatePath, outputPath string, data any) error {
	// Read template content from embedded filesystem
	templateContent, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read embedded template file %s: %w", templatePath, err)
	}

	// Parse template from content with custom functions
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(tr.funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template content: %w", err)
	}

	// Create file to write rendered content
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
