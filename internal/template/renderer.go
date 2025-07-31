package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Renderer handles template rendering operations
type Renderer struct {
	templateFS embed.FS
}

// NewRenderer creates a new template renderer
func NewRenderer(templateFS embed.FS) *Renderer {
	return &Renderer{
		templateFS: templateFS,
	}
}

// Data represents template data
type Data struct {
	Package     string
	ProjectRoot string
	ModuleName  string
	EntityName  string
	Entities    []string
	IsMonolith  bool
	UseGin      bool
	// Import paths for different architectures
	HandlerImport     string
	ServiceImport     string
	RepositoryImport  string
	EntityImport      string
	InfraImport       string
	RoutesImport      string
}

// RenderToFile renders a template to a file
func (r *Renderer) RenderToFile(templatePath, outputPath string, data Data) error {
	// Define custom template functions
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

	// Ensure template path uses forward slashes for embedded filesystem
	templatePathNormalized := strings.ReplaceAll(templatePath, "\\", "/")

	// Read template content from embedded filesystem
	templateContent, err := r.templateFS.ReadFile(templatePathNormalized)
	if err != nil {
		return fmt.Errorf("failed to read embedded template file %s: %w", templatePathNormalized, err)
	}

	// Parse template from content with custom functions
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).Parse(string(templateContent))
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
