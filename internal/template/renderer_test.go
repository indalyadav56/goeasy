package template

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	// Create a mock embed.FS for testing
	var mockFS embed.FS
	renderer := NewRenderer(mockFS)
	if renderer == nil {
		t.Error("NewRenderer() returned nil")
	}
}

func TestRenderToFile(t *testing.T) {
	// Create test templates
	testFS := createTestTemplates(t)
	renderer := NewRenderer(testFS)

	// Create temporary directory for test output
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "test_output.go")

	// Test data
	testData := Data{
		Package:    "test",
		ModuleName: "github.com/test/project",
		EntityName: "User",
		IsMonolith: false,
		UseGin:     false,
	}

	// Test rendering
	err := renderer.RenderToFile("testdata/simple.tmpl", outputFile, testData)
	if err != nil {
		t.Fatalf("RenderToFile() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("RenderToFile() did not create output file")
	}

	// Read and verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "package test\n\n// Module: github.com/test/project\n// Entity: User\n"
	if string(content) != expected {
		t.Errorf("RenderToFile() content = %q, want %q", string(content), expected)
	}
}

func TestRenderToFileWithInvalidTemplate(t *testing.T) {
	testFS := createTestTemplates(t)
	renderer := NewRenderer(testFS)

	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "test_output.go")

	testData := Data{Package: "test"}

	// Test with non-existent template
	err := renderer.RenderToFile("nonexistent.tmpl", outputFile, testData)
	if err == nil {
		t.Error("RenderToFile() expected error for non-existent template, got nil")
	}
}

func TestTemplateFunctions(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     Data
		expected string
	}{
		{
			name:     "ToPascalCase",
			template: "{{.EntityName | ToPascalCase}}",
			data:     Data{EntityName: "user"},
			expected: "User",
		},
		{
			name:     "ToCamelCase",
			template: "{{.EntityName | ToCamelCase}}",
			data:     Data{EntityName: "user"},
			expected: "user",
		},
		{
			name:     "ToLower",
			template: "{{.EntityName | ToLower}}",
			data:     Data{EntityName: "USER"},
			expected: "user",
		},
		{
			name:     "ToPascalCase with underscore",
			template: "{{.EntityName | ToPascalCase}}",
			data:     Data{EntityName: "user_profile"},
			expected: "UserProfile",
		},
		{
			name:     "ToPascalCase with hyphen",
			template: "{{.EntityName | ToPascalCase}}",
			data:     Data{EntityName: "user-profile"},
			expected: "UserProfile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFS := createTestTemplatesWithContent(t, map[string]string{
				"test.tmpl": tt.template,
			})
			renderer := NewRenderer(testFS)

			tempDir := t.TempDir()
			outputFile := filepath.Join(tempDir, "test_output.txt")

			err := renderer.RenderToFile("test.tmpl", outputFile, tt.data)
			if err != nil {
				t.Fatalf("RenderToFile() error = %v", err)
			}

			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			if strings.TrimSpace(string(content)) != tt.expected {
				t.Errorf("Template function result = %q, want %q", strings.TrimSpace(string(content)), tt.expected)
			}
		})
	}
}

// Helper functions for tests
func createTestTemplates(t *testing.T) embed.FS {
	return createTestTemplatesWithContent(t, map[string]string{
		"testdata/simple.tmpl": "package {{.Package}}\n\n// Module: {{.ModuleName}}\n// Entity: {{.EntityName}}\n",
	})
}

func createTestTemplatesWithContent(t *testing.T, templates map[string]string) embed.FS {
	// Create a temporary directory structure for test templates
	tempDir := t.TempDir()
	
	for path, content := range templates {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write template file %s: %v", fullPath, err)
		}
	}
	
	// Note: In a real scenario, you'd use embed.FS differently
	// For testing purposes, we'll create a mock embed.FS
	// This is a simplified approach for unit testing
	return embed.FS{}
}
