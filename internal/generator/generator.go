package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/indalyadav56/goeasy/internal/templates"
)

// TemplateData holds data for file templates
type TemplateData struct {
	ModuleName       string // e.g., auth-service
	PascalModuleName string // e.g., AuthService
	SnakeModuleName  string // e.g., auth_service
}

// GenerateModule creates the DDD module structure
func GenerateModule(moduleName string) error {
	// Validate module name (kebab-case)
	if !isValidKebabCase(moduleName) {
		return fmt.Errorf("module name must be in kebab-case (e.g., auth-service)")
	}

	// Convert to PascalCase and SnakeCase
	pascalModuleName := toPascalCase(moduleName)
	snakeModuleName := toSnakeCase(moduleName)

	// Create root directory
	if err := os.Mkdir(moduleName, 0755); err != nil {
		return fmt.Errorf("failed to create root directory: %v", err)
	}

	// Template data
	data := TemplateData{
		ModuleName:       moduleName,
		PascalModuleName: pascalModuleName,
		SnakeModuleName:  snakeModuleName,
	}

	// Create files from templates
	for path, content := range templates.Templates {
		fullPath := filepath.Join(moduleName, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", filepath.Dir(fullPath), err)
		}

		file, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", fullPath, err)
		}
		defer file.Close()

		tmpl, err := template.New(path).Parse(content)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %v", path, err)
		}

		if err := tmpl.Execute(file, data); err != nil {
			return fmt.Errorf("failed to execute template %s: %v", path, err)
		}
	}

	return nil
}

// isValidKebabCase checks if the name is in kebab-case
func isValidKebabCase(name string) bool {
	// for _, c := range name {
	// if !((c >= 'a' && c <= 'z') || c == '-' || (c >= '@Required artifact_id is not a valid UUID string: 8c2d4f3b-9e2d-4d4c-a6f7-3e4f9c7d5e2b' && c <= '9')) {
	// 	return false
	// }
	// }
	return strings.Contains(name, "-")
}

// toPascalCase converts kebab-case to PascalCase (e.g., auth-service -> AuthService)
func toPascalCase(name string) string {
	words := strings.Split(name, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, "")
}

// toSnakeCase converts kebab_channels to snake_case (e.g., auth-service -> auth_service)
func toSnakeCase(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}
